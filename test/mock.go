package test

import (
	"context"

	pstreamsv1 "github.com/videocoin/cloud-api/streams/private/v1"
	streamsv1 "github.com/videocoin/cloud-api/streams/v1"
	"google.golang.org/grpc"
)

const UserID = "12b1876f-341f-41b0-833f-5312f1e9c308"
const StreamID = "cdc1816b-0be8-44a6-80c3-3e43fbd441ee"

type MockPrivateStreamManager struct {
	id string
}

func (sm *MockPrivateStreamManager) Get(
	context.Context, *pstreamsv1.StreamRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID:     sm.id,
		UserID: UserID,
		Status: streamsv1.StreamStatusPrepared,
	}
	return &stream, nil
}

func (sm *MockPrivateStreamManager) Publish(
	context.Context, *pstreamsv1.StreamRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID: sm.id,
	}
	return &stream, nil
}

func (sm *MockPrivateStreamManager) PublishDone(
	context.Context, *pstreamsv1.StreamRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID: sm.id,
	}
	return &stream, nil
}

func (sm *MockPrivateStreamManager) Run(
	context.Context, *pstreamsv1.StreamRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID: sm.id,
	}
	return &stream, nil
}

func (sm *MockPrivateStreamManager) Stop(
	context.Context, *pstreamsv1.StreamRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID: sm.id,
	}
	return &stream, nil
}

func (sm *MockPrivateStreamManager) UpdateStatus(
	context.Context, *pstreamsv1.UpdateStatusRequest, ...grpc.CallOption,
) (*pstreamsv1.StreamResponse, error) {
	stream := pstreamsv1.StreamResponse{
		ID: sm.id,
	}
	return &stream, nil
}
