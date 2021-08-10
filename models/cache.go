package models

import "time"

type CacheItem struct {
	ID        int64       `json:"id"`
	FieldName string      `json:"field_name"`
	OldValue  interface{} `json:"old_value, omitempty"`
	NewValue  interface{} `json:"new_value"`
	UpdatedAt time.Time   `json:"updated_at"`
}
