package connector

import (
	"testing"

	"maunium.net/go/mautrix/event"
)

func TestLocationCapabilityMatchesTransport(t *testing.T) {
	if metaCaps.LocationMessage != event.CapLevelUnsupported {
		t.Fatalf("regular Messenger unexpectedly advertises location sending: %v", metaCaps.LocationMessage)
	}
	if metaCapsWithE2E.LocationMessage != event.CapLevelFullySupported || metaCapsWithE2EGroup.LocationMessage != event.CapLevelFullySupported {
		t.Fatalf("E2EE Messenger location capability not advertised")
	}
}
