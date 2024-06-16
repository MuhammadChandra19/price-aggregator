package ws

import (
	"context"

	"github.com/muhammadchandra19/price-aggregator/market-data-fetcher/module"
)

// MockFeed is a mock struct for a feed
type MockFeed struct {
	stopListener      context.CancelFunc
	expectedSubscribe error
	expectedTurnOff   error
	expectedRead      *module.Trade
	expectedReadError error
	allTrades         []module.Trade
}

// NewMockFeed creates a new MockFeed instance
func NewMockFeed(stopListener context.CancelFunc) *MockFeed {
	return &MockFeed{stopListener: stopListener}
}

// SetExpectedSubscribe sets the expected return value for the Subscribe method
func (m *MockFeed) SetExpectedSubscribe(err error) {
	m.expectedSubscribe = err
}

// SetExpectedTurnOff sets the expected return value for the TurnOff method
func (m *MockFeed) SetExpectedTurnOff(err error) {
	m.expectedTurnOff = err
}

// SetExpectedRead sets the expected return value and error for the Read method
func (m *MockFeed) SetExpectedRead(trade *module.Trade, err error) {
	m.expectedRead = trade
	m.expectedReadError = err
}

func (m *MockFeed) SetTrades(trades []module.Trade) {
	m.allTrades = trades
}

// Subscribe is the mock implementation of the Subscribe method
func (m *MockFeed) Subscribe() error {
	return m.expectedSubscribe
}

// TurnOff is the mock implementation of the TurnOff method
func (m *MockFeed) TurnOff() error {
	return m.expectedTurnOff
}

// Read is the mock implementation of the Read method
func (m *MockFeed) Read() (*module.Trade, error) {
	if count := len(m.allTrades); count <= 0 {
		m.stopListener()
		return &module.Trade{}, nil
	}
	if m.expectedReadError != nil {
		m.stopListener()
	}

	op := m.allTrades[0]
	m.allTrades = m.allTrades[1:]
	return &op, m.expectedReadError
}
