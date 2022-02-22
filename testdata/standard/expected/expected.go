package types

import (
	"time"
)

type PostUserRequest struct {
	Map   map[string]int64
	S     string `json:"s"`
	Slice []string
	T     time.Time
}
