# 生成api业务代码 ， 进入"服务.../api/desc"目录下，执行下面命令 generate to ../
APIFILE=$1
goctl api go -api $APIFILE -dir ../
