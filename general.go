package gomipush

const (
	VERSION     = "1.0.0"
	sdk_version = "GO_SDK_V1.0.0"
)

type TargetType int
type RequestType string
type SubscribeType int
type HttpMethod int
type BroadcastTopicOp string
type SenderTarget int
type NotifyType int

const (
	HTTP_GET HttpMethod = iota
	HTTP_POST
)

const (
	TARGET_TYPE_REGID TargetType = iota + 1
	TARGET_TYPE_ALIAS
	TARGET_TYPE_USER_ACCOUNT
	TARGET_TYPE_TOPIC
	TARGET_TYPE_PACKAGE
)

const (
	DEFAULT_ALL     NotifyType = -1
	DEFAULT_SOUND   NotifyType = 1 << iota / 2 // 使用默认提示音提示
	DEFAULT_VIBRATE                            // 使用默认震动提示
	DEFAULT_LIGHTS                             // 使用默认led灯光提示
)

const (
	Msg      RequestType = "1"
	Feedback RequestType = "2"
	Emq      RequestType = "3"
)
const (
	RegId SubscribeType = iota + 1
	Alias
)

// '''
//     Union 并集
//     Intersection 交集
//     Except 差集
// '''
const (
	Union        BroadcastTopicOp = "UNION"
	Intersection BroadcastTopicOp = "INTERSECTION"
	Except       BroadcastTopicOp = "EXCEPT"
)

const (
	max_message_length = 140
	auto_switch_host   = true
	access_timeout     = 5000
	http_protocol      = "https"

	// '''
	//     相关Push域名定义
	// '''
	host_emq                 = "emq.xmpush.xiaomi.com"
	host_sandbox             = "sandbox.xmpush.xiaomi.com"
	host_production          = "api.xmpush.xiaomi.com"
	host_production_lg       = "lg.api.xmpush.xiaomi.com"
	host_production_c3       = "c3.api.xmpush.xiaomi.com"
	host_production_feedback = "feedback.xmpush.xiaomi.com"

	is_sandbox                   = false
	refresh_server_host_interval = 5 * 60
)

var (
	METHOD_MAP = map[HttpMethod]string{
		HTTP_GET:  "GET",
		HTTP_POST: "POST",
	}

	host = ""

	V2_SEND          = []string{"/v2/send"}
	V2_REGID_MESSAGE = []string{"/v2/message/regid"}
	V3_REGID_MESSAGE = []string{"/v3/message/regid"}

	V2_SUBSCRIBE_TOPIC            = []string{"/v2/topic/subscribe"}
	V2_UNSUBSCRIBE_TOPIC          = []string{"/v2/topic/unsubscribe"}
	V2_SUBSCRIBE_TOPIC_BY_ALIAS   = []string{"/v2/topic/subscribe/alias"}
	V2_UNSUBSCRIBE_TOPIC_BY_ALIAS = []string{"/v2/topic/unsubscribe/alias"}

	V2_ALIAS_MESSAGE = []string{"/v2/message/alias"}
	V3_ALIAS_MESSAGE = []string{"/v3/message/alias"}

	V2_BROADCAST_TO_ALL         = []string{"/v2/message/all"}
	V3_BROADCAST_TO_ALL         = []string{"/v3/message/all"}
	V2_BROADCAST                = []string{"/v2/message/topic"}
	V3_BROADCAST                = []string{"/v3/message/topic"}
	V2_MULTI_TOPIC_BROADCAST    = []string{"/v2/message/multi_topic"}
	V3_MILTI_TOPIC_BROADCAST    = []string{"/v3/message/multi_topic"}
	V2_DELETE_BROADCAST_MESSAGE = []string{"/v2/message/delete"}

	V2_USER_ACCOUNT_MESSAGE = []string{"/v2/message/user_account"}

	V2_SEND_MULTI_MESSAGE_WITH_REGID   = []string{"/v2/multi_messages/regids"}
	V2_SEND_MULTI_MESSAGE_WITH_ALIAS   = []string{"/v2/multi_messages/aliases"}
	V2_SEND_MULTI_MESSAGE_WITH_ACCOUNT = []string{"/v2/multi_messages/user_accounts"}

	V1_VALIDATE_REGID  = []string{"/v1/validation/regids"}
	V1_GET_ALL_ACCOUNT = []string{"/v1/account/all"}
	V1_GET_ALL_TOPIC   = []string{"/v1/topic/all"}
	V1_GET_ALL_ALIAS   = []string{"/v1/alias/all"}

	V1_MESSAGES_STATUS      = []string{"/v1/trace/messages/status"}
	V1_MESSAGE_STATUS       = []string{"/v1/trace/message/status"}
	V1_GET_MESSAGE_COUNTERS = []string{"/v1/stats/message/counters"}

	V1_FEEDBACK_INVALID_ALIAS = []interface{}{"/v1/feedback/fetch_invalid_aliases", Feedback}
	V1_FEEDBACK_INVALID_REGID = []interface{}{"/v1/feedback/fetch_invalid_regids", Feedback}

	V1_REGID_PRESENCE = []string{"/v1/regid/presence"}
	V2_REGID_PRESENCE = []string{"/v1/regid/presence"}

	V2_DELETE_SCHEDULE_JOB      = []string{"/v2/schedule_job/delete"}
	V3_DELETE_SCHEDULE_JOB      = []string{"/v3/schedule_job/delete"}
	V2_CHECK_SCHEDULE_JOB_EXIST = []string{"/v2/schedule_job/exist"}
	V2_QUERY_SCHEDULE_JOB       = []string{"/v2/schedule_job/query"}

	V1_EMQ_ACK_INFO      = []string{"/msg/ack/info", string(Emq)}
	V1_EMQ_CLICK_INFO    = []string{"/msg/click/info", string(Emq)}
	V1_EMQ_INVALID_REGID = []string{"/app/invalid/regid", string(Emq)}
)
