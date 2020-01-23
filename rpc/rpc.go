package rpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	splitterv1 "github.com/videocoin/cloud-api/splitter/v1"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
)

func (s *RpcServer) Split(ctx context.Context, req *splitterv1.SplitRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.StreamID)

	localMediaDir := filepath.Join(s.hlsDir, req.StreamID)
	err := os.MkdirAll(localMediaDir, 0777)
	if err != nil {
		return nil, err
	}

	go func(filepath string, hlsDir string) {
		logger := s.logger.WithFields(logrus.Fields{
			"stream_id": req.StreamID,
			"path":      filepath,
		})
		logger.Info("splitting")

		err := s.splitMediafile(filepath, hlsDir)
		if err != nil {
			logger.Errorf("failed to split media file: %s", err)
			return
		}

		logger.Info("split has beend completed")

		streamReq := &pstreamsv1.StreamRequest{Id: req.StreamID}
		_, err = s.streams.Publish(context.Background(), streamReq)
		if err != nil {
			logger.Errorf("failed to publish stream: %s", err)
			return
		}

	}(req.Filepath, localMediaDir)

	return &protoempty.Empty{}, nil
}

func (s *RpcServer) splitMediafile(filepath string, hlsDir string) error {
	output := hlsDir + "/index.m3u8"
	args := []string{
		"-i",
		filepath,
		"-acodec",
		"copy",
		"-vcodec",
		"copy",
		"-reset_timestamps",
		"1",
		"-bsf:v",
		"h264_mp4toannexb",
		"-f",
		"hls",
		output,
	}
	fmt.Println("ffmpeg " + strings.Join(args, " "))
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
