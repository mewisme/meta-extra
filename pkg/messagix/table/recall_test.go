package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestCleanUpOnRecallDependency(t *testing.T) {
	const messageID = "mid.$phase1-recall"
	dependencies := table.SPToDepMap([]string{"cleanUpOnRecall"})
	if dependencies["cleanUpOnRecall"] != "LSCleanUpOnRecall" {
		t.Fatalf("unexpected dependency mapping: %q", dependencies["cleanUpOnRecall"])
	}

	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "cleanUpOnRecall", messageID})
	if len(result.LSCleanUpOnRecall) != 1 || result.LSCleanUpOnRecall[0].MessageId != messageID {
		t.Fatalf("unexpected recall decode: %#v", result.LSCleanUpOnRecall)
	}
}
