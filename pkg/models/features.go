package models

import (
	"encoding/json"
	"time"
)

//go:generate reform
//reform:features
type Features struct {
	ID          string          `reform:"id,pk"`
	Key         string          `reform:"key"`
	Description *string         `reform:"description"`
	Enabled     bool            `reform:"enabled"`
	Value       json.RawMessage `reform:"value"`
	CreatedAt   time.Time       `reform:"created_at"`
	UpdatedAt   time.Time       `reform:"updated_at"`
}
