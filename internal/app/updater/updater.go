package updater

import (
	"context"
	domain "github.com/a-castellano/home-ip-updater/internal/domain"
)

type Updater struct {
	recordUpdater domain.RecordUpdater
}

func NewUpdater(recordUpdater domain.RecordUpdater) *Updater {
	return &Updater{recordUpdater: recordUpdater}
}

func (u *Updater) UpdateIP(ctx context.Context, ipToUpdate string) error {
	return u.recordUpdater.UpdateRecord(ctx, ipToUpdate)
}
