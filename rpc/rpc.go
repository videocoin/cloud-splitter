package rpc

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
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

		logger.Info("split has been completed")

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
	args := []string{
		"-i",
		filepath,
		"-codec",
		"copy",
		"-f",
		"segment",
		"-segment_time",
		strconv.Itoa(s.segmentTime),
		"-segment_list",
		hlsDir + "/index.m3u8",
		hlsDir + "/%d.ts",
	}

	s.logger.Infof("ffmpeg: %s", strings.Join(args, " "))

	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
