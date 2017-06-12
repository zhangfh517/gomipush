package gomipush
import (
    "context"
    "fmt"
    "net/url"
    // log "github.com/Sirupsen/logrus"
    "errors"

)

//android
//msg := gomipush.NewAndroidMessage(title, description , passThrough , restrictedPackageName)
//gomipush.NewClient("security").Send(msg).ToRegId([]string{"a","b"}).Do(ctx)
//gomipush.NewClient("security").Send(msg).ToMultiTopic([]string{"a", "b"}, Union).Do(ctx)
//ios
//msg := gomipush.NewIOSMessage(description).TimeToLive(110)
//gomipush.NewClient("security").Send(msg).ToRegId([]string{"a", "b"}).Do(ctx)

const (
    BROADCAST_TOPIC_MAX = 5
)
type SenderService struct {
    client *Client
    message Message
    retryTimes int
    targetUrl []string
}
func NewSenderService(c *Client, msg Message) (*SenderService, error){
    if msg == nil {
        return nil,errors.New("messge is nil")
    }
    return &SenderService{
        client: c,
        retryTimes: 3,
        message: msg,
    }, nil
}

// func (ss *SenderService) Message(msg Message) *SenderService{
//     ss.message = msg
//     return ss
// }
func (ss *SenderService) RetryTimes(retryTimes int) *SenderService{
    ss.retryTimes = retryTimes
    return ss
}
func (ss *SenderService) ToRegID(regId []string) *SenderService{
    ss.message.RegId(regId)
    ss.targetUrl = V3_REGID_MESSAGE
    return ss
}
func (ss *SenderService) ToAlias(alias []string) *SenderService{
    ss.message.Alias(alias)
    ss.targetUrl = V3_ALIAS_MESSAGE
    return ss
}
func (ss *SenderService) ToUserAccount(userAccount []string) *SenderService{
    ss.message.UserAccount(userAccount)
    ss.targetUrl = V2_USER_ACCOUNT_MESSAGE
    return ss
}
func (ss *SenderService) ToAll() *SenderService{
    packageName := ss.message.getRestrictedPackageName()
    ss.targetUrl = V2_BROADCAST_TO_ALL
    if len(packageName) > 1 {
        ss.targetUrl = V3_BROADCAST_TO_ALL
    }
    return ss
}
func (ss *SenderService) ToTopic(topic string) *SenderService{
    ss.message.Topic(topic)
    ss.targetUrl = V2_BROADCAST

    return ss
}
func (ss *SenderService) ToMultiTopic(topic []string, topicOp BroadcastTopicOp) (*SenderService, error){
    if len(topic) > BROADCAST_TOPIC_MAX {
        return nil, fmt.Errorf("topics more than max topic 5")
    }
    ss.message.MulitTopic(topic, topicOp)
    ss.targetUrl = V3_MILTI_TOPIC_BROADCAST
    return ss, nil
}

func (ss *SenderService) Do(ctx context.Context) (*Response, error){
    // log.Infof("Do: message= %v",ss.message)
    p, err:= ss.message.Source()
    if err != nil {
        return nil, err
    }
    return ss.client.PerformRequest(ctx, ss.targetUrl, ss.retryTimes, HTTP_POST, p.(url.Values), "")
}
