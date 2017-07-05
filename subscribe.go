package gomipush

import (
	"context"
	"net/url"
	"strings"
	// "errors"
	log "github.com/Sirupsen/logrus"
)

//topic := gomipush.NewSubscribedTopic("topic").RestrictedPackageName([]string{"android"})

//resp, err := gomipush.NewClient("security").Subscribe(topic).SubscribeRegIds([]string{"regids"}).Do(context.Backgroud())

type SubscribedTopic struct {
	topic                 string
	regId                 []string
	alias                 []string
	category              string
	restrictedPackageName []string
}

func NewSubscribedTopic(topic string) *SubscribedTopic {
	return &SubscribedTopic{
		topic: topic,
	}
}

func (st *SubscribedTopic) RegId(regId []string) {
	st.regId = regId
	// return st
}
func (st *SubscribedTopic) Alias(alias []string) {
	st.alias = alias
	// return st
}
func (st *SubscribedTopic) Category(category string) *SubscribedTopic {
	st.category = category
	return st
}
func (st *SubscribedTopic) RestrictedPackageName(packageName []string) *SubscribedTopic {
	st.restrictedPackageName = packageName
	return st
}

func (st *SubscribedTopic) Source() interface{} {
	params := url.Values{}
	if len(st.topic) > 0 {
		params.Set("topic", st.topic)
	}
	if len(st.regId) > 0 {
		params.Set("registration_id", strings.Join(st.regId, ","))
	}
	if len(st.alias) > 0 {
		//不确定是alias 还是 aliases
		params.Set("aliases", strings.Join(st.alias, ","))
	}
	if len(st.category) > 0 {
		params.Set("category", st.category)
	}
	if len(st.restrictedPackageName) > 0 {
		params.Set("restricted_package_name", strings.Join(st.restrictedPackageName, ","))
	}
	log.Infof("subscribe source : %v", params)
	return params

}

type SubscribeService struct {
	client     *Client
	topic      SubscribedTopic
	retryTimes int
	targetUrl  []string
}

func NewSubscribeService(c *Client, topic SubscribedTopic) *SubscribeService {
	return &SubscribeService{
		client:     c,
		retryTimes: 3,
		topic:      topic,
	}
}

func (ss *SubscribeService) SubscribeRegIds(regId []string) *SubscribeService {
	ss.targetUrl = V2_SUBSCRIBE_TOPIC
	ss.topic.RegId(regId)
	return ss
}

func (ss *SubscribeService) UnsubscribeRegIds(regId []string) *SubscribeService {
	ss.targetUrl = V2_UNSUBSCRIBE_TOPIC
	ss.topic.RegId(regId)
	return ss
}

func (ss *SubscribeService) SubscribeAlias(alias []string) *SubscribeService {
	ss.targetUrl = V2_SUBSCRIBE_TOPIC_BY_ALIAS
	ss.topic.Alias(alias)
	return ss
}

func (ss *SubscribeService) UnsubscribeAlias(alias []string) *SubscribeService {
	ss.targetUrl = V2_UNSUBSCRIBE_TOPIC_BY_ALIAS
	ss.topic.Alias(alias)
	return ss
}

func (ss *SubscribeService) Do(ctx context.Context) (*Response, error) {
	p := ss.topic.Source()
	return ss.client.PerformRequest(ctx, ss.targetUrl, ss.retryTimes, HTTP_POST, p.(url.Values), "")
}
