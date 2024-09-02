package backoff

type RetryConfig struct {
	MaxRetries int
	Backoff    Strategy
}
