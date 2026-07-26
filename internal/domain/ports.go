package domain

import (
	"context"
)

type RecordUpdater interface {
	UpdateRecord(ctx context.Context, record string, value string) error
}
