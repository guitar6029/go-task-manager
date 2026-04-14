package model

import "encoding/json"

type Job struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
