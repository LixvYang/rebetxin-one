#生成的表名
MODEL=$1

# 数据库配置
host=localhost
port=3306
dbname=betxin
username=root
passwd=123456

echo "开始创建库：$dbname 的表：$1"

sql2pb -go_package ./pb -host $host -package pb -password $passwd -port $port -schema $dbname -service_name ${MODEL}srv -user $root -table=$MODEL > ${MODEL}srv.proto

# sql2pb -go_package ./pb -host localhost -package pb -password 123456 -port 3306 -schema betxin -service_name topicpurchasesrv -user root -table=topicpurchase > topicpurchasesrv.proto