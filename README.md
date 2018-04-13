go-jpush
===================

概述
----------------------------------- 
   这是JPush REST API 的 go 版本封装开发包,仅支持最新的REST API v3功能。
   REST API 文档：http://docs.jpush.cn/display/dev/Push-API-v3
  

安装
---
   `go get -u github.com/moorper/go-jpush`
   
   
## Push API

### 1.构建要推送的平台： jpushclient.Platform
	//Platform
	var pf jpushclient.Platform
	pf.Add(jpushclient.ANDROID)
	pf.Add(jpushclient.IOS)
	pf.Add(jpushclient.WINPHONE)
	//pf.All()
      
### 2.构建接收听众： jpushclient.Audience
	//Audience
	var ad jpushclient.Audience
	s := []string{"t1", "t2", "t3"}
	ad.SetTag(s)
	id := []string{"1", "2", "3"}
	ad.SetID(id)
	//ad.All()
      
### 3.构建通知 jpushclient.Notice，或者消息： jpushclient.Message
      
	//Notice
	var notice jpushclient.Notice
	notice.SetAlert("alert_test")
	notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: "AndroidNotice"})
	notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: "IOSNotice"})
	notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: "WinPhoneNotice"})
      
    //jpushclient.Message
    var msg jpushclient.Message
	msg.Title = "Hello"
	msg.Content = "你是ylywn"
      
### 4.构建jpushclient.PayLoad
    payload := jpushclient.NewPushPayLoad()
	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetMessage(&msg)
	payload.SetNotice(&notice)
      
      
### 5.构建PushClient，发出推送
	var client jpushclient.Client
	client.AppKey = appKey
	client.MasterSecret = secret
	r, err := client.PushSend(*payload)
	if err != nil {
		fmt.Printf("err:%s", err.Error())
	} else {
		fmt.Printf("ok:%s", r)
	}

  
### 6.完整demo
    package main

	import (
		"fmt"
		"github.com/ylywyn/jpush-api-go-client"
	)

	const (
		appKey = "you jpush appkey"
		secret = "you jpush secret"
	)

	func main() {

		//Platform
		var pf jpushclient.Platform
		pf.Add(jpushclient.ANDROID)
		pf.Add(jpushclient.IOS)
		pf.Add(jpushclient.WINPHONE)
		//pf.All()

		//Audience
		var ad jpushclient.Audience
		s := []string{"1", "2", "3"}
		ad.SetTag(s)
		ad.SetAlias(s)
		ad.SetID(s)
		//ad.All()

		//Notice
		var notice jpushclient.Notice
		notice.SetAlert("alert_test")
		notice.SetAndroidNotice(&jpushclient.AndroidNotice{Alert: "AndroidNotice"})
		notice.SetIOSNotice(&jpushclient.IOSNotice{Alert: "IOSNotice"})
		notice.SetWinPhoneNotice(&jpushclient.WinPhoneNotice{Alert: "WinPhoneNotice"})

		var msg jpushclient.Message
		msg.Title = "Hello"
		msg.Content = "你是ylywn"

		payload := jpushclient.NewPushPayLoad()
		payload.SetPlatform(&pf)
		payload.SetAudience(&ad)
		payload.SetMessage(&msg)
		payload.SetNotice(&notice)

		bytes, _ := payload.ToBytes()
		fmt.Printf("%s\r\n", string(bytes))

		//push
		c := jpushclient.NewPushClient(secret, appKey)
		str, err := c.Send(bytes)
		if err != nil {
			fmt.Printf("err:%s", err.Error())
		} else {
			fmt.Printf("ok:%s", str)
		}
	}

## Device API

### 1. 设置 key 和 secret
`deviceClient := device.Client{AppKey: "", MasterSecret:""}`
### 2. 正常请求

```golang
registrationID := "aaaaaaaaaaaaa"
alias := "abc"

// 获取设备信息
deviceClient.GetDevicesInfo(registrationID)

// 为当前设备绑定信息
var tags = make(map[string][]string)
tags["add"] = []string{"a"}
deviceClient.UpdateDevicesInfo(registrationID, "", alias, "")

// 获取绑定了指定 alias 的设备，可指定 platform
deviceClient.FetchDevicesByAlias(alias, []string{})

// 删除指定 alias 的绑定，可指定 platform
deviceClient.AliasDelete(alias, []string{})

// 当前账号下所有的 tags
deviceClient.FetchTags()
```