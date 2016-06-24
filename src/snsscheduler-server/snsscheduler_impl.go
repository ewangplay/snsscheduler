package main

import (
	"fmt"
	"jzlservice/snsscheduler"
	"strings"
)

const (
	RETRY_MAX_COUNT = 3
)

type SNSSchedulerImpl struct {
}

func (this *SNSSchedulerImpl) Ping() (r string, err error) {
	LOG_INFO("请求ping方法")
	return "pong", nil
}

func (this *SNSSchedulerImpl) SendMessage(msgs []*snsscheduler.SNSMessage) (r bool, err error) {
	LOG_INFO("请求SendMessage方法")

	for _, msg := range msgs {
		if msg.TriggerTime != "" || msg.IsPeriod {
			g_snsTimer.AddToTimer(msg)
		} else {
			g_snsQueue.AddToSendQueue(msg)
		}
	}

	return true, nil
}

func (this *SNSSchedulerImpl) GetSMSServiceChannel() (r []string, err error) {
	LOG_INFO("请求GetSMSServiceChannel方法")
	r, err = GetSMSChannels()
	return
}

func (this *SNSSchedulerImpl) GetSMSServiceNumber(channel string) (r string, err error) {
	LOG_INFO("请求GetSMSServiceNumber方法")

	if channel == "" {
		channel, _ = g_config.Get("service.smssender.major")
	}

	key := fmt.Sprintf("sms_sender_%v.serviceNumber", channel)
	r, _ = g_config.Get(key)
	if r == "" {
		return "", fmt.Errorf("该通道没有设置对应的服务号，请联系服务提供商")
	}

	return
}

func GetSMSChannels() (r []string, err error) {
	majorChannel, _ := g_config.Get("service.smssender.major")
	if majorChannel != "" {
		r = append(r, majorChannel)
	} else {
		return nil, fmt.Errorf("没有指定缺省的短信发送服务")
	}

	minorChannels, _ := g_config.Get("service.smssender.minor")
	if minorChannels != "" {
		channelList := strings.FieldsFunc(minorChannels, func(char rune) bool {
			switch char {
			case ',':
				return true
			}
			return false
		})
		for _, channel := range channelList {
			r = append(r, channel)
		}
	}
	return
}
