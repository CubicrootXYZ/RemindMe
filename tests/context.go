package tests

import (
	"context"
	"time"
)

func ContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*2)
}
