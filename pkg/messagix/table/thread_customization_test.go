package table_test

import (
	"testing"

	"go.mewis.me/meta-extra/pkg/messagix/lightspeed"
	"go.mewis.me/meta-extra/pkg/messagix/table"
)

func TestCurrentThreadCustomizationDependencies(t *testing.T) {
	dependencies := table.SPToDepMap([]string{"updateThreadParticipantNicknameV2", "updateThreadCustomEmoji"})
	result := &table.LSTable{}
	decoder := lightspeed.NewLightSpeedDecoder(dependencies, result)
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "updateThreadParticipantNicknameV2", []interface{}{float64(lightspeed.I64_FROM_STRING), "123"}, []interface{}{float64(lightspeed.I64_FROM_STRING), "456"}, "Mew", "mew"})
	decoder.Decode([]interface{}{float64(lightspeed.CALL_STORED_PROCEDURE), "updateThreadCustomEmoji", []interface{}{float64(lightspeed.I64_FROM_STRING), "123"}, "😀", []interface{}{float64(lightspeed.UNDEFINED)}})
	if len(result.LSUpdateThreadParticipantNicknameV2) != 1 {
		t.Fatalf("unexpected nickname decode: %#v", result.LSUpdateThreadParticipantNicknameV2)
	}
	nickname := result.LSUpdateThreadParticipantNicknameV2[0]
	if nickname.ThreadKey != 123 || nickname.ParticipantId != 456 || nickname.Nickname != "Mew" || nickname.NormalizedNickname != "mew" {
		t.Fatalf("unexpected nickname dependency: %#v", nickname)
	}
	if len(result.LSUpdateThreadCustomEmoji) != 1 {
		t.Fatalf("unexpected emoji decode: %#v", result.LSUpdateThreadCustomEmoji)
	}
	emoji := result.LSUpdateThreadCustomEmoji[0]
	if emoji.ThreadKey != 123 || emoji.CustomEmoji != "😀" || emoji.CustomEmojiImageUrl != "" {
		t.Fatalf("unexpected emoji dependency: %#v", emoji)
	}
}
