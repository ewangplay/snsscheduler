/* 
 * thrift interface for snsscheduler
 */

namespace cpp jzlservice.snsscheduler
namespace go jzlservice.snsscheduler
namespace py jzlservice.snsscheduler
namespace php jzlservice.snsscheduler
namespace perl jzlservice.snsscheduler
namespace java jzlservice.snsscheduler

/** 
* 发送消息类型
* 1、LandingPage；2、Email；3、Sms；4、Weibo；5、Weixin；6、CTA；7、Form
*/
enum SNS_Message_Type {
    UNSET = 0,                      //未设置
    SMS = 3,                        //短信类型
    Weibo = 4,                      //微博类型
    Weixin = 5,                     //微信类型
}

enum SMS_Message_Category {
    UNSET = 0,                      //未设置
    NORMAL = 1,                     //普通类短信(包括验证短信，通知短信等)
    MARKET = 2,                     //营销类短信
}

struct WeiboExtraInfo {
    1: i32 visible = 0,             //微博的可见性，0：所有人能看，1：仅自己可见，2：密友可见，3：指定分组可见，默认为0

    2: string list_id = "",         //微博的保护投递指定分组ID，只有当visible参数为3时生效且必选

    3: double latitude = 0.0,       //纬度，有效范围：-90.0到+90.0，+表示北纬，默认为0.0

    4: double longitude = 0.0,      //经度，有效范围：-180.0到+180.0，+表示东经，默认为0.0

    5: string annotations = "",     //或者多个元数据，必须以json字串的形式提交，字串长度不超过512个字符，具体内容可以自定

    6: string real_ip = "",         //开发者上报的操作用户真实IP，形如：211.156.0.1

    7: string pic = "",             //SNS消息附带的图片数据

    8: bool enable_short_url = false,       //是否启用短链接功能，即对微博中的url链接进行短链接处理
}

struct SMSExtraInfo {
    1: string signature = "",               //短信附加的签名信息，如：【枝兰】。注意不包含中文的方括号边界，只是签名信息本身，比如：枝兰

    2: string service_minor_number = "",    //分配给服务调用者的服务小号，用来标示不同的调用者身份

    3: SMS_Message_Category category = 0,   //短信所属的类别，比如验证类短信，营销类短息等等

    4: string channel = "",                 //指定短信发送的通道，比如hl, wxtl等，如果为空，则按照系统默认通道发送

    5: bool enable_short_url = false,       //是否启用短链接功能，即对短信中的url链接进行短链接处理
}

struct WeixinExtraInfo {
    1: string msgtype = "",                 //群发的消息类型, mpnews: 图文，text: 文本，voice: 语音，music: 音乐，image: 图片，video: 视频
}

/**
* struct SNSEntry
* sns entry structure description.
*/
struct SNSEntry {
    1: string content = "",                 //SNS消息的内容

    2: i64 receiver_id = 0,                 //联系人ID

    3: string receiver = "",                //SNS消息的接收者，多个之间用“，”分隔; 
                                            //如果是短信消息，该字段为要发送消息的手机号码列表
                                            //如果是微信和微博消息，该字段设置为空表示发布给发所有，设置为组ID表示发送给组成员

    4: string sender = "",                  //SNS消息的发送者，多个之间用“，”分隔; 
                                            //对于短信消息，该字段设置为空；
                                            //对于微博和微信消息，该字段为acess_token列表

    5: SMSExtraInfo sms_extra_info,         //短信消息的附加信息

    6: WeiboExtraInfo weibo_extra_info,     //微博消息的附加信息

    7: WeixinExtraInfo weixin_extra_info,   //微信消息的附加信息
}

/**
* struct SNSMessage
* sns message structure description.
*/
struct SNSMessage {
    1: i64 customer_id = 0,         //公司ID

    2: i64 operator_id = 0,         //操作者ID

    3: i64 task_id = 0,             //任务ID 

    4: i64 resource_id = 0,         //资源ID

    5: i64 group_id = 0,            //发送组ID

    6: SNS_Message_Type type,       //SNS消息的类型。3: 短信, 4: 微博, 5: 微信

    7: list<SNSEntry> entries,      //SNS消息列表

    8: bool is_period = false,	    //是否周期发送

    9: string trigger_time = "",	//SNS消息开始发送的时间，格式：%Y-%m-%d %H:%M:%S；如果设置为空，表示立即发送
    
    10: i32 period_interval = 0,	//周期发送的时间间隔，单位为秒；只有当设置is_period为true时才生效
}

/**
* 社交媒体调度服务
*/
service SNSScheduler {
    //================================================================================
    //Common API

    /** 
    * @描述: 
    *   服务连通性测试接口
    *
    * @返回: 
    *   返回pong表示服务正常；返回空或其它表示服务异常
    */
    string ping(),		            

    /**
    * @描述: 
    *   发送SNS消息接口，可以一次发送多条消息
    *
    * @参数: 
    *   msgs: SNSMessage列表
    *
    * @返回: 
    *   返回true只是表示成功添加到发送队列，返回false表示添加到发送队列失败
    */
    bool sendMessage(1: list<SNSMessage> msgs),        

    //===============================================================================
    // SMS API
    /** 
    * @描述: 
    *   获取短信服务的通道列表
    *
    * @返回:
    *   返回短信服务的通道列表
    */
    list<string> getSMSServiceChannel(),

    /** 
    * @描述: 
    *   获取短信服务的主号码
    *
    * @返回:
    *   返回短信服务的主号码
    *
    * @附加说明:
    *   实际发送短信的服务号 = 短信服务的主号码 + 用户分配的小号，比如短信服务主号码 = 10690266002, 用户分配的小号 = 001，
    *   那么最终发送短信的服务号就是: 10690266002001
    */
    string getSMSServiceNumber(1: string channel),
}

