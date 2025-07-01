package event

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type DbEvent struct {
	ID             uuid.UUID       `json:"id"`
	Timestamp      time.Time       `json:"timestamp"`
	Action         string          `json:"action"`
	DbSchema       string          `json:"db_schema"`
	CollectionName string          `json:"table"`
	Record         json.RawMessage `json:"record"`
	Old            json.RawMessage `json:"old"`
}
