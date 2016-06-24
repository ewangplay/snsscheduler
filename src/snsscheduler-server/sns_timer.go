package main

import (
	"jzlservice/snsscheduler"
	"time"
)

type SNSTimer struct {
	snsMessageQueue chan *snsscheduler.SNSMessage
}

func NewSNSTimer() (*SNSTimer, error) {
	snsTimer := &SNSTimer{}
	snsTimer.snsMessageQueue = make(chan *snsscheduler.SNSMessage, g_queueCapacity)
	return snsTimer, nil
}

func (this *SNSTimer) Release() {
	close(this.snsMessageQueue)
}

func (this *SNSTimer) AddToTimer(msg *snsscheduler.SNSMessage) error {
	LOG_DEBUG("活动[%v]设置了定时发送或者周期发送，将被添加到定时器中", msg.TaskId)

	this.snsMessageQueue <- msg

	return nil
}

func (this *SNSTimer) Run() error {
	LOG_INFO("消息定时器启动...")
	go func() {
		for {
			select {
			case msg := <-this.snsMessageQueue:
				go this.msgTimer(msg)
			}
		}
	}()

	return nil
}

func (this *SNSTimer) msgTimer(msg *snsscheduler.SNSMessage) error {
	if msg.TriggerTime != "" {
		LOG_DEBUG("活动[%v]设置了首次触发时间[%v]", msg.TaskId, msg.TriggerTime)

		startTime, err := time.ParseInLocation("2006-01-02 15:04:05", msg.TriggerTime, time.Local)
		if err != nil {
			LOG_ERROR("解析活动[%v]的首次触发时间失败", msg.TaskId)
			return err
		}
		startTimeStamp := startTime.Unix()
		startPoint := startTimeStamp - time.Now().Unix()

		LOG_DEBUG("活动[%v]的首次触发离当前还有[%v]秒", msg.TaskId, startPoint)

		if startPoint > 0 {
			time.AfterFunc(time.Duration(startPoint)*time.Second, func() {
				LOG_DEBUG("活动[%v]的首次触发时间已到，添加到消息发送队列", msg.TaskId)

				g_snsQueue.AddToSendQueue(msg)
			})
			//等待首次触发时间到来后再向下执行
			time.Sleep(time.Duration(startPoint) * time.Second)
		} else {
			LOG_DEBUG("活动[%v]的首次触发时间小于当前时间，立即添加到消息发送队列", msg.TaskId)

			g_snsQueue.AddToSendQueue(msg)
		}
	}

	if msg.IsPeriod && msg.PeriodInterval > 0 {
		LOG_DEBUG("活动[%v]设置了定时发送，定时发送周期为：%v", msg.TaskId, msg.PeriodInterval)

		ticker := time.NewTicker(time.Duration(msg.PeriodInterval) * time.Second)
		for {
			select {
			case <-ticker.C:
				LOG_DEBUG("活动[%v]的定时触发时间已到，添加到消息发送队列", msg.TaskId)

				g_snsQueue.AddToSendQueue(msg)
			}
		}
	}
	return nil
}
