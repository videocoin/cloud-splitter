// +build integration

package test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	protoempty "github.com/gogo/protobuf/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	splitterv1 "github.com/videocoin/cloud-api/splitter/v1"
	"github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/cloud-splitter/rpc"
	"github.com/videocoin/cloud-splitter/service"
)

var (
	ServiceName string = "splitter"
	Version     string = "test"
)

type APITestSuite struct {
	suite.Suite
	svc    *service.Service
	rpc    *rpc.Server
	config *service.Config
}

func (suite *APITestSuite) SetupSuite() {
	logger.Init(ServiceName, Version)

	log := logrus.NewEntry(logrus.New())
	suite.config = &service.Config{
		Name:    ServiceName,
		Version: Version,
	}

	err := envconfig.Process(ServiceName, suite.config)
	if err != nil {
		assert.FailNow(suite.T(), err.Error())
	}

	suite.config.Logger = log
	privateStreamManager := new(MockPrivateStreamManager)

	rpcConfig := &rpc.ServerOpts{
		HlsDir:      suite.config.HLSDir,
		Addr:        suite.config.RPCAddr,
		Logger:      suite.config.Logger.WithField("system", "rpc"),
		Streams:     privateStreamManager,
		SegmentTime: suite.config.SegmentTime,
	}

	suite.rpc, err = rpc.NewServer(rpcConfig)
	require.NoError(suite.T(), err)

}

func TestAPITestSuite(t *testing.T) {
	suite.Run(t, new(APITestSuite))
}

func (suite *APITestSuite) TestSplit() {
	req := &splitterv1.SplitRequest{Filepath: "testdata/small.mp4", StreamID: StreamID}
	span, _ := opentracing.StartSpanFromContext(context.Background(), "test")
	defer span.Finish()
	ctx := opentracing.ContextWithSpan(context.Background(), span)
	resp, err := suite.rpc.Split(ctx, req)

	if err != nil {
		assert.FailNow(suite.T(), err.Error())
	}
	assert.Equal(suite.T(), &protoempty.Empty{}, resp)
	if _, err := os.Stat(filepath.Join(suite.config.HLSDir, StreamID)); os.IsNotExist(err) {
		assert.FailNow(suite.T(), err.Error())
	}

}
