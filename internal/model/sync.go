package model

import "time"

type SyncArgs struct {
	Force bool
	All   bool
	Since time.Time
}
