package domain

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func NewUUIDV7() pgtype.UUID {
	id := uuid.Must(uuid.NewV7())
	var u pgtype.UUID
	copy(u.Bytes[:], id[:])
	u.Valid = true
	return u
}

func UUIDToString(u pgtype.UUID) string {
	if !u.Valid {
		return ""
	}
	var src [16]byte
	copy(src[:], u.Bytes[:])
	return fmt.Sprintf("%x-%x-%x-%x-%x", src[0:4], src[4:6], src[6:8], src[8:10], src[10:16])
}

func StringToUUID(s string) pgtype.UUID {
	var uuid pgtype.UUID
	_ = uuid.Scan(s)
	return uuid
}
