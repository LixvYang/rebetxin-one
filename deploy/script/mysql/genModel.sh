#!/usr/bin/env bash

# 使用方法：
# ../../../deploy/script/mysql/genModel.sh topicpurchase

#生成的表名
tables=$1
#表生成的genmodel目录
modeldir=./genModel

# 数据库配置
host=127.0.0.1
port=3306
dbname=betxin
username=root
passwd=123456

echo "开始创建库：$dbname 的表：$1"
# goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}" -dir=. --style=goZero
goctl model mysql datasource -url="${username}:${passwd}@tcp(${host}:${port})/${dbname}" -table="${tables}" -dir=. -cache=true --style=goZero
