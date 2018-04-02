package foundation

import "time"

// InternalTimer internal holded instance
var InternalTimer = NewLocalTimer()

// LocalTimer local struct will be wrapp time.Time
type LocalTimer struct{}

// NewLocalTimer return new LocalTimer instance
func NewLocalTimer() Timer {
	return &LocalTimer{}
}

// Timer interface to wrap time package
type Timer interface {
	Now() time.Time
}

// Now return when used in this app
func (lt *LocalTimer) Now() time.Time {
	return time.Now()
}
