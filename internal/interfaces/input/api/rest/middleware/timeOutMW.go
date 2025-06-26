package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// function TimeoutMiddleware accepts two args 1.-> timeout duration  2. ->httpHandler
func TimeoutMiddleware(timeout time.Duration, next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// create a new context ctx which would be automatically cancelled after timeout
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel() //to clean up resources

		// context is attached to the request so downnstream handler can also see the timeout
		r = r.WithContext(ctx)

		// run the handler in go routine
		done := make(chan struct{})
		go func() {
			t := time.Now()
			next.ServeHTTP(w, r)
			// next handler is run in a seperate go routine
			fmt.Printf("Processing time : %v\n", time.Since(t))
			close(done)
			// when the handler finishes it closes the done channel
		}()

		// wait for completion or timeout
		select {
		case <-done:
			// request completed successfully
		case <-ctx.Done():
			// request not completed and deadline achieved
			if ctx.Err() == context.DeadlineExceeded {
				http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			}
		}
	})
}
