package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string        `envconfig:"-"`
	Version string        `envconfig:"-"`
	Logger  *logrus.Entry `envconfig:"-"`

	RPCAddr        string `default:"0.0.0.0:5103" envconfig:"RPC_ADDR"`
	MQURI          string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`
	StreamsRPCAddr string `required:"true" envconfig:"STREAMS_RPC_ADDR" default:"127.0.0.1:5102"`
	HLSDir         string `required:"true" default:"/tmp/hls" envconfig:"HLS_DIR"`
	SegmentTime    int    `default:"10"`
}
