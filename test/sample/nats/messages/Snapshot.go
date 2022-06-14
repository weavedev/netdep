//nolint
package messages

import "github.com/gofrs/uuid"

// NewSnapshotStartdateChangedNotifyMsg  will create a notify message with an ID, subject and data
func NewSnapshotStartdateChangedNotifyMsg(snapshotID uuid.UUID, subject string) (string, error) {
	response := "reply mock"

	return response, nil
}
