package util

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func Now() time.Time {
	return time.Now()
}

// Makes a pgtype.Timestamptz from a time.Time
func MakePgTimestamp(timestamp time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:             timestamp,
		InfinityModifier: pgtype.Finite,
		Valid:            true,
	}
}

func MakePgUuid(uuid uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: uuid,
		Valid: true,
	}
}
