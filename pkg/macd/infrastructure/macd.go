package infrastructure

import (
	"errors"
	"sync"
	"time"

	"github.com/daddydemir/crypto/pkg/macd/domain"
)

// InMemoryMACDRepository implements MACDRepository using in-memory storage
type InMemoryMACDRepository struct {
	data  map[string][]*domain.MACDResult
	mutex sync.RWMutex
}

func NewInMemoryMACDRepository() *InMemoryMACDRepository {
	return &InMemoryMACDRepository{
		data: make(map[string][]*domain.MACDResult),
	}
}

func (r *InMemoryMACDRepository) SaveMACDResult(symbol string, result *domain.MACDResult) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.data[symbol] == nil {
		r.data[symbol] = make([]*domain.MACDResult, 0)
	}

	// Add timestamp to distinguish between results
	resultCopy := *result
	r.data[symbol] = append(r.data[symbol], &resultCopy)

	// Keep only last 100 results per symbol
	if len(r.data[symbol]) > 100 {
		r.data[symbol] = r.data[symbol][len(r.data[symbol])-100:]
	}

	return nil
}

func (r *InMemoryMACDRepository) GetMACDResult(symbol string, from, to time.Time) (*domain.MACDResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	results, exists := r.data[symbol]
	if !exists || len(results) == 0 {
		return nil, errors.New("no MACD data found for symbol")
	}

	// Return the most recent result for now
	// In a real implementation, you'd filter by date range
	return results[len(results)-1], nil
}

func (r *InMemoryMACDRepository) GetLatestMACDResult(symbol string) (*domain.MACDResult, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	results, exists := r.data[symbol]
	if !exists || len(results) == 0 {
		return nil, errors.New("no MACD data found for symbol")
	}

	return results[len(results)-1], nil
}

func (r *InMemoryMACDRepository) DeleteOldMACDData(olderThan time.Time) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// For in-memory implementation, this is a no-op
	// In a real database implementation, you'd delete records older than the specified time
	return nil
}
