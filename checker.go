package healthcheck

// Checker is an interface that handles health check status
type Checker interface {
	Check() error
}
