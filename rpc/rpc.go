package rpc

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/vansante/go-ffprobe"
	splitterv1 "github.com/videocoin/cloud-api/splitter/v1"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
)

func (s *Server) Split(ctx context.Context, req *splitterv1.SplitRequest) (*protoempty.Empty, error) {
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

			_, stopErr := s.streams.Stop(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamID})
			if stopErr != nil {
				logger.Errorf("failed to stop stream: %s", err)
			}

			return
		}

		logger.Info("split has been completed")

		mediadata, err := ffprobe.GetProbeData(filepath, time.Second*5)
		if err != nil {
			logger.Errorf("failed to get probe data: %s", err)

			_, stopErr := s.streams.Stop(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamID})
			if stopErr != nil {
				logger.Errorf("failed to stop stream: %s", err)
			}

			return
		}

		stream := mediadata.GetFirstVideoStream()
		if stream == nil {
			logger.Errorf("failed to get stream data: %s", err)

			_, stopErr := s.streams.Stop(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamID})
			if stopErr != nil {
				logger.Errorf("failed to stop stream: %s", err)
			}

			return
		}

		duration, err := strconv.ParseFloat(stream.Duration, 64)
		if err != nil {
			logger.Errorf("failed to parse duration: %s", err)

			_, stopErr := s.streams.Stop(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamID})
			if stopErr != nil {
				logger.Errorf("failed to stop stream: %s", err)
			}

			return
		}

		logger.Infof("duration is %f", duration)

		streamReq := &pstreamsv1.StreamRequest{
			Id:       req.StreamID,
			Duration: duration,
		}
		_, err = s.streams.Publish(context.Background(), streamReq)
		if err != nil {
			logger.Errorf("failed to publish stream: %s", err)

			_, stopErr := s.streams.Stop(context.Background(), &pstreamsv1.StreamRequest{Id: req.StreamID})
			if stopErr != nil {
				logger.Errorf("failed to stop stream: %s", err)
			}

			return
		}

	}(req.Filepath, localMediaDir)

	return &protoempty.Empty{}, nil
}

func (s *Server) splitMediafile(filepath string, hlsDir string) error {
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
	out, err := cmd.CombinedOutput()
	if err != nil {
		s.logger.Error(string(out))
		return err
	}

	return nil
}
