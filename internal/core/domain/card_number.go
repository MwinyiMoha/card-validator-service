package domain

import (
	"fmt"
	"strings"
)

type CardNumberPayload struct {
	CardNumber string `validate:"required,valid_card_number"`
}

type CardProviderInfo struct {
	Name    string
	IconURL string
}

type CardInfo struct {
	CardNumber          string
	ProviderInformation *CardProviderInfo
}

func NewCardInfo(cardNumber string) *CardInfo {
	providers := map[string]string{
		"3": "AMEX",
		"4": "VISA",
		"5": "MASTERCARD",
		"6": "DISCOVER",
	}
	providerName := providers[cardNumber[:1]]

	return &CardInfo{
		CardNumber: cardNumber,
		ProviderInformation: &CardProviderInfo{
			Name:    providerName,
			IconURL: fmt.Sprintf("https://example.dummy/card-provider-icons/%s.png", strings.ToLower(providerName)),
		},
	}
}
