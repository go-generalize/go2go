package standard

import "time"

type Status int

const (
	StatusOK      = "OK"
	StatusFailure = "Failure"
)

type PostUserRequest struct {
	T     time.Time
	S     string `json:"s"`
	Slice []string
	Map   map[string]int64
}
