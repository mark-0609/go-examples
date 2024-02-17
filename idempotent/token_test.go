package idempotent

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

// MockRedisClient is a mock implementation of the Redis client.
type MockRedisClient struct {
	data map[string]string
}

func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	val, _ := m.data[key]
	return redis.NewStringResult(val, redis.Nil)
}

func (m *MockRedisClient) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	_, exists := m.data[key]
	m.data[key] = value.(string)
	return redis.NewBoolResult(!exists, nil)
}

func TestProcessRequest(t *testing.T) {
	// Mock the Redis client
	mockRedisClient := &MockRedisClient{
		data: make(map[string]string),
	}
	// redisClient = mockRedisClient // Set the global redisClient variable to use the mock

	// Test case for an existing result in Redis
	t.Run("ExistingResultInRedis", func(t *testing.T) {
		// Set up initial data in Redis
		requestID := "existing_request"
		expectedResult := "Existing result in Redis"
		mockRedisClient.data[requestID] = expectedResult

		// Call the function
		result, err := ProcessRequest(requestID)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
	})

	// Test case for processing a new request
	t.Run("NewRequestProcessing", func(t *testing.T) {
		// Set up test data
		requestID := GenerateToken()

		// Call the function
		result, err := ProcessRequest(requestID)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, "success result", result)
		assert.Equal(t, "success result", mockRedisClient.data[requestID])
	})

	// Test case for timeout during processing
	t.Run("TimeoutDuringProcessing", func(t *testing.T) {
		// Set up test data
		requestID := GenerateToken()

		// Mock the Work function to simulate a timeout
		// originalWork := Work
		// Work = func(requestID string) (string, error) {
		// 	time.Sleep(time.Second * 3)
		// 	return "", errors.New("timeout")
		// }
		// defer func() {
		// 	// Restore the original Work function
		// 	Work = originalWork
		// }()

		// Call the function
		result, err := ProcessRequest(requestID)

		// Assertions
		assert.Error(t, err)
		assert.Equal(t, "", result)
	})
}
