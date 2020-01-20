package rpc

import (
	"net"

	"github.com/sirupsen/logrus"
	splitterv1 "github.com/videocoin/cloud-api/splitter/v1"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type RpcServerOpts struct {
	Addr    string
	Logger  *logrus.Entry
	HlsDir  string
	Streams pstreamsv1.StreamsServiceClient
}

type RpcServer struct {
	addr    string
	hlsDir  string
	logger  *logrus.Entry
	grpc    *grpc.Server
	listen  net.Listener
	streams pstreamsv1.StreamsServiceClient
}

func NewRpcServer(opts *RpcServerOpts) (*RpcServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	rpcServer := &RpcServer{
		addr:    opts.Addr,
		hlsDir:  opts.HlsDir,
		logger:  opts.Logger.WithField("system", "splitterv1"),
		grpc:    grpcServer,
		listen:  listen,
		streams: opts.Streams,
	}

	splitterv1.RegisterSplitterServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *RpcServer) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}
