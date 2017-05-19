# gomipush

## Usage
### android client
#### First create client and message:  
```go
title := "gomipush sdk send to alias123"
description := "gomipush sdk send to alias123"
passThrough := 0
client := NewClient("security")
msg := NewAndroidMessage(title, description , passThrough , []string{restrictedPackageName})
```
#### Sencod do send  
* Send to topic:  
```go
  topic := "abc"
  rsp, err := client.Send(msg).ToTopic(topic).Do(context.Background())
```
* Send to multi topics
```go
  topic := []string{"abc","abc2"}
  rsp, err := client.Send(msg).ToMultiTopic(topic, op).Do(context.Background())
  //op is type BroadcastTopicOp ,and can must be one in [Union, Intersection, Except]
```
* Send to RegId
```go
  regId := []string{"abc"}
  rsp, err := client.Send(msg).ToToRegID(regId).Do(context.Background())
```
* Send to Alias
```go
  alias := []string{"abc"}
  rsp, err := client.Send(msg).ToAlias(alias).Do(context.Background())
```
* Send to userAccount
```go
userAccount := []string{"abc"}
rsp, err := client.Send(msg).ToUserAccount(userAccount).Do(context.Background())
```
* Send to All
```go
rsp, err := client.Send(msg).ToAll().Do(context.Background())
```
### ios client
## Important
大部分用例未经测试

