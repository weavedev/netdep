package natsanalyzer

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"lab.weave.nl/internships/tud-2022/netDep/helpers"
)

func TestFindNATSCalls(t *testing.T) {
	svcDir := filepath.Join(helpers.RootDir, "test", "sample", "nats", "svc")
	consumers, producers, _ := FindNATSCalls(svcDir)

	assert.Equal(t, len(consumers), 1)
	assert.Equal(t, len(producers), 1)
	assert.Equal(t, consumers[0].Subject, "SnapshotStartdateChangedSubject")
	assert.Equal(t, producers[0].Subject, "SnapshotStartdateChangedSubject")
	assert.Equal(t, producers[0].ServiceName, "snapshot")
}
