package service

import (
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-splitter/eventbus"
	"github.com/videocoin/cloud-splitter/rpc"
)

type Service struct {
	cfg *Config
	rpc *rpc.Server
	eb  *eventbus.EventBus
}

func NewService(cfg *Config, streams pstreamsv1.StreamsServiceClient) (*Service, error) {
	rpcConfig := &rpc.ServerOpts{
		HlsDir:      cfg.HLSDir,
		Addr:        cfg.RPCAddr,
		Logger:      cfg.Logger.WithField("system", "rpc"),
		Streams:     streams,
		SegmentTime: cfg.SegmentTime,
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

	rpc, err := rpc.NewServer(rpcConfig)
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

func (s *Service) Start(errCh chan error) {
	go func() {
		errCh <- s.rpc.Start()
	}()

	go func() {
		errCh <- s.eb.Start()
	}()
}

func (s *Service) Stop() error {
	return s.eb.Stop()
}
