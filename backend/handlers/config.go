package handlers

import "time"

// Login rate limiting configuration
const (
	MaxLoginAttempts = 5
	LockoutDuration  = 15 * time.Minute
)
