package application

import (
	"card-validator-service/internal/core/domain"

	"github.com/go-playground/validator/v10"
	"github.com/mwinyimoha/card-validator-utils/errors"
)

type ValidatorService interface {
	ValidateCardNumber(payload *domain.CardNumberPayload) (*domain.CardInfo, error)
}

type Service struct {
	Validator *validator.Validate
}

func NewService(v *validator.Validate) ValidatorService {
	return &Service{Validator: v}
}

func (s *Service) ValidateCardNumber(payload *domain.CardNumberPayload) (*domain.CardInfo, error) {
	if err := s.Validator.Struct(payload); err != nil {
		if verr, ok := err.(validator.ValidationErrors); ok {
			violations := errors.BuildViolations(verr)
			return nil, errors.NewValidationError(violations)
		}

		return nil, err
	}

	return domain.NewCardInfo(payload.CardNumber), nil
}
