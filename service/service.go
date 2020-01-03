package service

import (
	"github.com/videocoin/cloud-splitter/eventbus"
	"github.com/videocoin/cloud-splitter/rpc"
)

type Service struct {
	cfg *Config
	rpc *rpc.RpcServer
	eb  *eventbus.EventBus
}

func NewService(cfg *Config) (*Service, error) {

	rpcConfig := &rpc.RpcServerOpts{
		HlsDir: cfg.HLSDir,
		Addr:   cfg.RPCAddr,
		Logger: cfg.Logger.WithField("system", "rpc"),
	}

	ebConfig := &eventbus.Config{
		URI:    cfg.MQURI,
		Name:   cfg.Name,
		Logger: cfg.Logger.WithField("system", "eventbus"),
	}
	eb, err := eventbus.New(ebConfig)
	if err != nil {
		return nil, err
	}

	rpc, err := rpc.NewRpcServer(rpcConfig)
	if err != nil {
		return nil, err
	}

	svc := &Service{
		cfg: cfg,
		rpc: rpc,
		eb:  eb,
	}

	return svc, nil
}

func (s *Service) Start() error {
	go s.rpc.Start()
	go s.eb.Start()
	return nil
}

func (s *Service) Stop() error {
	s.eb.Stop()
	return nil
}
