package service

import (
	"github.com/videocoin/cloud-splitter/rpc"
)

type Service struct {
	cfg        *Config
	privateRPC *rpc.PrivateRPCServer
}

func NewService(cfg *Config) (*Service, error) {

	privateRPCConfig := &rpc.PrivateRPCServerOpts{
		Bucket:    cfg.Bucket,
		OutputDir: cfg.OutputDir,
		Addr:      cfg.PrivateRPCAddr,
		Logger:    cfg.Logger.WithField("system", "privaterpc"),
	}

	privateRPC, err := rpc.NewPrivateRPCServer(privateRPCConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg:        cfg,
		privateRPC: privateRPC,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.privateRPC.Start()
	return nil
}

func (s *Service) Stop() error {
	return nil
}
