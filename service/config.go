package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	RPCAddr string `default:"0.0.0.0:5103" envconfig:"RPC_ADDR"`
	MQURI   string `default:"amqp://guest:guest@127.0.0.1:5672" envconfig:"MQURI"`

	HLSDir string `required:"true" default:"/tmp/hls" envconfig:"HLS_DIR"`

	Logger *logrus.Entry `envconfig:"-"`
}
