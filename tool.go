package gomipush
import(
	"strings"
	"net/url"
	"strconv"
	"context"
)

type Tool struct {
	client *Client
	targetUrl []string
	retryTimes int
	params url.Values
}
func (t *Tool) AddParam(k, v string) *Tool{
	t.params.Set(k, v)
	return t
}

func NewTool(c *Client) *Tool {
	return &Tool{
		client: c,
	}
}
//post
func (t *Tool) CheckScheduleJobExist(jobId string) *Tool {
	t.targetUrl = V2_CHECK_SCHEDULE_JOB_EXIST
	t.AddParam("job_id", jobId)
	return t
}

//post
func (t *Tool) DeleteScheduleJob(jobId string) *Tool {
	t.targetUrl = V2_DELETE_SCHEDULE_JOB
	t.AddParam("job_id", jobId)
	return t
}
//post
func (t *Tool) DeleteScheduleJobKey(jobKey string) *Tool {
	t.targetUrl = V3_DELETE_SCHEDULE_JOB
	t.AddParam("jobkey", jobKey)
	return t
}
//post
func (t *Tool) DeleteTopic(msgId string) *Tool {
	t.targetUrl = V2_DELETE_BROADCAST_MESSAGE
	t.AddParam("id", msgId)
	return t
}
//get
func (t *Tool) QueryDeviceAliases(packageName, regId string) *Tool {
	t.targetUrl = V1_GET_ALL_TOPIC
	t.AddParam("restricted_package_name", packageName)
	t.AddParam("registration_id", regId)

	return t
}
//get
func (t *Tool) QueryDeviceUserAccounts(packageName, regId string) *Tool {
	t.targetUrl = V1_GET_ALL_ACCOUNT
	t.AddParam("restricted_package_name", packageName)
	t.AddParam("registration_id", regId)
	return t
}
//get
func (t *Tool) QueryDevicePresence(packageName string, regId []string) *Tool {
	var rid string

	if len(regId) == 1{
		rid = regId[0]
		t.targetUrl = V1_REGID_PRESENCE

	}
	if len(regId) > 1{
		t.targetUrl = V2_REGID_PRESENCE
		rid = strings.Join(regId, ",")
	}
	t.AddParam("registration_id", rid)
	t.AddParam("restricted_package_name", packageName)
	return t
}
//get
func (t *Tool) QueryInvalidRegIds() *Tool {
	t.targetUrl = []string{V1_FEEDBACK_INVALID_REGID[0].(string)}
	return t
}

func (t *Tool) QueryMessageStatus(msgId string) *Tool {
	t.targetUrl = V1_MESSAGE_STATUS
	t.AddParam("msg_id", msgId)
	return t
}
func (t *Tool) QueryMessageGroupStatus(jobKey string) *Tool {
	t.targetUrl = V1_MESSAGE_STATUS
	t.AddParam("job_key", jobKey)
	return t
}
func (t *Tool) query_message_status_time_range(beginTime int64, endTime int64) *Tool {
	t.targetUrl = V1_MESSAGES_STATUS
	t.AddParam("begin_time", strconv.FormatInt(beginTime, 64))
	t.AddParam("end_time", strconv.FormatInt(endTime, 64))
	return t
}
func (t *Tool) QueryStatData(startDate string, endDate string, packageName string) *Tool {
	t.targetUrl = V1_GET_MESSAGE_COUNTERS
	t.AddParam("start_date", startDate)
	t.AddParam("end_date", endDate)
	t.AddParam("restricted_package_name", packageName)
	return t
}
func (t *Tool) ValidateRegIds(regId []string) *Tool {
	t.targetUrl = V1_VALIDATE_REGID
	t.AddParam("registration_ids", strings.Join(regId, ","))
	return t
}
func (t *Tool) FetchAckInfo(packageName string) *Tool {
	t.targetUrl = []string{V1_EMQ_ACK_INFO[0].(string)}
	//有问题，为啥不是restrict_package_name
	t.AddParam("package_name", packageName)
	return t
}
func (t *Tool) FetchClickInfo(packageName string) *Tool {
	t.targetUrl = []string{V1_EMQ_CLICK_INFO[0].(string)}
	t.AddParam("package_name", packageName)
	return t
}
func (t *Tool) FetchInvalidRegId(packageName string) *Tool {
	t.targetUrl = []string{V1_EMQ_INVALID_REGID[0].(string)}
	t.AddParam("package_name", packageName)
	return t
}

func (t *Tool) DoGet(ctx context.Context) (*Response, error){
    return t.client.PerformRequest(ctx, t.targetUrl, t.retryTimes, HTTP_GET, t.params, "")
}
func (t *Tool) DoPost(ctx context.Context) (*Response, error){
    return t.client.PerformRequest(ctx, t.targetUrl, t.retryTimes, HTTP_POST, t.params, "")
}


