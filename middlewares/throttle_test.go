package middlewares

import (
	"context"
	"testing"
	"time"
)

func TestSlidingLogThrottle(t *testing.T) {
	clientID := "clientA"
	endpoint := "foo"
	requestsPerWindow := 3
	windowDuration := time.Minute

	throttle := NewSlidingLogThrottle(windowDuration)
	for i := 0; i < requestsPerWindow; i++ {
		allowed, err := throttle.Allow(context.Background(), clientID, endpoint, requestsPerWindow)
		if err != nil {
			t.Fatal(err)
		}
		if !allowed {
			t.Fatalf("Expected allowed request %d", i+1)
		}
	}

	allowed, err := throttle.Allow(context.Background(), clientID, endpoint, requestsPerWindow)
	if err != nil {
		t.Fatal(err)
	}
	if allowed {
		t.Fatal("Expected rate limit exceeded")
	}
}
