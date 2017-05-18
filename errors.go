package gomipush
import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
)
// {
//   "result": "error",
//   "reason": "Must input one of: regid/topic/alias/userAccount/miid/geoId/groupId/region",
//   "trace_id": "Xcm35b38494233429183nK",
//   "code": 21305,
//   "description": "缺少必要的参数"
// }


type Error struct {
	AppStatus  	int 	`json:"-"`
	AppReason   string  `json:"-"`

	Result  	string 	`json:"result"`
	Reason  	string 	`json:"reason"`
	TraceId 	string 	`json:"trace_id"`
	Code    	int     `json:"code"`
	Description string  `json:"description"`
}
func (e *Error) Error() string {
	return fmt.Sprintf("mipush: Error %d (%s): %s [result: %s, reason: %s, description: %s, code: %d, tracdID: %s]", e.AppStatus, http.StatusText(e.AppStatus), e.AppReason, e.Result, e.Reason, e.Description, e.Code, e. TraceId)
}

func checkResponse(res *http.Response) error {
	// 200-299 are valid status codes
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	return createResponseError(res)
}

func createResponseError(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return &Error{AppStatus: res.StatusCode}
	}

	if res.Body == nil {
		return &Error{AppStatus: res.StatusCode}
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Read From res.Body error"}
	}
	errReply := new(Error)
	err = json.Unmarshal(data, errReply)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "unmarshal error"}
	}

	if errReply.Code == 0 {
		return nil
	}

	errReply.AppStatus = res.StatusCode
	return errReply
}

