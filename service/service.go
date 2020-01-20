package service

import (
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-splitter/eventbus"
	"github.com/videocoin/cloud-splitter/rpc"
)

type Service struct {
	cfg *Config
	rpc *rpc.RpcServer
	eb  *eventbus.EventBus
}

func NewService(cfg *Config) (*Service, error) {
	conn, err := grpcutil.Connect(cfg.StreamsRPCAddr, cfg.Logger.WithField("system", "streamscli"))
	if err != nil {
		return nil, err
	}
	streams := pstreamsv1.NewStreamsServiceClient(conn)

	rpcConfig := &rpc.RpcServerOpts{
		HlsDir:  cfg.HLSDir,
		Addr:    cfg.RPCAddr,
		Logger:  cfg.Logger.WithField("system", "rpc"),
		Streams: streams,
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
