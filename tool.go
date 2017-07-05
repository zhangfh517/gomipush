package gomipush

import (
	"context"
	// log "github.com/Sirupsen/logrus"
	"net/url"
	"strconv"
	"strings"
)

//gomipush.NewClient("security").Tool().FetchClickInfo("package").DoGet(ctx)
type Tool struct {
	client        *Client
	targetUrl     []string
	retryTimes    int
	params        url.Values
	requestMethod HttpMethod
}

func (t *Tool) RetryTimes(retryTimes int) *Tool {
	t.retryTimes = retryTimes
	return t
}
func (t *Tool) RequestMethod(m HttpMethod) *Tool {
	t.requestMethod = m
	return t
}

func (t *Tool) addParam(k, v string) *Tool {
	t.params.Set(k, v)
	return t
}

func NewTool(c *Client) *Tool {
	return &Tool{
		client:        c,
		retryTimes:    2,
		params:        url.Values{},
		requestMethod: HTTP_GET,
	}
}

//post test
func (t *Tool) CheckScheduleJobExist(jobId string) *Tool {
	t.targetUrl = V2_CHECK_SCHEDULE_JOB_EXIST
	t.addParam("job_id", jobId)
	t.requestMethod = HTTP_POST
	return t
}

//post test
func (t *Tool) DeleteScheduleJob(jobId string) *Tool {
	t.targetUrl = V2_DELETE_SCHEDULE_JOB
	t.addParam("job_id", jobId)
	t.requestMethod = HTTP_POST

	return t
}

//!!post no test
func (t *Tool) DeleteScheduleJobKey(jobKey string) *Tool {
	t.targetUrl = V3_DELETE_SCHEDULE_JOB
	t.addParam("jobkey", jobKey)
	t.requestMethod = HTTP_POST

	return t
}

//!!post request successfully but not work
func (t *Tool) DeleteTopic(msgId string) *Tool {
	t.targetUrl = V2_DELETE_BROADCAST_MESSAGE
	t.addParam("msg_id", msgId)
	t.requestMethod = HTTP_POST

	return t
}

//get
func (t *Tool) QueryDeviceAliases(packageName, regId string) *Tool {
	t.targetUrl = V1_GET_ALL_TOPIC
	t.addParam("restricted_package_name", packageName)
	t.addParam("registration_id", regId)
	t.requestMethod = HTTP_GET
	return t
}

//get
func (t *Tool) QueryDeviceUserAccounts(packageName, regId string) *Tool {
	t.targetUrl = V1_GET_ALL_ACCOUNT
	t.addParam("restricted_package_name", packageName)
	t.addParam("registration_id", regId)
	t.requestMethod = HTTP_GET
	return t
}

//get
func (t *Tool) QueryDevicePresence(packageName string, regId []string) *Tool {
	var rid string

	if len(regId) == 1 {
		rid = regId[0]
		t.targetUrl = V1_REGID_PRESENCE

	}
	if len(regId) > 1 {
		t.targetUrl = V2_REGID_PRESENCE
		rid = strings.Join(regId, ",")
	}
	t.addParam("registration_id", rid)
	t.addParam("restricted_package_name", packageName)
	t.requestMethod = HTTP_GET
	return t
}

//get
func (t *Tool) QueryInvalidRegIds() *Tool {
	t.targetUrl = []string{V1_FEEDBACK_INVALID_REGID[0].(string)}
	t.requestMethod = HTTP_GET
	return t
}

func (t *Tool) QueryMessageStatus(msgId string) *Tool {
	t.targetUrl = V1_MESSAGE_STATUS
	t.addParam("msg_id", msgId)
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) QueryMessageGroupStatus(jobKey string) *Tool {
	t.targetUrl = V1_MESSAGE_STATUS
	t.addParam("job_key", jobKey)
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) QueryMessageStatusTimeRange(beginTime int64, endTime int64) *Tool {
	t.targetUrl = V1_MESSAGES_STATUS
	t.addParam("begin_time", strconv.FormatInt(beginTime, 10))
	t.addParam("end_time", strconv.FormatInt(endTime, 10))
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) QueryStatData(beginTime int64, endTime int64, packageName string) *Tool {
	t.targetUrl = V1_GET_MESSAGE_COUNTERS
	t.addParam("start_date", strconv.FormatInt(beginTime, 10))
	t.addParam("end_date", strconv.FormatInt(endTime, 10))
	t.addParam("restricted_package_name", packageName)
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) ValidateRegIds(regId []string) *Tool {
	t.targetUrl = V1_VALIDATE_REGID
	t.addParam("registration_ids", strings.Join(regId, ","))
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) FetchAckInfo(packageName string) *Tool {
	t.targetUrl = V1_EMQ_ACK_INFO
	//有问题，为啥不是restrict_package_name
	t.addParam("package_name", packageName)
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) FetchClickInfo(packageName string) *Tool {
	t.targetUrl = V1_EMQ_CLICK_INFO
	t.addParam("package_name", packageName)
	t.requestMethod = HTTP_GET
	return t
}
func (t *Tool) FetchInvalidRegId(packageName string) *Tool {
	t.targetUrl = V1_EMQ_INVALID_REGID
	t.addParam("package_name", packageName)
	t.requestMethod = HTTP_GET
	return t
}

func (t *Tool) Do(ctx context.Context) (*Response, error) {
	return t.client.PerformRequest(ctx, t.targetUrl, t.retryTimes, t.requestMethod, t.params, "")
}
