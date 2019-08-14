package client

import (
	"context"
)

var (
	// DefaultClient is a default client to use out of the box
	DefaultClient Client
	// DefaultBackoff is the default backoff function for retries
	DefaultBackoff = exponentialBackoff
	// DefaultRetry is the default check-for-retry function for retries
	DefaultRetry = RetryAlways
	// DefaultRetries is the default number of times a request is tried
	DefaultRetries = 1
)

// Call Makes a asynchronous call to a service using the default client
func Call(ctx context.Context, request Request, opts ...CallOption) error {
	return DefaultClient.Call(ctx, request, opts...)
}

// String string
func String() string {
	return DefaultClient.String()
}
