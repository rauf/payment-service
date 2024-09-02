package nullutil

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/sqlc-dev/pqtype"
)

func NewNullString(s string) sql.NullString {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func NewNullRawMessage(m json.RawMessage) pqtype.NullRawMessage {
	if m == nil {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{RawMessage: m, Valid: true}
}
