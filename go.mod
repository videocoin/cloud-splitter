module github.com/videocoin/cloud-splitter

go 1.12

require (
	cloud.google.com/go/storage v1.4.0 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/google/btree v1.0.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20190404075320-75d898a42a94
	github.com/videocoin/cloud-api v0.2.15
	github.com/videocoin/cloud-pkg v0.0.6
	google.golang.org/grpc v1.26.0
)

replace github.com/videocoin/cloud-pkg => ../cloud-pkg

replace github.com/videocoin/cloud-api => ../cloud-api
