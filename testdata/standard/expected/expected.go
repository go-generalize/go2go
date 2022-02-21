package types

import (
	"time"
)

type PostUserRequest struct {
	S string `json:"s"`
	T time.Time
}
