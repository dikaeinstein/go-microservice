package throttling

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestHandler(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		<-r.Context().Done()
	})
}

func setup(ctx context.Context) (*http.Request, *httptest.ResponseRecorder) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(ctx)
	return r, rw
}

func TestReturnsBusyWhenZeroConnections(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	handler := NewLimitHandler(0, newTestHandler(ctx))
	r, rw := setup(ctx)

	time.AfterFunc(10*time.Millisecond, func() {
		cancel()
	})

	handler.ServeHTTP(rw, r)
	assert.Equal(t, rw.Code, http.StatusTooManyRequests)
}

func TestReturnsBusyWhenConnectionsExhausted(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	handler := NewLimitHandler(1, newTestHandler(ctx))
	r, rw := setup(ctx)
	r2, rw2 := setup(ctx2)

	time.AfterFunc(10*time.Millisecond, func() {
		cancel()
		cancel2()
	})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		handler.ServeHTTP(rw, r)
		wg.Done()
	}()
	go func() {
		handler.ServeHTTP(rw2, r2)
		wg.Done()
	}()
	wg.Wait()

	if rw.Code == http.StatusOK && rw2.Code == http.StatusOK {
		t.Fatalf("One request should have been busy, request 1: %v, request 2: %v",
			rw.Code, rw2.Code)
	}
}

func TestReleasesConnectionLockWhenFinished(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	handler := NewLimitHandler(1, newTestHandler(ctx))
	r, rw := setup(ctx)
	r2, rw2 := setup(ctx2)

	cancel()
	cancel2()

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(2)

	go func() {
		handler.ServeHTTP(rw, r)
		waitGroup.Done()
		handler.ServeHTTP(rw2, r2)
		waitGroup.Done()
	}()

	waitGroup.Wait()

	if rw.Code != http.StatusOK || rw2.Code != http.StatusOK {
		t.Fatalf("One request should have been busy, request 1: %v, request 2: %v", rw.Code, rw2.Code)
	}
}
