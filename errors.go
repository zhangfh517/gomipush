package gomipush

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// {
//   "result": "error",
//   "reason": "Must input one of: regid/topic/alias/userAccount/miid/geoId/groupId/region",
//   "trace_id": "Xcm35b38494233429183nK",
//   "code": 21305,
//   "description": "缺少必要的参数"
// }

type Error struct {
	AppStatus int    `json:"-"`
	AppReason string `json:"-"`

	Result      string `json:"result,omitempty"`
	Reason      string `json:"reason,omitempty"`
	TraceId     string `json:"trace_id,omitempty"`
	Code        int    `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
	ErrorCode   int    `json:"errorCode,omitempty"`
}

func (e *Error) Error() string {
	var code int = e.Code
	if e.Code == 0 {
		code = e.ErrorCode
	}
	return fmt.Sprintf("mipush: Error %d (%s): %s [result: %s, reason: %s, description: %s, code: %d, tracdID: %s]", e.AppStatus, http.StatusText(e.AppStatus), e.AppReason, e.Result, e.Reason, e.Description, code, e.TraceId)

}

func checkResponse(res *http.Response) error {
	// 200-299 are valid status codes
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return createResponseError(res)
	}
	return &Error{AppStatus: res.StatusCode}
}

func createResponseError(res *http.Response) error {

	if res.Body == nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Response body is nil"}
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Read From response body error"}
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(data))

	errReply := new(Error)
	err = json.Unmarshal(data, &errReply)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Unmarshall response body content error"}
	}
	if errReply.Code == 0 && errReply.ErrorCode == 0 {
		return nil
	}
	errReply.AppStatus = res.StatusCode
	return errReply
}
