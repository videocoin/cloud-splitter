package rpc

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/opentracing/opentracing-go"
	splitterv1 "github.com/videocoin/cloud-api/splitter/v1"
)

func (s *RpcServer) Split(ctx context.Context, req *splitterv1.SplitRequest) (*protoempty.Empty, error) {
	span := opentracing.SpanFromContext(ctx)
	span.SetTag("id", req.StreamID)

	localMediaDir := filepath.Join(s.hlsDir, req.StreamID)
	os.MkdirAll(localMediaDir, 0777)

	go func(filepath string, hlsDir string) {
		err := splitMediafile(filepath, hlsDir)
		if err != nil {
			s.logger.Errorf("failed to split media file: %s", err)
			return
		}
	}(req.Filepath, localMediaDir)

	return &protoempty.Empty{}, nil
}

func splitMediafile(filepath string, hlsDir string) error {
	output := fmt.Sprintf("%s/%%d.ts", hlsDir)
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
		"-bsf:v",
		"h264_mp4toannexb",
		"-map",
		"0",
		output,
	}
	cmd := exec.Command("ffmpeg", args...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
