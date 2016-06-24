package main

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"jzlservice/smssender"
	"jzlservice/snsscheduler"
	"jzlservice/weibosender"
	"jzlservice/weixinsender"
)

type SMSSenderClient struct {
	smsProvider string
}

func (this *SMSSenderClient) Send(msg *snsscheduler.SNSMessage) (r []*smssender.SMSStatus, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("短信消息[%v]开始发送中...", msg.TaskId)

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = "没有提供短信服务提供商"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败")
		return nil, err
	}

	var value []*smssender.SMSEntry
	for _, sns_entry := range msg.Entries {
		value = append(value, &smssender.SMSEntry{
			TaskId:             msg.TaskId,
			SerialNumber:       "1",
			Content:            sns_entry.Content,
			Receiver:           sns_entry.Receiver,
			Signature:          sns_entry.SmsExtraInfo.Signature,
			ServiceMinorNumber: sns_entry.SmsExtraInfo.ServiceMinorNumber,
			Category:           int8(sns_entry.SmsExtraInfo.Category),
		})
	}

	r, err = client.SendSMS(value)
	if err != nil {
		LOG_ERROR("短信消息[%v]发送失败. 失败原因：%s", msg.TaskId, err)
		return nil, err
	}

	LOG_INFO("短信消息[%v]发送成功. 发送状态码：%v", msg.TaskId, r)

	return r, nil
}

func (this *SMSSenderClient) SendMessage(task_id int64, serial_number string, receiver string, content string, signature string, service_number string, catetory int8) (r *smssender.SMSStatus, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("短信消息[%v]开始发送中...", content)

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = "没有提供短信服务提供商"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败, 失败原因：%v", networkAddr, err)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败，失败原因：%v", err)
		return nil, err
	}

	entry := &smssender.SMSEntry{
		TaskId:             task_id,
		SerialNumber:       serial_number,
		Receiver:           receiver,
		Content:            content,
		Signature:          signature,
		ServiceMinorNumber: service_number,
		Category:           catetory,
	}
	r, err = client.SendMessage(entry)
	if err != nil {
		LOG_ERROR("短信消息[%v]发送失败. 失败原因：%s", content, err)
		return nil, err
	}

	LOG_INFO("短信消息[%v]发送成功. 发送状态码：%v", content, r)

	return r, nil
}

func (this *SMSSenderClient) GetReport(category int16) (r *smssender.SMSReport, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取发送短信的状态报告开始...")

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = "没有提供短信服务提供商"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败, 失败原因: %v", networkAddr, err)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败, 失败原因：%v", err)
		return nil, err
	}

	r, err = client.GetReport(category)
	if err != nil {
		LOG_ERROR("获取发送短信的状态报告失败。 失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取发送短信的状态报告成功。状态报告为：%v", r)

	return r, nil
}

func (this *SMSSenderClient) GetMOMessage(category int16) (r *smssender.SMSMOMessage, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取用户回复的短信开始...")

	if this.smsProvider != "" {
		addr, addrIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".addr")
		port, portIsSet = g_config.Get("sms_sender_" + this.smsProvider + ".port")
	} else {
		outputStr = "没有提供短信服务提供商"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址设置错误", this.smsProvider)
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = fmt.Sprintf("短信服务提供商[%v]的网络连接地址没有设置", this.smsProvider)
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到短信发送服务[%v]的连接失败，失败原因: %v", networkAddr, err)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := smssender.NewSMSSenderClientFactory(trans, protocolFactory)
	if err := trans.Open(); err != nil {
		LOG_ERROR("打开到短信发送服务的连接失败, 失败原因：%v", err)
		return nil, err
	}

	r, err = client.GetMOMessage(category)
	if err != nil {
		LOG_ERROR("获取用户回复的短信失败。 失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取上行短信成功。上行短信为：%v", r)

	return r, nil
}

type WeiboSenderClient struct {
}

func (this *WeiboSenderClient) Send(msg *snsscheduler.SNSMessage) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("微博消息[%v]开始发送中...", msg.TaskId)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	for _, sns_entry := range msg.Entries {
		weibo_status := &weibosender.WeiboStatus{
			TaskId:      msg.TaskId,
			AccessToken: sns_entry.Sender,
			Status:      sns_entry.Content,
			Visible:     sns_entry.WeiboExtraInfo.Visible,
			ListId:      sns_entry.WeiboExtraInfo.ListId,
			Latitude:    sns_entry.WeiboExtraInfo.Latitude,
			Longitude:   sns_entry.WeiboExtraInfo.Longitude,
			Annotations: sns_entry.WeiboExtraInfo.Annotations,
			RealIp:      sns_entry.WeiboExtraInfo.RealIp,
			Pic:         sns_entry.WeiboExtraInfo.Pic,
		}

		r, err = client.SendStatus(weibo_status)
		if err != nil {
			LOG_ERROR("微博消息[%v]发送失败. 失败原因：%v", msg.TaskId, err)
			return "", err
		}
	}

	LOG_INFO("微博消息[%v]发送成功.", msg.TaskId)

	return r, nil
}

func (this *WeiboSenderClient) SendMessage(status *weibosender.WeiboStatus) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("微博消息[%v]开始发送中...", status.TaskId)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败，失败原因：%v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.SendStatus(status)
	if err != nil {
		LOG_ERROR("微博消息[%v]发送失败. 失败原因：%v", status.TaskId, err)
		return "", err
	}

	LOG_INFO("微博消息[%v]发送成功.", status.TaskId)

	return r, nil
}

func (this *WeiboSenderClient) GetUserInfoById(access_token string, uid int64) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取用户[%v]的详细信息开始...", uid)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserInfoById(access_token, uid)
	if err != nil {
		LOG_ERROR("获取用户[%v]的详细信息失败. 失败原因：%v", uid, err)
		return "", err
	}

	LOG_INFO("获取用户[%v]的详细信息成功. 返回结果：%v", uid, r)

	return r, nil
}

func (this *WeiboSenderClient) GetUserInfoByName(access_token string, screen_name string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取用户[%v]的详细信息开始...", screen_name)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败，失败原因：%v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserInfoByName(access_token, screen_name)
	if err != nil {
		LOG_ERROR("获取用户[%v]的详细信息失败. 失败原因：%s", screen_name, err)
		return "", err
	}

	LOG_INFO("获取用户[%v]的详细信息成功. 返回结果：%v", screen_name, r)

	return r, nil
}

func (this *WeiboSenderClient) SendPrivateMsg(access_token string, type_a1 string, data string, receiver_id int64, save_sender_box int32) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("发送用户[%v]的私信开始. 私信类型：%v, 私信内容：%v", receiver_id, type_a1, data)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败, 失败原因: %v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.SendMessage(access_token, type_a1, data, receiver_id, save_sender_box)
	if err != nil {
		LOG_ERROR("发送用户[%v]的私信失败，失败原因：%s", receiver_id, err)
		return "", err
	}

	LOG_INFO("发送用户[%v]的私信成功. 返回结果：%v", receiver_id, r)

	return r, nil
}

func (this *WeiboSenderClient) GetStatuses(access_token string, since_id int64, max_id int64, count int32, page int32) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取最新发布的微博列表开始.")

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败, 失败原因：%v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetStatuses(access_token, since_id, max_id, count, page)
	if err != nil {
		LOG_ERROR("获取最近发布的微博列表失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("获取最近发布的微博列表成功. 返回结果：%v", r)

	return r, nil
}

func (this *WeiboSenderClient) GetConcernStatuses(access_token string, since_id int64, max_id int64, count int32, page int32) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取用户[%v]关注用户最新发布的微博列表开始.", access_token)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败, 失败原因: %v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetConcernStatuses(access_token, since_id, max_id, count, page)
	if err != nil {
		LOG_ERROR("获取用户[%v]关注用户的最近发布的微博列表失败，失败原因：%v", access_token, err)
		return "", err
	}

	LOG_INFO("获取用户[%v]关注用户的最近发布的微博列表成功. 返回结果：%v", access_token, r)

	return r, nil
}

func (this *WeiboSenderClient) GetStatusInteractCount(access_token string, ids string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微博[%v]互动数开始.", ids)

	addr, addrIsSet = g_config.Get("weibo_sender.addr")
	port, portIsSet = g_config.Get("weibo_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "微博发送服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "微博发送服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到微博发送服务[%v]的连接失败, 失败原因: %v", networkAddr, err)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weibosender.NewWeiboSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到微博发送服务WeiboSender的连接失败, 微博发送服务WeiboSender可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetStatusInteractCount(access_token, ids)
	if err != nil {
		LOG_ERROR("获取微博[%v]互动数失败，失败原因：%v", ids, err)
		return "", err
	}

	LOG_INFO("获取微博[%v]的互动数成功. 返回结果：%v", ids, r)

	return r, nil
}

type WeixinSenderClient struct {
}

func (this *WeixinSenderClient) GetAccessToken(appid string, appsecret string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取公众号[%v]的AccessToken开始.", appid)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetAccessToken(appid, appsecret)
	if err != nil {
		LOG_ERROR("获取公众号[%v]的AccessToken失败，失败原因：%v", appid, err)
		return "", err
	}

	LOG_INFO("获取公众号[%v]的AccessToken成功", appid)

	return

}

func (this *WeixinSenderClient) AddKFAccount(access_token string, kfaccount string, nickname string, password string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("添加微信客服[%v]开始.", kfaccount)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.AddKFAccount(access_token, kfaccount, nickname, password)
	if err != nil {
		LOG_ERROR("添加微信客服[%v]失败，失败原因：%v", kfaccount, err)
		return "", err
	}

	LOG_INFO("添加微信客服[%v]成功", kfaccount)

	return
}

func (this *WeixinSenderClient) UpdateKFAccount(access_token string, kfaccount string, nickname string, password string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("更新微信客服[%v]开始.", kfaccount)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.UpdateKFAccount(access_token, kfaccount, nickname, password)
	if err != nil {
		LOG_ERROR("更新微信客服[%v]失败，失败原因：%v", kfaccount, err)
		return "", err
	}

	LOG_INFO("更新微信客服[%v]成功", kfaccount)

	return
}

func (this *WeixinSenderClient) DeleteKFAccount(access_token string, kfaccount string, nickname string, password string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("删除微信客服[%v]开始.", kfaccount)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.DeleteKFAccount(access_token, kfaccount, nickname, password)
	if err != nil {
		LOG_ERROR("删除微信客服[%v]失败，失败原因：%v", kfaccount, err)
		return "", err
	}

	LOG_INFO("删除微信客服[%v]成功", kfaccount)

	return
}

func (this *WeixinSenderClient) SetKFHeadImg(access_token string, kfaccount string, media string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("设置微信客服[%v]的头像开始.", kfaccount)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.SetKFHeadImg(access_token, kfaccount, media)
	if err != nil {
		LOG_ERROR("设置微信客服[%v]的头像失败，失败原因：%v", kfaccount, err)
		return "", err
	}

	LOG_INFO("设置微信客服[%v]的头像成功", kfaccount)

	return
}

func (this *WeixinSenderClient) GetKFAccountList(access_token string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微信客服列表开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetKFAccountList(access_token)
	if err != nil {
		LOG_ERROR("获取微信客服列表失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("获取微信客服列表成功")

	return
}

func (this *WeixinSenderClient) SendMessage(access_token string, touser string, type_a1 string, data string, kfaccount string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("发送微信消息[%v]给[%v]开始.", data, touser)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.SendMessage(access_token, touser, type_a1, data, kfaccount)
	if err != nil {
		LOG_ERROR("发送微信消息[%v]给[%v]失败，失败原因：%v", data, touser, err)
		return "", err
	}

	LOG_INFO("发送微信消息[%v]给[%v]成功", data, touser)

	return
}

func (this *WeixinSenderClient) UploadTempMedia(access_token string, type_a1 string, media string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("上传临时媒体资源开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.UploadTempMedia(access_token, type_a1, media)
	if err != nil {
		LOG_ERROR("上传临时媒体资源失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("上传临时媒体资源成功")

	return
}

func (this *WeixinSenderClient) DownloadTempMedia(access_token string, media_id string) (r []byte, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("下载临时媒体资源[%v]开始.", media_id)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return nil, err
	}

	r, err = client.DownloadTempMedia(access_token, media_id)
	if err != nil {
		LOG_ERROR("下载临时媒体资源[%v]失败，失败原因：%v", media_id, err)
		return nil, err
	}

	LOG_INFO("下载临时媒体资源[%v]成功", media_id)

	return
}

func (this *WeixinSenderClient) UploadPermanentMedia(access_token string, type_a1 string, media string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("上传永久媒体资源开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.UploadPermanentMedia(access_token, type_a1, media)
	if err != nil {
		LOG_ERROR("上传永久媒体资源失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("上传永久媒体资源成功")

	return
}

func (this *WeixinSenderClient) DownloadPermanentMedia(access_token string, media_id string) (r []byte, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("下载永久媒体资源[%v]开始.", media_id)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return nil, err
	}

	r, err = client.DownloadPermanentMedia(access_token, media_id)
	if err != nil {
		LOG_ERROR("下载永久媒体资源[%v]失败，失败原因：%v", media_id, err)
		return nil, err
	}

	LOG_INFO("下载永久媒体资源[%v]成功", media_id)

	return
}

func (this *WeixinSenderClient) DeletePermanentMedia(access_token string, media_id string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("删除永久媒体资源开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.DeletePermanentMedia(access_token, media_id)
	if err != nil {
		LOG_ERROR("删除永久媒体资源失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("删除永久媒体资源成功")

	return
}

func (this *WeixinSenderClient) UploadNews(access_token string, news []byte) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("上传图文信息[%v]开始.", string(news))

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.UploadNews(access_token, news)
	if err != nil {
		LOG_ERROR("上传图文信息[%v]失败，失败原因：%v", string(news), err)
		return "", err
	}

	LOG_INFO("上传图文信息[%v]成功", string(news))

	return
}

func (this *WeixinSenderClient) SendNews(access_token string, is_to_all bool, group_id string, msg_type string, content string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("群发图文信息[%v]开始.", content)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		outputStr = fmt.Sprintf("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		LOG_ERROR(outputStr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.SendNews(access_token, is_to_all, group_id, msg_type, content)
	if err != nil {
		LOG_ERROR("群发图文信息[%v]失败，失败原因：%v", content, err)
		return "", err
	}

	LOG_INFO("群发图文信息[%v]成功", content)

	return
}

func (this *WeixinSenderClient) DeleteNews(access_token string, msg_id string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("删除图文信息[%v]开始.", msg_id)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.DeleteNews(access_token, msg_id)
	if err != nil {
		LOG_ERROR("删除图文信息[%v]失败，失败原因：%v", msg_id, err)
		return "", err
	}

	LOG_INFO("删除图文信息[%v]成功", msg_id)

	return
}

func (this *WeixinSenderClient) PreviewNewsByOpenId(access_token string, touser string, msg_type string, content string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("预览发送图文消息[%v]给[%v]开始.", content, touser)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.PreviewNewsByOpenId(access_token, touser, msg_type, content)
	if err != nil {
		LOG_ERROR("预览发送图文消息[%v]给[%v]失败，失败原因：%v", content, touser, err)
		return "", err
	}

	LOG_INFO("预览发送图文消息[%v]给[%v]成功", content, touser)

	return
}

func (this *WeixinSenderClient) GetNewsStatus(access_token string, msg_id string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取图文消息[%v]的状态开始.", msg_id)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetNewsStatus(access_token, msg_id)
	if err != nil {
		LOG_ERROR("获取图文消息[%v]的状态失败，失败原因：%v", msg_id, err)
		return "", err
	}

	LOG_INFO("获取图文消息[%v]的状态成功", msg_id)

	return
}

func (this *WeixinSenderClient) CreateUserGroup(access_token string, group_name string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("创建微信用户组[%v]开始.", group_name)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.CreateUserGroup(access_token, group_name)
	if err != nil {
		LOG_ERROR("创建微信用户组[%v]失败，失败原因：%v", group_name, err)
		return "", err
	}

	LOG_INFO("创建微信用户组[%v]成功", group_name)

	return
}

func (this *WeixinSenderClient) UpdateUserGroup(access_token string, group_id string, new_group_name string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("微信用户组[%v]更名为[%v]开始.", group_id, new_group_name)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.UpdateUserGroup(access_token, group_id, new_group_name)
	if err != nil {
		LOG_ERROR("微信用户组[%v]更名为[%v]失败，失败原因：%v", group_id, new_group_name, err)
		return "", err
	}

	LOG_INFO("微信用户组[%v]更名为[%v]成功", group_id, new_group_name)

	return
}

func (this *WeixinSenderClient) GetUserGroupList(access_token string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微信用户组列表开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserGroupList(access_token)
	if err != nil {
		LOG_ERROR("获取微信用户组列表失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("获取微信用户组列表成功")

	return
}

func (this *WeixinSenderClient) GetUserGroupByOpenID(access_token string, openid string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微信用户[%v]的组开始.", openid)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserGroupByOpenID(access_token, openid)
	if err != nil {
		LOG_ERROR("获取微信用户[%v]的组失败，失败原因：%v", openid, err)
		return "", err
	}

	LOG_INFO("获取微信用户[%v]的组成功", openid)

	return
}

func (this *WeixinSenderClient) MoveUserToGroup(access_token string, openid_list []string, to_groupid string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("移动用户%v到组[%v]开始.", openid_list, to_groupid)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.MoveUserToGroup(access_token, openid_list, to_groupid)
	if err != nil {
		LOG_ERROR("移动用户%v到组[%v]失败，失败原因：%v", openid_list, to_groupid, err)
		return "", err
	}

	LOG_INFO("移动用户%v到组[%v]成功", openid_list, to_groupid)

	return
}

func (this *WeixinSenderClient) RemarkUser(access_token string, openid string, remark string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("备注用户[%v]为[%v]开始.", openid, remark)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.RemarkUser(access_token, openid, remark)
	if err != nil {
		LOG_ERROR("备注用户[%v]为[%v]失败，失败原因：%v", openid, remark, err)
		return "", err
	}

	LOG_INFO("备注用户[%v]为[%v]成功", openid, remark)

	return
}

func (this *WeixinSenderClient) GetUserInfo(access_token string, openid string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取用户[%v]的详细信息开始.", openid)

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserInfo(access_token, openid)
	if err != nil {
		LOG_ERROR("获取用户[%v]的详细信息失败，失败原因：%v", openid, err)
		return "", err
	}

	LOG_INFO("获取用户[%v]的详细信息成功", openid)

	return
}

func (this *WeixinSenderClient) GetUserList(access_token string, next_openid string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微信用户列表开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.GetUserList(access_token, next_openid)
	if err != nil {
		LOG_ERROR("获取微信用户列表失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("获取微信用户列表成功")

	return
}

func (this *WeixinSenderClient) CreateMenu(access_token string, menu_data []byte) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("创建微信菜单开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.CreateMenu(access_token, menu_data)
	if err != nil {
		LOG_ERROR("创建微信菜单失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("创建微信菜单成功")

	return

}

func (this *WeixinSenderClient) DeleteMenu(access_token string) (r string, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("删除微信菜单开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return "", fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return "", fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return "", err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return "", err
	}

	r, err = client.DeleteMenu(access_token)
	if err != nil {
		LOG_ERROR("删除微信菜单失败，失败原因：%v", err)
		return "", err
	}

	LOG_INFO("删除微信菜单成功")

	return

}

func (this *WeixinSenderClient) GetMenu(access_token string) (r []byte, err error) {
	var outputStr string
	var networkAddr string
	var addr, port string
	var addrIsSet, portIsSet bool

	LOG_INFO("获取微信菜单开始.")

	addr, addrIsSet = g_config.Get("weixin_sender.addr")
	port, portIsSet = g_config.Get("weixin_sender.port")

	if addrIsSet && portIsSet {
		if addr != "" && port != "" {
			networkAddr = fmt.Sprintf("%s:%s", addr, port)
		} else {
			outputStr = "WeixinSender服务的网络地址设置错误"
			LOG_ERROR(outputStr)
			return nil, fmt.Errorf(outputStr)
		}
	} else {
		outputStr = "WeixinSender服务的网络地址没有设置"
		LOG_ERROR(outputStr)
		return nil, fmt.Errorf(outputStr)
	}

	trans, err := thrift.NewTSocket(networkAddr)
	if err != nil {
		LOG_ERROR("创建到WeixinSender服务[%v]的连接失败", networkAddr)
		return nil, err
	}
	defer trans.Close()

	var protocolFactory thrift.TProtocolFactory
	protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()

	client := weixinsender.NewWeixinSenderClientFactory(trans, protocolFactory)
	if err = trans.Open(); err != nil {
		LOG_ERROR("打开到WeixinSender服务的连接失败, WeixinSender服务可能没有就绪，请检查服务状态")
		return nil, err
	}

	r, err = client.GetMenu(access_token)
	if err != nil {
		LOG_ERROR("获取微信菜单失败，失败原因：%v", err)
		return nil, err
	}

	LOG_INFO("获取微信菜单成功")

	return
}

type SNSClient struct {
	smsSenderClient    map[string]*SMSSenderClient
	weiboSenderClient  *WeiboSenderClient
	weixinSenderClient *WeixinSenderClient
}

func NewSNSClient() (*SNSClient, error) {
	snsClient := &SNSClient{}

	snsClient.smsSenderClient = make(map[string]*SMSSenderClient)
	channels, err := GetSMSChannels()
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		snsClient.smsSenderClient[channel] = &SMSSenderClient{smsProvider: channel}
	}

	snsClient.weiboSenderClient = &WeiboSenderClient{}
	snsClient.weixinSenderClient = &WeixinSenderClient{}

	return snsClient, nil
}

func (this *SNSClient) GetSMSSender(provider string) *SMSSenderClient {
	r, ok := this.smsSenderClient[provider]
	if ok && r != nil {
		return r
	} else {
		provider, _ := g_config.Get("service.smssender.major")
		r, _ = this.smsSenderClient[provider]
		return r
	}
}

func (this *SNSClient) GetWeiboSender() *WeiboSenderClient {
	return this.weiboSenderClient
}

func (this *SNSClient) GetWeixinSender() *WeixinSenderClient {
	return this.weixinSenderClient
}
