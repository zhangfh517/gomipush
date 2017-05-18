package gomipush

import (
"net/url"
"strings"
"errors"
"strconv"
log "github.com/Sirupsen/logrus"

)

type Message interface {
	Source() (interface{}, error)
	RegId(regId []string)
	Alias(alias []string)
	UserAccount(userAcct []string)
	Topic(topic string)
	MulitTopic(topic []string, op BroadcastTopicOp)
	getRestrictedPackageName() []string
}

type BaseMessage struct {
	regId []string
	alias []string
	userAccount []string
	topic []string
	topicOp BroadcastTopicOp
	restrictedPackageName []string

}

func (base *BaseMessage) RegId(regId []string){
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
func (base *BaseMessage) getRestrictedPackageName() []string{
	return base.restrictedPackageName
}

func (base *BaseMessage) Source() (interface{}, error) {
	params := url.Values{}
	if len(base.regId) > 0 {
		params.Set("registration_id",strings.Join(base.regId,  ","))
		return params, nil
	}
	if len(base.alias) > 0 {
		params.Set("alias", strings.Join(base.alias,  ","))
		return params, nil
	}
	if len(base.userAccount) > 0 {
		params.Set("user_account", strings.Join(base.userAccount,  ","))
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
		}else {
			return nil, errors.New("need topicOp")
		}
	}
	return nil, errors.New("need target")
}


type AndroidMessage struct {
	BaseMessage
	payload string
	passThrough int
	title string
	description string
	notifyType int
	timeToLive int64
	timeToSend int64
	notifyId int
	extra AndroidExtra
}
func NewAndroidMessage(title, description string, passThrough int, restrictedPackageName []string) *AndroidMessage {
	msg := &AndroidMessage{
		title : title,
		description : description,
		passThrough : passThrough,
	}
	msg.restrictedPackageName = restrictedPackageName
	return msg
}

type AndroidExtra struct {
	tricker string
	notifyForeground string
	notifyEffect string
	intentUri string
	webUri string
	flowControl int
	layoutName int
	layoutValue int
	jobkey string
	callback string
	locale string
	localeNotIn string
	model string
	modelNotIn string
	appVersion string
	appVersionNotIn string
	connpt string
}

func (ad *AndroidMessage) Source() (interface{}, error) {
	var rq = url.Values{}

	baseRq, err := ad.BaseMessage.Source()
	if err != nil {
		return nil, err
	}
	for k, v := range baseRq.(url.Values) {
		log.Infof("use base Source: k %s, v %s",k,v[0])
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
	if ad.notifyType != -1 {
		rq.Set("notify_type", strconv.Itoa(ad.notifyType))
	}
	if ad.timeToLive > 0 {
		rq.Set("time_to_live", strconv.FormatInt(ad.timeToLive, 10))
	}
	if ad.timeToSend > 0 {
		rq.Set("time_to_Send", strconv.FormatInt(ad.timeToSend, 10))
	}
	if ad.notifyId != -1 {
		rq.Set("notify_id", strconv.Itoa(ad.notifyId))
	}
	if ad.extra.tricker != "" {
		rq.Set("extra.tricker", ad.extra.tricker)
	}
	if ad.extra.notifyForeground != "" {
		rq.Set("extra.notify_foreground" ,ad.extra.notifyForeground)
	}
	if ad.extra.notifyEffect != "" {
		rq.Set("extra.notify_effect", ad.extra.notifyEffect)
	}
	if ad.extra.intentUri != "" {
		rq.Set("extra.intent_uri", ad.extra.intentUri)
	}
	if ad.extra.webUri != "" {
		rq.Set("extra.web_uri", ad.extra.webUri)
	}
	if ad.extra.flowControl != -1 {
		rq.Set("extra.flow_control", strconv.Itoa(ad.extra.flowControl))
	}
	if ad.extra.layoutName != -1 {
		rq.Set("extra.layout_name", strconv.Itoa(ad.extra.layoutName))
	}
	if ad.extra.layoutValue != -1 {
		rq.Set("extra.layout_value", strconv.Itoa(ad.extra.layoutValue))
	}
	if ad.extra.jobkey != "" {
		rq.Set("extra.jobkey", ad.extra.jobkey)
	}
	if ad.extra.callback != "" {
		rq.Set("extra.callback", ad.extra.callback)
	}
	if ad.extra.locale != "" {
		rq.Set("extra.locale", ad.extra.locale)
	}
	if ad.extra.localeNotIn != "" {
		rq.Set("extra.locale_not_in", ad.extra.localeNotIn)
	}
	if ad.extra.model != "" {
		rq.Set("extra.model",ad.extra.model)
	}
	if ad.extra.modelNotIn != "" {
		rq.Set("extra.model_not_in",ad.extra.modelNotIn)
	}
	if ad.extra.appVersion != "" {
		rq.Set("extra.app_version",ad.extra.appVersion)
	}
	if ad.extra.appVersionNotIn != "" {
		rq.Set("extra.app_version_not_in",ad.extra.appVersionNotIn)
	}
	if ad.extra.connpt != "" {
		rq.Set("extra.connpt",ad.extra.connpt)
	}
	if len(ad.restrictedPackageName) > 0 {
		rq.Set("restricted_package_name", strings.Join(ad.restrictedPackageName, ","))
	}


	// ad.Source()
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
func (ad *AndroidMessage) NotifyType(notifyType int) *AndroidMessage {
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
func (ad *AndroidMessage) NotifyId(notifyId int) *AndroidMessage {
	ad.notifyId = notifyId
	return ad
}
func (ad *AndroidMessage) Extratricker(tricker string) *AndroidMessage {
	ad.extra.tricker = tricker
	return ad
}
func (ad *AndroidMessage) ExtraNotifyForeground(notifyForeground string) *AndroidMessage {
	ad.extra.notifyForeground = notifyForeground
	return ad
}
func (ad *AndroidMessage) ExtraNotifyEffect(notifyEffect string) *AndroidMessage {
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
func (ad *AndroidMessage) ExtraLayoutName(layoutName int) *AndroidMessage {
	ad.extra.layoutName = layoutName
	return ad
}
func (ad *AndroidMessage) ExtraLayoutValue(layoutValue int) *AndroidMessage {
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
func (ad *AndroidMessage) ExtraLocale(locale string) *AndroidMessage {
	ad.extra.locale = locale
	return ad
}
func (ad *AndroidMessage) ExtraLocaleNotIn(localeNotIn string) *AndroidMessage {
	ad.extra.localeNotIn = localeNotIn
	return ad
}
func (ad *AndroidMessage) ExtraModel(model string) *AndroidMessage {
	ad.extra.model = model
	return ad
}
func (ad *AndroidMessage) ExtraModelNotIn(modelNotIn string) *AndroidMessage {
	ad.extra.modelNotIn = modelNotIn
	return ad
}
func (ad *AndroidMessage) ExtraAppVersion(appVersion string) *AndroidMessage {
	ad.extra.appVersion = appVersion
	return ad
}
func (ad *AndroidMessage) ExtraAppVersionNotIn(appVersionNotIn string) *AndroidMessage {
	ad.extra.appVersionNotIn = appVersionNotIn
	return ad
}
func (ad *AndroidMessage) ExtraConnpt(connpt string) *AndroidMessage {
	ad.extra.connpt = connpt
	return ad
}


type IOSMessage struct {
	BaseMessage
	description string
	apsProperFields IOSApsProperField
	timeToLive int64
	timeToSend int64
	extra IOSExtra
}

func NewIOSMessage(description string) *IOSMessage{
	return &IOSMessage{
		description : description,
	}
}

type IOSApsProperField struct {
	title string
	subtitle string
	body string
	mutableContent string
}
type IOSExtra struct {
	soundUrl string
	badge string
	category string
}
func (ad *IOSMessage) Source() (interface{}, error) {
	rq := url.Values{}

	if ad.description != "" {
		rq.Set("description", ad.description)
	}
	if ad.timeToLive != -1 {
		rq.Set("time_to_live", strconv.FormatInt(ad.timeToLive, 10))
	}
	if ad.timeToSend != -1 {
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
	if ad.extra.badge != "" {
		rq.Set("extra.badge", ad.extra.badge)
	}
	if ad.extra.category != "" {
		rq.Set("extra.category", ad.extra.category)
	}
	return rq, nil

}

func (ios *IOSMessage) Description(description string) *IOSMessage{
	ios.description = description
	return ios
}
func (ios *IOSMessage) TimeToLive(timeToLive int64) *IOSMessage{
	ios.timeToLive = timeToLive
	return ios
}
func (ios *IOSMessage) TimeToSend(timeToSend int64) *IOSMessage{
	ios.timeToSend = timeToSend
	return ios
}

func (ios *IOSMessage) ApsTitle(title string) *IOSMessage{
	ios.apsProperFields.title = title
	return ios
}
func (ios *IOSMessage) ApsSubtitle(subtitle string) *IOSMessage{
	ios.apsProperFields.subtitle = subtitle
	return ios
}
func (ios *IOSMessage) ApsBody(body string) *IOSMessage{
	ios.apsProperFields.body = body
	return ios
}
func (ios *IOSMessage) ApsMutableContent(mutableContent string) *IOSMessage{
	ios.apsProperFields.mutableContent = mutableContent
	return ios
}
func (ios *IOSMessage) ExtraSoundUrl(soundUrl string) *IOSMessage{
	ios.extra.soundUrl = soundUrl
	return ios
}
func (ios *IOSMessage) ExtraBadge(badge string) *IOSMessage{
	ios.extra.badge = badge
	return ios
}
func (ios *IOSMessage) ExtraCategory(category string) *IOSMessage{
	ios.extra.category = category
	return ios
}


