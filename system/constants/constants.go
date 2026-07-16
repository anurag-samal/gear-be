package constants

import "time"

const (
	PG_CONNECTION_TIMEOUT = 10 * time.Second
	CONTEXT_TIMEOUT       = 10 * time.Second
	SRV_READ_TIME_OUT     = 10 * time.Second
	SRV_WRITE_TIME_OUT    = 10 * time.Second
	SRV_IDLE_TIME_OUT     = 10 * time.Minute
	SRV_ADDRESS           = ":8000"
)
