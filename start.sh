#!/bin/sh

# 因为smssender-server服务程序中使用了新浪微博的短链接服务t.cn，
# 所以需要设置环境变量SINA_WEIBO_APP_KEY，对应的值为客户申请的AppKey
export SINA_WEIBO_APP_KEY=4228546392

# 在后台启动snsscheduler-server服务
nohup ./bin/snsscheduler-server &
