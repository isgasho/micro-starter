module github.com/micro-community/auth

go 1.15

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/dgraph-io/dgo/v200 v200.0.0-20200916081436-9ff368ad829a
	github.com/go-redis/redis/v8 v8.1.3
	github.com/golang/protobuf v1.4.2
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200929133051-87e898f4fc62
	github.com/micro/micro/v3 v3.0.0-beta.5
	github.com/olivere/elastic/v7 v7.0.20
	github.com/sirupsen/logrus v1.6.0
	github.com/urfave/cli/v2 v2.2.0
	go.mongodb.org/mongo-driver v1.4.1
	go.uber.org/dig v1.10.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/sohlich/elogrus.v7 v7.0.0
	gorm.io/driver/mysql v1.0.1
	gorm.io/driver/sqlite v1.1.2
	gorm.io/gorm v1.20.1
)
