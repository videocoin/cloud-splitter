package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/cloud-splitter/service"
)

var (
	ServiceName string = "splitter"
	Version     string = "dev"
)

func main() {
	logger.Init(ServiceName, Version)  //nolint

	log := logrus.NewEntry(logrus.New())
	log = logrus.WithFields(logrus.Fields{
		"service": ServiceName,
		"version": Version,
	})

	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		log.Info(err.Error())
	} else {
		defer closer.Close()
	}

	config := &service.Config{
		Name:    ServiceName,
		Version: Version,
	}

	err = envconfig.Process(ServiceName, config)
	if err != nil {
		log.Fatal(err.Error())
	}

	config.Logger = log

	conn, err := grpcutil.Connect(config.StreamsRPCAddr, config.Logger.WithField("system", "streamscli"))
	if err != nil {
		log.Fatal(err.Error())
	}
	streams := pstreamsv1.NewStreamsServiceClient(conn)

	svc, err := service.NewService(config, streams)
	if err != nil {
		log.Fatal(err.Error())
	}

	signals := make(chan os.Signal, 1)
	exit := make(chan bool, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-signals

		log.Infof("recieved signal %s", sig)
		exit <- true
	}()

	log.Info("starting")
	go svc.Start()  //nolint

	<-exit

	log.Info("stopping")
	err = svc.Stop()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("stopped")
}
