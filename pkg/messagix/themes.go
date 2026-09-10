package messagix

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mewis.me/meta-extra/pkg/messagix/types"
)

type ThreadThemeAsset struct {
	ID    string `json:"id"`
	Image struct {
		URI string `json:"uri"`
	} `json:"image"`
}

type ThreadThemeVariant struct {
	ID                           string            `json:"id"`
	AccessibilityLabel           string            `json:"accessibility_label"`
	Description                  *string           `json:"description"`
	AppColorMode                 string            `json:"app_color_mode"`
	BackgroundAsset              *ThreadThemeAsset `json:"background_asset"`
	BackgroundGradientColors     []string          `json:"background_gradient_colors"`
	ComposerBackgroundColor      *string           `json:"composer_background_color"`
	ComposerInputBackgroundColor *string           `json:"composer_input_background_color"`
	ComposerTintColor            *string           `json:"composer_tint_color"`
	FallbackColor                string            `json:"fallback_color"`
	GradientColors               []string          `json:"gradient_colors"`
	HotLikeColor                 *string           `json:"hot_like_color"`
	IconAsset                    *ThreadThemeAsset `json:"icon_asset"`
	InboundMessageBorderColor    *string           `json:"inbound_message_border_color"`
	InboundMessageBorderWidth    *int64            `json:"inbound_message_border_width"`
	InboundMessageGradientColors []string          `json:"inbound_message_gradient_colors"`
	InboundMessageTextColor      *string           `json:"inbound_message_text_color"`
	IsDeprecated                 bool              `json:"is_deprecated"`
	MessageBorderColor           *string           `json:"message_border_color"`
	MessageBorderWidth           *int64            `json:"message_border_width"`
	MessageTextColor             *string           `json:"message_text_color"`
	NormalThemeID                string            `json:"normal_theme_id"`
	PrimaryButtonBackgroundColor *string           `json:"primary_button_background_color"`
	ReactionPillBackgroundColor  *string           `json:"reaction_pill_background_color"`
	ReverseGradientsForRadial    bool              `json:"reverse_gradients_for_radial"`
	SecondaryTextColor           *string           `json:"secondary_text_color"`
	TertiaryTextColor            *string           `json:"tertiary_text_color"`
	TitleBarAttributionColor     *string           `json:"title_bar_attribution_color"`
	TitleBarBackgroundColor      *string           `json:"title_bar_background_color"`
	TitleBarButtonTintColor      *string           `json:"title_bar_button_tint_color"`
	TitleBarTextColor            *string           `json:"title_bar_text_color"`
}

type ThreadTheme struct {
	ThreadThemeVariant
	AlternativeThemes []ThreadThemeVariant `json:"alternative_themes"`
}

func (c *Client) ListThreadThemes(ctx context.Context) ([]ThreadTheme, error) {
	if c == nil {
		return nil, ErrClientIsNil
	} else if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	_, body, err := c.http.MakeGraphQLRequest(ctx, "MWPThreadThemeQuery_AllThemesQuery", map[string]any{"version": "default"})
	if err != nil {
		return nil, err
	}
	var response struct {
		Data struct {
			Themes []ThreadTheme `json:"messenger_thread_themes"`
		} `json:"data"`
		types.ErrorResponse
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to decode thread themes: %w", err)
	} else if err := response.AsError(); err != nil {
		return nil, err
	}
	return response.Data.Themes, nil
}
