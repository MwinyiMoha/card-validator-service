package api

import (
	"card-validator-service/internal/core/application"
	"card-validator-service/internal/core/domain"
	protos "card-validator-service/internal/gen"
	"context"

	"github.com/mwinyimoha/card-validator-utils/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	protos.UnimplementedCardValidatorServiceServer
	Service application.ValidatorService
}

func NewServer(svc application.ValidatorService) *Server {
	return &Server{Service: svc}
}

func (s *Server) ValidateNumber(ctx context.Context, req *protos.ValidateNumberRequest) (*protos.ValidateNumberResponse, error) {
	payload := domain.CardNumberPayload{
		CardNumber: req.GetCardNumber(),
	}

	cardInfo, err := s.Service.ValidateCardNumber(&payload)
	if err != nil {
		if verr, ok := err.(*errors.ValidationError); ok {
			st, _ := verr.GRPCStatus() // TODO: Handle error
			return nil, st.Err()
		}

		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &protos.ValidateNumberResponse{
		Valid: true,
		Data: &protos.CardData{
			CardNumber: req.CardNumber,
			ProviderInformation: &protos.ProviderInformation{
				Name:    cardInfo.ProviderInformation.Name,
				IconUrl: cardInfo.ProviderInformation.IconURL,
			},
		},
	}, nil
}
