package updater

import (
	"context"
	"errors"
	"testing"
)

type mockUpdater struct {
	fail bool
}

func (m mockUpdater) UpdateRecord(ctx context.Context, value string) error {
	if m.fail {
		return errors.New("Fail")
	}
	return nil
}

func TestUpdater(t *testing.T) {
	ctx := context.Background()

	mock := mockUpdater{fail: true}

	updater := NewUpdater(mock)

	err := updater.UpdateIP(ctx, "anyvalue")

	if err == nil {
		t.Fatalf("TestUpdater UpdateIP test should falaar as mockUpdater fails to")
	}
}
