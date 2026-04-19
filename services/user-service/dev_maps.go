package main

import "sync"

// In-memory OTP / reset state when InsForge is disabled (dev only).
var (
	devMapsMu      sync.Mutex
	devVerifyCodes = map[string]string{}
	devResetCodes  = map[string]string{}
	devResetTokens = map[string]string{}
)
