package gomipush

import(
	"net/http"
	"context"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"strings"
	"io"
	"fmt"
	"bytes"
	"net/url"
	"time"
)

const MAX_BACKOFF_DELAY = 1024000


type Response struct {
	AppStatus int `json:"-"`
	AppReason string `json:"-"`

	Result string `json:"result"`
	TraceID string `json:"trace_id"`
	Code int `json:"code"`
	Data struct {
		Id string `json:"id"`
	} `json:"data"`
	Description string `json:"description"`
	Info string `json:"info"`
}


func newResponse(res *http.Response) (*Response, error) {
	r := &Response{
		AppStatus: res.StatusCode,
	}
	if res.Body != nil {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		if len(body) > 0 {
			if err := json.Unmarshal(body, &r); err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}

func setBodyString(req *http.Request, bodyStr string) {
	body := strings.NewReader(bodyStr)
	rc, ok := (io.Reader)(body).(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	req.Body = rc
	if body != nil {
		switch v := (io.Reader)(body).(type) {
		case *strings.Reader:
			req.ContentLength = int64(v.Len())
		case *bytes.Buffer:
			req.ContentLength = int64(v.Len())
		}
	}
}

func newRequest(method string , url string, contentType string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", contentType)
	return req, nil
}

// func buildUrl(url string) string {
// 	return "https://api.xmpush.xiaomi.com"+url
// }

func httpCall(ctx context.Context, c *http.Client, url string, method HttpMethod, authorization string, params url.Values, body string, token string) (*Response, error){
	var resp *Response

	var urlWithParams string = url
	if len(params) > 0 {
		urlWithParams += "?" + params.Encode()
	}
	req, err := newRequest(METHOD_MAP[method], urlWithParams, "application/x-www-form-urlencoded;charset=UTF-8")
	if err != nil {
		return nil, err
	}
	setBodyString(req, body)
	if len(authorization) > 0 {
		req.Header.Add("Authorization", fmt.Sprintf("key=%s", authorization))
	}
	if len(token) > 0 {
		req.Header.Add("X-PUSH-AUDIT-TOKEN", token)
	}
	if auto_switch_host && NewServerSwitch().NeedRefreshHostList() {
		req.Header.Add("X-PUSH-HOST-LIST", "true")
	}
	res, err := c.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	if err :=  checkResponse(res); err != nil {
		return nil, err
	}
	resp, err = newResponse(res)
	if err != nil {
		return nil, err
	}
	hostList := res.Header.Get("X-PUSH-HOST-LIST")
	if len(hostList) > 0 {
		NewServerSwitch().Initialize(hostList)
	}
	return resp, nil
}


type Client struct{
	c *http.Client
	security string
	token string
	proxyIp string
	proxyPort string
}

func NewClient(security string) *Client{
	return &Client{
		c : http.DefaultClient,
		security : security,
	}
}

func (c *Client) Send(msg Message) (*SenderService, error) {
	return NewSenderService(c, msg)
}
func (c *Client) Subscribe(topic SubscribedTopic) *SubscribeService {
	return NewSubscribeService(c, topic)
}
func (c *Client) Tool() *Tool {
	return NewTool(c)
}

func (c *Client) buildRequestUrl(server *Server, requestPath []string) string {
	return http_protocol + "://" + server.GetHost() + requestPath[0]
}
func (c *Client) PerformRequest(ctx context.Context, requestPath []string, retryTimes int, method HttpMethod, params url.Values, body string) (*Response, error){

	isFail := true
	tryTime := 0
	var resp *Response = nil
	sleepTime := 1
	var err error
	start := time.Now()
	log.Infof("select server for request :%v - %v", requestPath, params)
	server := NewServerSwitch().SelectServer(requestPath)

	for isFail && tryTime <= retryTimes {
		// log.Infof("tryTime: %d", tryTime)
		resp, err = httpCall(ctx, c.c, c.buildRequestUrl(server, requestPath), method, c.security, params, body, "")
		if err != nil || time.Now().Sub(start).Seconds() > 5 {
			server.DecrPriority()
		}else {
			server.IncrPriority()
		}
		if err == nil {
			isFail = false
		}else {
			tryTime += 1
			time.Sleep(time.Duration(sleepTime) * time.Second)

			if 2 * sleepTime < MAX_BACKOFF_DELAY {
				sleepTime *= 2
			}
		}
	}
	return resp, err
}

func (c *Client) Proxy(proxyIp string, proxyPort string) (*Client, error) {
	c.proxyIp = proxyIp
	c.proxyPort = proxyPort
	proxyUrl, err := url.Parse(fmt.Sprintf("%s:%s", c.proxyIp, c.proxyPort))
	if err != nil {
		return nil, fmt.Errorf("parse proxy url error")
	}

	proxyClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	c.c = proxyClient
	return c
}
func (c *Client) Token(token string) *Client {
	c.token = token
	return c
}

