module github.com/xhyonline/xdq

go 1.16

require (
	github.com/BurntSushi/toml v0.4.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.7.4
	github.com/go-redis/redis/v7 v7.2.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/jasonlvhit/gocron v0.0.1
	github.com/xhyonline/xutil v0.1.2021111118
	go.etcd.io/etcd v3.3.27+incompatible
	google.golang.org/grpc v1.38.0
	gorm.io/gorm v1.21.16
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
