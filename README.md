# midmsg
the message transfer service

config.yaml 修改服务器端口 数据库ip 端口 密码等

maxwoker : 1000  设置协程数量
maxqueue : 10    设置队列数据


git clone https://github.com/fieldlee/midmsg.git

go mod download
go mod vendor

go build ./ 
