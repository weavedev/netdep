//nolint
package main

import natsconfig "lab.weave.nl/internships/tud-2022/static-analysis-project/test/sample/nats"

func main() {
	observant := &natsconfig.Observant{
		ID: "1",
	}

	observant.Subscribe("consumer", natsconfig.SnapshotStartdateChangedSubject)
}
