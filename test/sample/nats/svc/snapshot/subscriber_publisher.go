//nolint
package main

import (
	"github.com/gofrs/uuid"
	natsconfig "lab.weave.nl/internships/tud-2022/netDep/test/sample/nats"
	"lab.weave.nl/internships/tud-2022/netDep/test/sample/nats/messages"
)

func main() {
	messages.NewSnapshotStartdateChangedNotifyMsg(uuid.Must(uuid.NewV4()), natsconfig.SnapshotStartdateChangedSubject)
}
