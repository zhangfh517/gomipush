package gomipush

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
	// log "github.com/Sirupsen/logrus"
	"fmt"
)

type Message interface {
	Source() (interface{}, error)
	RegId(regId []string)                           //发送给一组设备，不同的registration_id之间用“,”分割
	Alias(alias []string)                           //可以提供多个alias，发送给一组设备，不同的alias之间用“,”分割。
	UserAccount(userAcct []string)                  //发送消息给设置了该user_account的所有设备。可以提供多个user_account，user_account之间用“,”分割。
	Topic(topic string)                             //发送消息给订阅了该topic的所有设备。
	MulitTopic(topic []string, op BroadcastTopicOp) //Stringtopicstopic列表，使用;$;分割。注: topics参数需要和topic_op参数配合使用，另外topic的数量不能超过5。UNION并集 INTERSECTION交集 EXCEPT差集
	getRestrictedPackageName() []string
}

type BaseMessage struct {
	regId                 []string
	alias                 []string
	userAccount           []string
	topic                 []string
	topicOp               BroadcastTopicOp
	restrictedPackageName []string
}

func (base *BaseMessage) RegId(regId []string) {
	base.regId = regId
}
func (base *BaseMessage) Alias(alias []string) {
	base.alias = alias
}
func (base *BaseMessage) UserAccount(userAcct []string) {
	base.userAccount = userAcct
}
func (base *BaseMessage) Topic(topic string) {
	base.topic = []string{topic}
}
func (base *BaseMessage) MulitTopic(topic []string, op BroadcastTopicOp) {
	base.topic = topic
	base.topicOp = op
}
func (base *BaseMessage) getRestrictedPackageName() []string {
	return base.restrictedPackageName
}

func (base *BaseMessage) Source() (interface{}, error) {
	params := url.Values{}
	if len(base.regId) > 0 {
		params.Set("registration_id", strings.Join(base.regId, ","))
		return params, nil
	}
	if len(base.alias) > 0 {
		params.Set("alias", strings.Join(base.alias, ","))
		return params, nil
	}
	if len(base.userAccount) > 0 {
		params.Set("user_account", strings.Join(base.userAccount, ","))
		return params, nil
	}
	if len(base.topic) == 1 {
		params.Set("topic", base.topic[0])
		return params, nil
	}
	if len(base.topic) > 1 {
		params.Set("topics", strings.Join(base.topic, ";$;"))
		if len(base.topicOp) > 0 {
			params.Set("topic_op", string(base.topicOp))
			return params, nil
		} else {
			return nil, errors.New("need topicOp")
		}
	}
	// return nil, errors.New("need target")
	return params, nil
}

type AndroidMessage struct {
	BaseMessage
	payload     string
	passThrough int //0 表示通知栏消息 1 表示透传消息
	title       string
	description string
	notifyType  NotifyType //DEFAULT_ALL = -1;
	//DEFAULT_SOUND  = 1;   // 使用默认提示音提示
	//DEFAULT_VIBRATE = 2;   // 使用默认震动提示
	//DEFAULT_LIGHTS = 4;    // 使用默认led灯光提示
	timeToLive int64 //可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	timeToSend int64 //定时消息，最大支持七天
	notifyId   int   //如果通知栏要显示多条推送消息，需要针对不同的消息设置不同的notify_id（相同notify_id的通知栏
	extra      AndroidExtra
}

// func NewAndroidMessage(title, description string, passThrough int, restrictedPackageName []string) *AndroidMessage {
// 	msg := &AndroidMessage{
// 		title : title,
// 		description : description,
// 		passThrough : passThrough,
// 	}
// 	msg.restrictedPackageName = restrictedPackageName
// 	return msg
// }

func NewAndroidMessage(title, description string) *AndroidMessage {
	msg := &AndroidMessage{
		passThrough: 1,
		title:       title,
		description: description,
	}
	return msg
}
func NewAndroidMessagePassThrough(payload string) *AndroidMessage {
	msg := &AndroidMessage{
		passThrough: 1,
		payload:     payload,
	}
	return msg
}

type AndroidExtra struct {
	ticker           string //开启通知消息在状态栏滚动显示
	notifyForeground int    //开启/关闭app在前台时的通知弹出。当extra.notify_foreground值为”1″时，app会弹出通知栏消息；当extra.notify_foreground值为”0″时，app不会弹出通知栏消息。注：默认情况下会弹出通知栏消息。
	notifyEffect     int    //“1″：通知栏点击后打开app的Launcher Activity。
	//“2″：通知栏点击后打开app的任一Activity（开发者还需要传入extra.intent_uri）。
	//“3″：通知栏点击后打开网页（开发者还需要传入extra.web_uri）。
	intentUri       string //打开一个app组件
	webUri          string //打开一个网页
	flowControl     int    //平滑推送的速度
	layoutName      string //
	layoutValue     string
	jobkey          string //推送批次，聚合消息
	callback        string
	locale          []string //可以接收消息的设备的语言范围，用逗号分隔
	localeNotIn     []string //无法收到消息的设备的语言范围，逗号分隔。
	model           []string //1. 可以收到消息的设备的机型范围，逗号分隔,2.以收到消息的设备的品牌范围，逗号分割。3.可以收到消息的设备的价格范围，逗号分隔。
	modelNotIn      []string //1.无法收到消息的设备的机型范围，逗号分隔
	appVersion      []string //可以接收消息的app版本号，用逗号分割。安卓app版本号来源于manifest文件中的”android:versionName”的值
	appVersionNotIn []string //无法接收消息的app版本号，用逗号分割。
	connpt          string   //指定在特定的网络环境下才能接收到消息。目前仅支持指定”wifi”。
}

func (ad *AndroidMessage) Source() (interface{}, error) {
	var rq = url.Values{}

	baseRq, err := ad.BaseMessage.Source()
	if err != nil {
		return nil, err
	}
	for k, v := range baseRq.(url.Values) {
		// log.Infof("use base Source: k %s, v %s",k,v[0])
		rq.Set(k, v[0])
	}
	if ad.payload != "" {
		rq.Set("payload", ad.payload)
	}
	rq.Set("pass_through", strconv.Itoa(ad.passThrough))

	if ad.title != "" {
		rq.Set("title", ad.title)
	}
	if ad.description != "" {
		rq.Set("description", ad.description)
	}
	if ad.notifyType != 0 {
		rq.Set("notify_type", fmt.Sprint(ad.notifyType))
	}
	if ad.timeToLive > 0 {
		rq.Set("time_to_live", strconv.FormatInt(ad.timeToLive, 10))
	}
	if ad.timeToSend > 0 {
		rq.Set("time_to_send", strconv.FormatInt(ad.timeToSend, 10))
	}
	if ad.notifyId != 0 {
		rq.Set("notify_id", strconv.Itoa(ad.notifyId))
	}
	if ad.extra.ticker != "" {
		rq.Set("extra.ticker", ad.extra.ticker)
	}
	if ad.extra.notifyForeground != 0 {
		rq.Set("extra.notify_foreground", strconv.Itoa(ad.extra.notifyForeground))
	}
	if ad.extra.notifyEffect != 0 {
		rq.Set("extra.notify_effect", strconv.Itoa(ad.extra.notifyEffect))
	}
	if ad.extra.intentUri != "" {
		rq.Set("extra.intent_uri", ad.extra.intentUri)
	}
	if ad.extra.webUri != "" {
		rq.Set("extra.web_uri", ad.extra.webUri)
	}
	if ad.extra.flowControl != 0 {
		rq.Set("extra.flow_control", strconv.Itoa(ad.extra.flowControl))
	}
	if ad.extra.layoutName != "" {
		rq.Set("extra.layout_name", ad.extra.layoutName)
	}
	if ad.extra.layoutValue != "" {
		rq.Set("extra.layout_value", ad.extra.layoutValue)
	}
	if ad.extra.jobkey != "" {
		rq.Set("extra.jobkey", ad.extra.jobkey)
	}
	if ad.extra.callback != "" {
		rq.Set("extra.callback", ad.extra.callback)
	}
	if len(ad.extra.locale) != 0 {
		rq.Set("extra.locale", strings.Join(ad.extra.locale, ","))
	}
	if len(ad.extra.localeNotIn) != 0 {
		rq.Set("extra.locale_not_in", strings.Join(ad.extra.localeNotIn, ","))
	}
	if len(ad.extra.model) != 0 {
		rq.Set("extra.model", strings.Join(ad.extra.model, ","))
	}
	if len(ad.extra.modelNotIn) != 0 {
		rq.Set("extra.model_not_in", strings.Join(ad.extra.modelNotIn, ","))
	}
	if len(ad.extra.appVersion) != 0 {
		rq.Set("extra.app_version", strings.Join(ad.extra.appVersion, ","))
	}
	if len(ad.extra.appVersionNotIn) != 0 {
		rq.Set("extra.app_version_not_in", strings.Join(ad.extra.appVersionNotIn, ","))
	}
	if ad.extra.connpt != "" {
		rq.Set("extra.connpt", ad.extra.connpt)
	}
	if len(ad.restrictedPackageName) > 0 {
		rq.Set("restricted_package_name", strings.Join(ad.restrictedPackageName, ","))
	}

	return rq, nil
}

func (ad *AndroidMessage) Payload(payload string) *AndroidMessage {
	ad.payload = payload
	return ad
}

func (ad *AndroidMessage) PassThrough(passThrough int) *AndroidMessage {
	ad.passThrough = passThrough
	return ad
}

func (ad *AndroidMessage) Title(title string) *AndroidMessage {
	ad.title = title
	return ad
}
func (ad *AndroidMessage) Description(description string) *AndroidMessage {
	ad.description = description
	return ad
}
func (ad *AndroidMessage) NotifyType(notifyType NotifyType) *AndroidMessage {
	ad.notifyType = notifyType
	return ad
}
func (ad *AndroidMessage) TimeToLive(timeToLive int64) *AndroidMessage {
	ad.timeToLive = timeToLive
	return ad
}
func (ad *AndroidMessage) TimeToSend(timeToSend int64) *AndroidMessage {
	ad.timeToSend = timeToSend
	return ad
}
func (ad *AndroidMessage) RestrictedPackageName(packages []string) *AndroidMessage {
	ad.restrictedPackageName = packages
	return ad
}

func (ad *AndroidMessage) NotifyId(notifyId int) *AndroidMessage {
	ad.notifyId = notifyId
	return ad
}
func (ad *AndroidMessage) Extraticker(ticker string) *AndroidMessage {
	ad.extra.ticker = ticker
	return ad
}
func (ad *AndroidMessage) ExtraNotifyForeground(notifyForeground int) *AndroidMessage {
	ad.extra.notifyForeground = notifyForeground
	return ad
}
func (ad *AndroidMessage) ExtraNotifyEffect(notifyEffect int) *AndroidMessage {
	ad.extra.notifyEffect = notifyEffect
	return ad
}
func (ad *AndroidMessage) ExtraIntentUri(intentUri string) *AndroidMessage {
	ad.extra.intentUri = intentUri
	return ad
}
func (ad *AndroidMessage) ExtraWebUri(webUri string) *AndroidMessage {
	ad.extra.webUri = webUri
	return ad
}
func (ad *AndroidMessage) ExtraFlowControl(flowControl int) *AndroidMessage {
	ad.extra.flowControl = flowControl
	return ad
}
func (ad *AndroidMessage) ExtraLayoutName(layoutName string) *AndroidMessage {
	ad.extra.layoutName = layoutName
	return ad
}
func (ad *AndroidMessage) ExtraLayoutValue(layoutValue string) *AndroidMessage {
	ad.extra.layoutValue = layoutValue
	return ad
}
func (ad *AndroidMessage) ExtraJobkey(jobkey string) *AndroidMessage {
	ad.extra.jobkey = jobkey
	return ad
}
func (ad *AndroidMessage) ExtraCallback(callback string) *AndroidMessage {
	ad.extra.callback = callback
	return ad
}
func (ad *AndroidMessage) ExtraLocale(locale []string) *AndroidMessage {
	ad.extra.locale = locale
	return ad
}
func (ad *AndroidMessage) ExtraLocaleNotIn(localeNotIn []string) *AndroidMessage {
	ad.extra.localeNotIn = localeNotIn
	return ad
}
func (ad *AndroidMessage) ExtraModel(model []string) *AndroidMessage {
	ad.extra.model = model
	return ad
}
func (ad *AndroidMessage) ExtraModelNotIn(modelNotIn []string) *AndroidMessage {
	ad.extra.modelNotIn = modelNotIn
	return ad
}
func (ad *AndroidMessage) ExtraAppVersion(appVersion []string) *AndroidMessage {
	ad.extra.appVersion = appVersion
	return ad
}
func (ad *AndroidMessage) ExtraAppVersionNotIn(appVersionNotIn []string) *AndroidMessage {
	ad.extra.appVersionNotIn = appVersionNotIn
	return ad
}
func (ad *AndroidMessage) ExtraConnpt(connpt string) *AndroidMessage {
	ad.extra.connpt = connpt
	return ad
}

type IOSMessage struct {
	BaseMessage
	description     string //通知栏展示的通知的描述。
	apsProperFields IOSApsProperField
	timeToLive      int64 //可选项。如果用户离线，设置消息在服务器保存的时间，单位：ms。服务器默认最长保留两周。
	timeToSend      int64 //可选项。定时发送消息。用自1970年1月1日以来00:00:00.0 UTC时间表示（以毫秒为单位的时间）。注：仅支持七天内的定时消息。
	extra           IOSExtra
}

func NewIOSMessage(description string) *IOSMessage {
	return &IOSMessage{
		description: description,
	}
}

type IOSApsProperField struct {
	title          string //	在通知栏展示的通知的标题（支持iOS10及以上版本，如有该字段，会覆盖掉description字段）。
	subtitle       string //	展示在标题下方的子标题（支持iOS10及以上版本，如有该字段，会覆盖掉description字段）。
	body           string //	在通知栏展示的通知的内容（支持iOS10及以上版本，如有该字段，会覆盖掉description字段）。
	mutableContent string //通知可以修改选项，设置之后，在展示远程通知之前会进入Notification Service Extension中允许程序对通知内容修改（支持iOS10及以上版本）。

}
type IOSExtra struct {
	soundUrl string //	可选项，自定义消息铃声。当值为空时为无声，default为系统默认声音。
	badge    int    //可选项。通知角标。
	category string //可选项。iOS8推送消息快速回复类别。
}

func (ad *IOSMessage) Source() (interface{}, error) {
	rq := url.Values{}
	baseRq, err := ad.BaseMessage.Source()
	if err != nil {
		return nil, err
	}
	for k, v := range baseRq.(url.Values) {
		// log.Infof("ios use base Source: k %s, v %s",k,v[0])
		rq.Set(k, v[0])
	}
	if ad.description != "" {
		rq.Set("description", ad.description)
	}
	if ad.timeToLive > 0 {
		rq.Set("time_to_live", strconv.FormatInt(ad.timeToLive, 10))
	}
	if ad.timeToSend > 0 {
		rq.Set("time_to_send", strconv.FormatInt(ad.timeToSend, 10))
	}
	if ad.apsProperFields.title != "" {
		rq.Set("aps_proper_fields.title", ad.apsProperFields.title)
	}
	if ad.apsProperFields.subtitle != "" {
		rq.Set("aps_proper_fields.subtitle", ad.apsProperFields.subtitle)
	}
	if ad.apsProperFields.body != "" {
		rq.Set("aps_proper_fields.body", ad.apsProperFields.body)
	}
	if ad.apsProperFields.mutableContent != "" {
		rq.Set("aps_proper_fields.mutable-content", ad.apsProperFields.mutableContent)
	}
	if ad.extra.soundUrl != "" {
		rq.Set("extra.sound_url", ad.extra.soundUrl)
	}
	if ad.extra.badge != 0 {
		rq.Set("extra.badge", strconv.Itoa(ad.extra.badge))
	}
	if ad.extra.category != "" {
		rq.Set("extra.category", ad.extra.category)
	}
	return rq, nil

}

func (ios *IOSMessage) Description(description string) *IOSMessage {
	ios.description = description
	return ios
}
func (ios *IOSMessage) TimeToLive(timeToLive int64) *IOSMessage {
	ios.timeToLive = timeToLive
	return ios
}
func (ios *IOSMessage) TimeToSend(timeToSend int64) *IOSMessage {
	ios.timeToSend = timeToSend
	return ios
}

func (ios *IOSMessage) ApsTitle(title string) *IOSMessage {
	ios.apsProperFields.title = title
	return ios
}
func (ios *IOSMessage) ApsSubtitle(subtitle string) *IOSMessage {
	ios.apsProperFields.subtitle = subtitle
	return ios
}
func (ios *IOSMessage) ApsBody(body string) *IOSMessage {
	ios.apsProperFields.body = body
	return ios
}
func (ios *IOSMessage) ApsMutableContent(mutableContent string) *IOSMessage {
	ios.apsProperFields.mutableContent = mutableContent
	return ios
}
func (ios *IOSMessage) ExtraSoundUrl(soundUrl string) *IOSMessage {
	ios.extra.soundUrl = soundUrl
	return ios
}
func (ios *IOSMessage) ExtraBadge(badge int) *IOSMessage {
	ios.extra.badge = badge
	return ios
}
func (ios *IOSMessage) ExtraCategory(category string) *IOSMessage {
	ios.extra.category = category
	return ios
}
