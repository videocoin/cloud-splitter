package service

import (
	"github.com/sirupsen/logrus"
)

type Config struct {
	Name    string `envconfig:"-"`
	Version string `envconfig:"-"`

	PrivateRPCAddr string `default:"0.0.0.0:5103" envconfig:"PRIVATE_RPC_ADDR"`
	//DownloadDir    string `required:"true" default:"/tmp" envconfig:"DOWNLOAD_DIR"`
	OutputDir string `required:"true" default:"/tmp" envconfig:"OUTPUT_DIR"  description:"local folder for ts chunks"`
	Bucket    string `required:"true" default:"testvc01" envconfig:"BUCKET"`

	Logger *logrus.Entry `envconfig:"-"`
}
