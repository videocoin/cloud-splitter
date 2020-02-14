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

type ServerOpts struct {
	Addr        string
	Logger      *logrus.Entry
	HlsDir      string
	Streams     pstreamsv1.StreamsServiceClient
	SegmentTime int
}

type Server struct {
	addr        string
	hlsDir      string
	logger      *logrus.Entry
	grpc        *grpc.Server
	listen      net.Listener
	streams     pstreamsv1.StreamsServiceClient
	segmentTime int
}

func NewServer(opts *ServerOpts) (*Server, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	gServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	server := &Server{
		addr:        opts.Addr,
		hlsDir:      opts.HlsDir,
		logger:      opts.Logger.WithField("system", "splitterv1"),
		grpc:        gServer,
		listen:      listen,
		streams:     opts.Streams,
		segmentTime: opts.SegmentTime,
	}

	splitterv1.RegisterSplitterServiceServer(gServer, server)
	reflection.Register(gServer)

	return server, nil
}

func (s *Server) Start() error {
	s.logger.Infof("starting rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}
