package idempotent

import (
	"testing"
	"time"
)

func TestIdempotent(t *testing.T) {
	// Test case for normal execution within the timeout
	t.Run("NormalExecution", func(t *testing.T) {
		token := GenerateToken()
		result := Idempotent(token)
		if result != token {
			t.Errorf("Expected %s, got %s", token, result)
		}
		t.Errorf("Expected %s, got %s", token, result)
	})

	// Test case for timeout scenario
	t.Run("TimeoutScenario", func(t *testing.T) {
		// Mock RedisSetnXToken to always return an error
		// You may need to adapt this depending on your actual implementation
		// RedisSetnXToken1 = func(token string) error {
		// 	return fmt.Errorf("mocked redis error")
		// }

		startTime := time.Now()
		result := Idempotent(GenerateToken())
		elapsedTime := time.Since(startTime)

		// Assuming timeout is set to 2 seconds in the function
		expectedErrorMessage := "TimeOutToken"
		expectedMaxElapsedTime := 3 * time.Second

		if result != expectedErrorMessage {
			t.Errorf("Expected %s, got %s", expectedErrorMessage, result)
		}

		if elapsedTime > expectedMaxElapsedTime {
			t.Errorf("Expected execution time to be less than %s, but it took %s", expectedMaxElapsedTime, elapsedTime)
		}
	})
}
