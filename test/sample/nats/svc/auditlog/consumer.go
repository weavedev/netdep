//nolint
package main

import (
	"fmt"

	natsconfig "lab.weave.nl/internships/tud-2022/netDep/test/sample/nats"
)

func main() {
	observant := &natsconfig.Observant{
		ID: "1",
	}

	fmt.Println("The tests are running just fine")

	observant.Subscribe("consumer", natsconfig.SnapshotStartdateChangedSubject)
}
