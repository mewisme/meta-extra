package messagix

import (
	"encoding/json"
	"testing"
)

func TestThreadThemeDecodesCompleteShape(t *testing.T) {
	raw := []byte(`{"id":"1","accessibility_label":"Theme","description":"desc","app_color_mode":"DARK","background_asset":{"id":"2","image":{"uri":"background"}},"background_gradient_colors":["a"],"composer_background_color":"b","composer_input_background_color":"c","composer_tint_color":"d","fallback_color":"e","gradient_colors":["f"],"hot_like_color":"g","icon_asset":{"id":"3","image":{"uri":"icon"}},"inbound_message_border_color":"h","inbound_message_border_width":1,"inbound_message_gradient_colors":["i"],"inbound_message_text_color":"j","is_deprecated":false,"message_border_color":"k","message_border_width":2,"message_text_color":"l","normal_theme_id":"4","primary_button_background_color":"m","reaction_pill_background_color":"n","reverse_gradients_for_radial":true,"secondary_text_color":"o","tertiary_text_color":"p","title_bar_attribution_color":"q","title_bar_background_color":"r","title_bar_button_tint_color":"s","title_bar_text_color":"t","alternative_themes":[{"id":"5","accessibility_label":"Alternative","description":null,"app_color_mode":"LIGHT","background_asset":null,"background_gradient_colors":[],"composer_background_color":null,"composer_input_background_color":null,"composer_tint_color":null,"fallback_color":"u","gradient_colors":[],"hot_like_color":null,"icon_asset":null,"inbound_message_border_color":null,"inbound_message_border_width":null,"inbound_message_gradient_colors":[],"inbound_message_text_color":null,"is_deprecated":true,"message_border_color":null,"message_border_width":null,"message_text_color":null,"normal_theme_id":"6","primary_button_background_color":null,"reaction_pill_background_color":null,"reverse_gradients_for_radial":false,"secondary_text_color":null,"tertiary_text_color":null,"title_bar_attribution_color":null,"title_bar_background_color":null,"title_bar_button_tint_color":null,"title_bar_text_color":null}]}`)
	var theme ThreadTheme
	if err := json.Unmarshal(raw, &theme); err != nil {
		t.Fatal(err)
	}
	if theme.ID != "1" || theme.AccessibilityLabel != "Theme" || theme.BackgroundAsset == nil || theme.BackgroundAsset.ID != "2" || theme.BackgroundAsset.Image.URI != "background" {
		t.Fatalf("unexpected primary theme: %+v", theme)
	}
	if theme.InboundMessageBorderWidth == nil || *theme.InboundMessageBorderWidth != 1 || theme.MessageBorderWidth == nil || *theme.MessageBorderWidth != 2 || !theme.ReverseGradientsForRadial {
		t.Fatalf("missing primary theme fields: %+v", theme)
	}
	if len(theme.AlternativeThemes) != 1 || theme.AlternativeThemes[0].ID != "5" || !theme.AlternativeThemes[0].IsDeprecated || theme.AlternativeThemes[0].Description != nil {
		t.Fatalf("unexpected alternative theme: %+v", theme.AlternativeThemes)
	}
}
