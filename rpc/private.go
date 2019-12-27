package rpc

import (
	"context"
	"net"
	"path/filepath"
	//"fmt"
	"os/exec"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	//"github.com/videocoin/cloud-api/rpc"
	protoempty "github.com/gogo/protobuf/types"
	privatev1 "github.com/videocoin/cloud-api/splitter/private/v1"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type PrivateRPCServerOpts struct {
	Addr      string
	Logger    *logrus.Entry
	Bucket    string
	OutputDir string
}

type PrivateRPCServer struct {
	addr      string
	bucket    string
	OutputDir string
	logger    *logrus.Entry
	grpc      *grpc.Server
	listen    net.Listener
	uploader  *Uploader
}

func NewPrivateRPCServer(opts *PrivateRPCServerOpts) (*PrivateRPCServer, error) {
	grpcOpts := grpcutil.DefaultServerOpts(opts.Logger)
	grpcServer := grpc.NewServer(grpcOpts...)
	healthService := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthService)
	listen, err := net.Listen("tcp", opts.Addr)
	if err != nil {
		return nil, err
	}
	uploader, _ := NewUploader(opts.Bucket, opts.Logger.WithField("system", "privaterpc"))
	rpcServer := &PrivateRPCServer{
		addr:      opts.Addr,
		OutputDir: opts.OutputDir,
		logger:    opts.Logger.WithField("system", "privaterpc"),
		grpc:      grpcServer,
		listen:    listen,
		uploader:  uploader,
	}

	privatev1.RegisterSplitterServiceServer(grpcServer, rpcServer)
	reflection.Register(grpcServer)

	return rpcServer, nil
}

func (s *PrivateRPCServer) Start() error {
	s.logger.Infof("starting private rpc server on %s", s.addr)
	return s.grpc.Serve(s.listen)
}

func (s *PrivateRPCServer) Split(ctx context.Context, req *privatev1.SplitRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.StreamID)
	//logger := s.logger.WithField("id", req.StreamID)
	localMediaDir := filepath.Join(s.OutputDir, req.StreamID)

	err := splitMediafile(req.Filepath, localMediaDir)
	if err != nil {
		return nil, err
	}
	s.uploader.uploadSegments(req.StreamID, localMediaDir)
	//stream, err := s.streams.GetStreamByID(ctx, req.StreamID)
	//if err != nil {
	//	logFailedTo(logger, "get stream", err)
	//	return nil, rpc.ErrRpcInternal
	//}
	return &protoempty.Empty{}, nil
}

func splitMediafile(filepath string, outputDir string) error {
	//output := fmt.Sprintf("%s.png", uuid.New().String())
	args := []string{
		"-i",
		filepath,
		"-acodec",
		"copy",
		"-f",
		"segment",
		"-vcodec",
		"copy",
		"-reset_timestamps",
		"1",
		"-map",
		"0",
		outputDir,
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
