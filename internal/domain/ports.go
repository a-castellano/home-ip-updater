package domain

import (
	"context"
)

type RecordUpdater interface {
	UpdateRecord(ctx context.Context, value string) error
}
