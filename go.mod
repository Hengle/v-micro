module github.com/fananchong/v-micro

go 1.12

require (
	github.com/fananchong/gotcp v0.0.0-20190826033715-2799541f0ab3
	github.com/fananchong/protoc-gen-vmicro v0.0.4
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/lestrrat-go/file-rotatelogs v2.2.0+incompatible
	github.com/lestrrat-go/strftime v0.0.0-20190725011945-5c849dd2c51d // indirect
	github.com/micro/cli v0.2.0
	github.com/micro/mdns v0.3.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/sirupsen/logrus v1.4.2
)

replace github.com/micro/mdns v0.3.0 => github.com/fananchong/mdns v0.3.1-0.20190816031422-53f87c9ecf9a82e2464601ef5305751a5da9b3ad
