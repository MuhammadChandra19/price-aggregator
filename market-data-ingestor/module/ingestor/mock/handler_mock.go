// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go

// Package ingestor is a generated GoMock package.
package ingestor

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	module "github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
)

// MockIngestorHandler is a mock of IngestorHandler interface.
type MockIngestorHandler struct {
	ctrl     *gomock.Controller
	recorder *MockIngestorHandlerMockRecorder
}

// MockIngestorHandlerMockRecorder is the mock recorder for MockIngestorHandler.
type MockIngestorHandlerMockRecorder struct {
	mock *MockIngestorHandler
}

// NewMockIngestorHandler creates a new mock instance.
func NewMockIngestorHandler(ctrl *gomock.Controller) *MockIngestorHandler {
	mock := &MockIngestorHandler{ctrl: ctrl}
	mock.recorder = &MockIngestorHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIngestorHandler) EXPECT() *MockIngestorHandlerMockRecorder {
	return m.recorder
}

// WriteMarketFeed mocks base method.
func (m *MockIngestorHandler) WriteMarketFeed(ctx context.Context, feeds map[string]*module.Trade) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WriteMarketFeed", ctx, feeds)
}

// WriteMarketFeed indicates an expected call of WriteMarketFeed.
func (mr *MockIngestorHandlerMockRecorder) WriteMarketFeed(ctx, feeds interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMarketFeed", reflect.TypeOf((*MockIngestorHandler)(nil).WriteMarketFeed), ctx, feeds)
}
