package eventbus

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	"github.com/videocoin/cloud-pkg/mqmux"
	tracerext "github.com/videocoin/cloud-pkg/tracer"
)

type Config struct {
	Logger *logrus.Entry
	URI    string
	Name   string
	HlsDir string
}

type EventBus struct {
	logger *logrus.Entry
	mq     *mqmux.WorkerMux
	hlsDir string
}

func New(c *Config) (*EventBus, error) {
	mq, err := mqmux.NewWorkerMux(c.URI, c.Name)
	if err != nil {
		return nil, err
	}
	if c.Logger != nil {
		mq.Logger = c.Logger
	}

	return &EventBus{
		logger: c.Logger,
		mq:     mq,
		hlsDir: c.HlsDir,
	}, nil
}

func (e *EventBus) Start() error {
	err := e.mq.Consumer("streams.completed", 1, false, e.handleStreamCompleted)
	if err != nil {
		return err
	}
	return e.mq.Run()
}

func (e *EventBus) Stop() error {
	return e.mq.Close()
}

func (e *EventBus) handleStreamCompleted(d amqp.Delivery) error {
	var span opentracing.Span
	tracer := opentracing.GlobalTracer()
	spanCtx, err := tracer.Extract(opentracing.TextMap, mqmux.RMQHeaderCarrier(d.Headers))

	e.logger.Debugf("handling body: %+v", string(d.Body))

	if err != nil {
		span = tracer.StartSpan("eventbus.handleStreamEvent")
	} else {
		span = tracer.StartSpan("eventbus.handleStreamEvent", ext.RPCServerOption(spanCtx))
	}

	defer span.Finish()

	req := new(pstreamsv1.Event)
	err = json.Unmarshal(d.Body, req)
	if err != nil {
		tracerext.SpanLogError(span, err)
		return err
	}

	span.SetTag("stream_id", req.StreamID)
	span.SetTag("event_type", req.Type.String())

	e.logger.Debugf("handling request %+v", req)

	switch req.Type {
	case pstreamsv1.EventTypeUpdate:
		{
			e.logger.Info("cleanup directory")
			StreamDir := filepath.Join(e.hlsDir, req.StreamID)
			err = os.RemoveAll(StreamDir)
			if err != nil {
				tracerext.SpanLogError(span, err)
				return err
			}

		}
	case pstreamsv1.EventTypeUnknown:
		e.logger.Error("event type is unknown")
	}

	return nil
}
