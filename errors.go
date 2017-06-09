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

	Result  	string 	`json:"result, omitempty"`
	Reason  	string 	`json:"reason, omitempty"`
	TraceId 	string 	`json:"trace_id, omitempty"`
	Code    	int     `json:"code, omitempty"`
	Description string  `json:"description, omitempty"`
}
func (e *Error) Error() string {
	if e.Code != 0 {
		return fmt.Sprintf("mipush Error: request detail(StatusCode %d (%s), reason %s), from server[result: %s, reason: %s, description: %s, code: %d, tracdID: %s]", e.AppStatus, http.StatusText(e.AppStatus), e.AppReason, e.Result, e.Reason, e.Description, e.Code, e. TraceId)
	}else {
		return fmt.Sprintf("mipush Error: request detail(StatusCode %d (%s), reason%s)", e.AppStatus, http.StatusText(e.AppStatus), e.AppReason)
	}
}

func checkResponse(res *http.Response) error {
	// 200-299 are valid status codes
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
			return createResponseError(res)
	}else {
		return &Error{AppStatus: res.StatusCode}
	}
}

func createResponseError(res *http.Response) error {
	if res.Body == nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Response body is nil"}
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Read From response body error"}
	}
	errReply := new(Error)
	err = json.Unmarshal(data, errReply)
	if err != nil {
		return &Error{AppStatus: res.StatusCode, AppReason: "Unmarshall response body content error"}
	}
	if errReply.Code == 0 {
		return nil
	}

	errReply.AppStatus = res.StatusCode
	return errReply
}

