package conduit

import (
	"net/http"

	"golang.org/x/time/rate"
)

func MakeTransportWithRateLimit(t http.RoundTripper, rps int) http.RoundTripper {
	return &rateLimitedTransport{
		RoundTripper: t,
		limiter:      rate.NewLimiter(rate.Limit(rps), 1),
	}
}

type rateLimitedTransport struct {
	http.RoundTripper

	limiter *rate.Limiter
}

func (t rateLimitedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(r.Context()); err != nil {
		return &http.Response{
			Status:     "429 Too Many Requests",
			StatusCode: http.StatusTooManyRequests,
			Request:    r,
			Body:       http.NoBody,
			Header:     make(http.Header),
		}, nil
	}
	return t.RoundTripper.RoundTrip(r)
}
