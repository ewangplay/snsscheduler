#!/bin/sh

# 根据thrift接口定义文件自动生成目标代码
thrift -out src -r --gen go smssender.thrift
thrift -out src -r --gen go weibosender.thrift
thrift -out src -r --gen go weixinsender.thrift
thrift -out src -r --gen go snsscheduler.thrift

# 编译安装snsscheduler-server服务程序
go install snsscheduler-server
