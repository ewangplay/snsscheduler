package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/ewangplay/go-utils"
	"jzlservice/smssender"
	"jzlservice/snsscheduler"
	"jzlservice/weibosender"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type SNSQueue struct {
	smsQueue    chan *snsscheduler.SNSMessage
	weiboQueue  chan *snsscheduler.SNSMessage
	weixinQueue chan *snsscheduler.SNSMessage
}

func NewSNSQueue() (*SNSQueue, error) {
	snsQueue := &SNSQueue{}

	snsQueue.smsQueue = make(chan *snsscheduler.SNSMessage, g_queueCapacity)
	snsQueue.weiboQueue = make(chan *snsscheduler.SNSMessage, g_queueCapacity)
	snsQueue.weixinQueue = make(chan *snsscheduler.SNSMessage, g_queueCapacity)

	return snsQueue, nil
}

func (this *SNSQueue) Release() {
	close(this.smsQueue)
	close(this.weiboQueue)
	close(this.weixinQueue)
}

func (this *SNSQueue) AddToSendQueue(msg *snsscheduler.SNSMessage) error {
	switch msg.TypeA1 {
	case snsscheduler.SNS_Message_Type_SMS:
		LOG_DEBUG("任务[%v]为短信类型，被添加到短信发送队列", msg.TaskId)

		this.smsQueue <- msg
	case snsscheduler.SNS_Message_Type_Weibo:
		LOG_DEBUG("任务[%v]为微博类型，被添加到微博发送队列", msg.TaskId)

		this.weiboQueue <- msg
	case snsscheduler.SNS_Message_Type_Weixin:
		LOG_DEBUG("任务[%v]为微信类型，被添加到微信发送队列", msg.TaskId)

		this.weixinQueue <- msg
	default:
		LOG_ERROR("任务[%v]设置了错误的任务类型: %v, 该任务被丢弃", msg.TaskId, msg.TypeA1)
		return fmt.Errorf("错误的任务类型: %v", msg.TypeA1)
	}

	return nil
}

func (this *SNSQueue) Run() {
	LOG_INFO("消息发送器启动...")
	go this.runSMSSendQueue()
	go this.runWeiboSendQueue()
	go this.runWeixinSendQueue()
}

func (this *SNSQueue) runSMSSendQueue() {
	for {
		select {
		case msg := <-this.smsQueue:
			LOG_DEBUG("任务[%v]在短信发送队列中被调度，准备发送...", msg.TaskId)

			//发送短信
			for _, entry := range msg.Entries {
				LOG_DEBUG("任务[%v]中包含短信消息：%v", msg.TaskId, entry)

				//产生唯一的短信发送流水号
				rand := rand.New(rand.NewSource(time.Now().UnixNano()))
				smsSerialNumberInt64 := rand.Int63n(9999999999) + 1
				LOG_DEBUG("新产生的短信发送批次号: %v", smsSerialNumberInt64)
				smsSerialNumber := strconv.FormatInt(smsSerialNumberInt64, 10)

				//发送短信前，先把短信批次号和任务对应信息缓存起来
				go func(spnumber string, localMsg snsscheduler.SNSMessage, localEntry snsscheduler.SNSEntry) {
					content := fmt.Sprintf("%v:%v:%v:%v:%v", localMsg.CustomerId, localEntry.ReceiverId, localMsg.TaskId, localMsg.ResourceId, localMsg.GroupId)
					g_snsCache.Set(spnumber, content)
				}(smsSerialNumber, *msg, *entry)

				smsExtraInfo := &snsscheduler.SMSExtraInfo{}
				if entry.SmsExtraInfo != nil {
					*smsExtraInfo = *entry.SmsExtraInfo
				}

                //默认启用短链接处理
                smsExtraInfo.EnableShortUrl = true

				/*
					//Deleted by wxh at 2015-10.30
					//不再支持通过参数选择发送通道，直接通过配置文件来控制

					channel := smsExtraInfo.Channel
					if channel == "" {
						channel, ok := g_config.Get("service.smssender.major")
						if !ok || channel == "" {
							channel = "hl"
						}
					}
				*/

				channel, ok := g_config.Get("service.smssender.major")
				if !ok || channel == "" {
					channel, ok = g_config.Get("service.smssender.minor")
					if !ok || channel == "" {
						channel = "hl"
					}
				}

				receiverListAll := strings.FieldsFunc(entry.Receiver, func(char rune) bool {
					switch char {
					case ',':
						return true
					}
					return false
				})

                //手机号码去重
                receiverList := utils.DedupAndSort(receiverListAll)


				switch smsExtraInfo.Category {
				case snsscheduler.SMS_Message_Category_MARKET:

					sms_service_provider_maxphonenum, ok := g_config.Get("service.smssender.maxphonenum")
					if !ok || sms_service_provider_maxphonenum == "" {
						LOG_INFO("允许最大的用户发送数没有配置，设置成默认值20")
						sms_service_provider_maxphonenum = "20"
					}

					maxphonenum, err := strconv.ParseInt(sms_service_provider_maxphonenum, 0, 0)
					if err != nil {
						maxphonenum = 20
					}

					//assemble the request url
					var receiverTemp string
					var total, groupNum, i int64

					total = int64(len(receiverList))

					if total <= maxphonenum {
						go this.SendSMSMessage(
							msg.CustomerId,
							msg.TaskId,
							entry.ReceiverId,
							msg.ResourceId,
							msg.GroupId,
							channel,
							smsSerialNumber,
							entry.Receiver,
							entry.Content,
							smsExtraInfo.Signature,
							smsExtraInfo.ServiceMinorNumber,
							int8(smsExtraInfo.Category),
							smsExtraInfo.EnableShortUrl)
					} else {
						LOG_DEBUG("超出了短信服务提供商允许发送的最大用户数，需要对发送人列表进行拆分")

						groupNum = int64(total / maxphonenum)

						for i = 0; i < groupNum; i++ {
							receiverTemp = strings.Join(receiverList[i*maxphonenum:i*maxphonenum+maxphonenum], ",")

							go this.SendSMSMessage(
								msg.CustomerId,
								msg.TaskId,
								entry.ReceiverId,
								msg.ResourceId,
								msg.GroupId,
								channel,
								smsSerialNumber,
								receiverTemp,
								entry.Content,
								smsExtraInfo.Signature,
								smsExtraInfo.ServiceMinorNumber,
								int8(smsExtraInfo.Category),
								smsExtraInfo.EnableShortUrl)
						}

						if len(receiverList[groupNum*maxphonenum:]) > 0 {
							receiverTemp = strings.Join(receiverList[groupNum*maxphonenum:], ",")

							go this.SendSMSMessage(
								msg.CustomerId,
								msg.TaskId,
								entry.ReceiverId,
								msg.ResourceId,
								msg.GroupId,
								channel,
								smsSerialNumber,
								receiverTemp,
								entry.Content,
								smsExtraInfo.Signature,
								smsExtraInfo.ServiceMinorNumber,
								int8(smsExtraInfo.Category),
								smsExtraInfo.EnableShortUrl)
						}
					}

				case snsscheduler.SMS_Message_Category_NORMAL:
					fallthrough
				default:
					for _, receiver := range receiverList {
						LOG_DEBUG("任务[%v]中的短信消息[%v]开始发送给[%v]", msg.TaskId, entry, receiver)

						go this.SendSMSMessage(
							msg.CustomerId,
							msg.TaskId,
							entry.ReceiverId,
							msg.ResourceId,
							msg.GroupId,
							channel,
							smsSerialNumber,
							receiver,
							entry.Content,
							smsExtraInfo.Signature,
							smsExtraInfo.ServiceMinorNumber,
							int8(smsExtraInfo.Category),
							smsExtraInfo.EnableShortUrl)
					}
				}
			}
		}
	}
}

func (this *SNSQueue) runWeiboSendQueue() {
	for {
		select {
		case msg := <-this.weiboQueue:
			LOG_DEBUG("任务[%v]在微博发送队列中被调度，准备发送...", msg.TaskId)

			//发送微博
			for _, entry := range msg.Entries {
				weiboExtraInfo := &snsscheduler.WeiboExtraInfo{}
				if entry.WeiboExtraInfo != nil {
					*weiboExtraInfo = *entry.WeiboExtraInfo
				}

                //默认启用短链接处理
                weiboExtraInfo.EnableShortUrl = true

				senderList := strings.FieldsFunc(entry.Sender, func(char rune) bool {
					switch char {
					case ',':
						return true
					}
					return false
				})

				for _, sender := range senderList {
					go func(cid, task_id, resource_id, group_id int64, access_token string, content string, extraInfo *snsscheduler.WeiboExtraInfo) {
						var encoded_content string

						if extraInfo.EnableShortUrl {
							//针对BI做URL编码
							url_click_service_addr, ok := g_config.Get("url_click_service.addr")
							if !ok || url_click_service_addr == "" {
								url_click_service_addr = "http://url-staging.jiuzhilan.com/redirect"
							}
							encoded_content = utils.BIUrlEncode("1",
								strconv.FormatInt(cid, 10),
								"",
								strconv.FormatInt(task_id, 10),
								"",
								strconv.FormatInt(resource_id, 10),
								strconv.FormatInt(group_id, 10),
								fmt.Sprintf("%v", 4),
								content,
								access_token,
								url_click_service_addr,
								false,
								false)
						} else {
							encoded_content = content
						}

						weibo_status := &weibosender.WeiboStatus{
							TaskId:      task_id,
							AccessToken: access_token,
							Status:      encoded_content,
							Visible:     extraInfo.Visible,
							ListId:      extraInfo.ListId,
							Latitude:    extraInfo.Latitude,
							Longitude:   extraInfo.Longitude,
							Annotations: extraInfo.Annotations,
							RealIp:      extraInfo.RealIp,
							Pic:         extraInfo.Pic,
						}

						r, err := g_snsClient.GetWeiboSender().SendMessage(weibo_status)
						if err == nil {
							LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=200]", msg.CustomerId, msg.TaskId, 4, msg.GroupId, msg.ResourceId)
							LOG_INFO("Send weibo message success: %v", r)

							//把获取的msg_id保存到数据库，以方便后续获取该微博的互动数（评论数、转发数、点赞数)
							this.saveWeiboStatus(msg.CustomerId, msg.ResourceId, r)
						} else {
							LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=300]", msg.CustomerId, msg.TaskId, 4, msg.GroupId, msg.ResourceId)
							LOG_ERROR("Send weibo message failure: %v", err)
						}
					}(msg.CustomerId, msg.TaskId, msg.ResourceId, msg.GroupId, sender, entry.Content, weiboExtraInfo)
				}
			}
		}
	}
}

type WeiboResultInfo struct {
	Created_at string `json:"created_at"`
	Id         int64  `json:"id"`
	Idstr      string `json:"idstr"`
}

const SET_WEIBO_MSG_ID_SQL string = "UPDATE jzl_weibo SET is_publish=1, msg_id=?, publish_time=? WHERE cid=? AND weibo_id=? AND is_delete=0"

func (this *SNSQueue) saveWeiboStatus(cid, weibo_id int64, result string) error {
	var err error
	var info WeiboResultInfo
	err = json.Unmarshal([]byte(result), &info)
	if err != nil {
		LOG_ERROR("Unmarshal error: %v", err)
		return err
	}

	//格式化微博的创建时间
	t, err := time.Parse("Mon Jan 2 15:04:05 -0700 2006", info.Created_at)
	if err != nil {
		return err
	}
	create_time := t.Format("2006-01-02 15:04:05")

	err = g_mysqladaptor.ExecFormat(SET_WEIBO_MSG_ID_SQL, info.Idstr, create_time, cid, weibo_id)
	if err != nil {
		LOG_ERROR("Update weibo msg_id error: %v", err)
		return err
	}

	LOG_INFO("Save weibo msg_id succ!")

	return nil
}

func (this *SNSQueue) runWeixinSendQueue() {
	for {
		select {
		case msg := <-this.weixinQueue:
			LOG_DEBUG("任务[%v]在微信发送队列中被调度，准备发送...", msg.TaskId)

			//发送微信
			for _, entry := range msg.Entries {
				weixinExtraInfo := &snsscheduler.WeixinExtraInfo{}
				if entry.WeixinExtraInfo != nil {
					*weixinExtraInfo = *entry.WeixinExtraInfo
				}

				senderList := strings.FieldsFunc(entry.Sender, func(char rune) bool {
					switch char {
					case ',':
						return true
					}
					return false
				})

				for _, sender := range senderList {
					go func(task_id int64, access_token string, group_ids string, content string, extraInfo *snsscheduler.WeixinExtraInfo) {
						if group_ids != "" {
							groupidList := strings.FieldsFunc(group_ids, func(char rune) bool {
								switch char {
								case ',':
									return true
								}
								return false
							})

							for _, group_id := range groupidList {
								r, err := g_snsClient.GetWeixinSender().SendNews(access_token, false, group_id, extraInfo.Msgtype, content)
								if err == nil {
									LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=200]", msg.CustomerId, msg.TaskId, 5, msg.GroupId, msg.ResourceId)
									LOG_INFO("Send weixin message success: %v", r)

									//把获取的msg_id保存到数据库，以方便后续获取该微信的阅读数和分享数
									this.saveWeixinStatus(msg.CustomerId, msg.ResourceId, r)
								} else {
									LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=300]", msg.CustomerId, msg.TaskId, 5, msg.GroupId, msg.ResourceId)
									LOG_ERROR("Send weixin message failure: %v", err)
								}
							}
						} else {
							r, err := g_snsClient.GetWeixinSender().SendNews(access_token, true, "", extraInfo.Msgtype, content)
							if err == nil {
								LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=200]", msg.CustomerId, msg.TaskId, 5, msg.GroupId, msg.ResourceId)
								LOG_INFO("Send weixin message success: %v", r)

								//把获取的msg_id保存到数据库，以方便后续获取该微信的阅读数和分享数
								this.saveWeixinStatus(msg.CustomerId, msg.ResourceId, r)
							} else {
								LOG_INFO("bi[c=%v t=%v y=%v g=%v d=%v s=300]", msg.CustomerId, msg.TaskId, 5, msg.GroupId, msg.ResourceId)
								LOG_ERROR("Send weixin message failure: %v", err)
							}
						}
					}(msg.TaskId, sender, entry.Receiver, entry.Content, weixinExtraInfo)
				}
			}
		}
	}
}

type WeixinResultInfo struct {
	Errcode     int64  `json:"errcode"`
	Errmsg      string `json:"errmsg"`
	Msg_id      int64  `json:"msg_id"`
	Msg_data_id int64  `json:"msg_data_id"`
}

const SET_WEIXIN_MSG_ID_SQL string = "UPDATE jzl_weixin SET is_publish=1, msg_id=?, msg_data_id=?, publish_time=? WHERE cid=? AND weixin_id=? AND is_delete=0"

func (this *SNSQueue) saveWeixinStatus(cid, weixin_id int64, result string) error {
	var err error
	var info WeixinResultInfo
	err = json.Unmarshal([]byte(result), &info)
	if err != nil {
		LOG_ERROR("Unmarshal error: %v", err)
		return err
	}

	//微信的创建时间
	create_time := time.Now().Format("2006-01-02 15:04:05")

	err = g_mysqladaptor.ExecFormat(SET_WEIXIN_MSG_ID_SQL, info.Msg_id, info.Msg_data_id, create_time, cid, weixin_id)
	if err != nil {
		LOG_ERROR("Update weixin msg_id error: %v", err)
		return err
	}

	LOG_INFO("Save weixin msg_id succ!")

	return nil
}

func (this *SNSQueue) SendSMSMessage(cid, task_id, contact_id, resource_id, group_id int64, channel string, spnumber string, phone string, content string, signature string, service_number string, catetory int8, enable_short_url bool) {
	var err error
	var r *smssender.SMSStatus

	smsSender := g_snsClient.GetSMSSender(channel)
	if smsSender != nil {

		var encoded_content string

		//针对BI做URL编码
		if enable_short_url {
			url_click_service_addr, ok := g_config.Get("url_click_service.addr")
			if !ok || url_click_service_addr == "" {
				url_click_service_addr = "http://url-staging.jiuzhilan.com/redirect"
			}
			encoded_content = utils.BIUrlEncode("1",
				strconv.FormatInt(cid, 10),
				strconv.FormatInt(contact_id, 10),
				strconv.FormatInt(task_id, 10),
				spnumber,
				strconv.FormatInt(resource_id, 10),
				strconv.FormatInt(group_id, 10),
				fmt.Sprintf("%v", 3),
				content,
				phone,
				url_click_service_addr,
				false,
				true)

		} else {
			encoded_content = content
		}

		r, err = smsSender.SendMessage(task_id, spnumber, phone, encoded_content, signature, service_number, catetory)

		phoneList := strings.FieldsFunc(phone, func(char rune) bool {
			switch char {
			case ',':
				return true
			}
			return false
		})

		if err == nil {
			for _, mobile := range phoneList {
				if r.Status == "0" {
					LOG_INFO("bi[c=%v t=%v y=%v g=%v u=%v d=%v w=%v n=%v s=100]", cid, task_id, 3, group_id, contact_id, resource_id, mobile, channel)
					LOG_INFO("Send sms message success: %v", r.Message)
				} else {
					LOG_INFO("bi[c=%v t=%v y=%v g=%v u=%v d=%v w=%v n=%v s=300]", cid, task_id, 3, group_id, contact_id, resource_id, mobile, channel)
					LOG_ERROR("Send sms message failure: %v", r.Message)
				}
			}
		} else {
			for _, mobile := range phoneList {
				LOG_INFO("bi[c=%v t=%v y=%v g=%v u=%v d=%v w=%v n=%v s=300]", cid, task_id, 3, group_id, contact_id, resource_id, mobile, channel)
				LOG_ERROR("Send sms message failure: %v", err)
			}
		}
	}
}
