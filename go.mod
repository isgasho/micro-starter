module github.com/micro-community/auth

go 1.15

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/dgraph-io/dgo/v200 v200.0.0-20200825025457-a38d5eaacbf8
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/micro-in-cn/starter-kit v1.18.0
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v3 v3.0.0-beta.2
	github.com/micro/micro/v3 v3.0.0-beta.3
	github.com/olivere/elastic/v7 v7.0.20
	github.com/sirupsen/logrus v1.6.0
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/redis.v5 v5.2.9
	gopkg.in/sohlich/elogrus.v7 v7.0.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/driver/sqlite v1.1.2
	gorm.io/gorm v1.20.1
)
