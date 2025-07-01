package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
	"gitlab.com/enCapital/models"
	"gitlab.com/enCapital/models/dnse"
	"google.golang.org/protobuf/types/known/timestamppb"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/suggested_offer/repository"
	"financing-offer/internal/event"
	"financing-offer/pkg/pb"
)

var _ repository.SuggestedOfferEventRepository = (*SuggestedOfferEventPublisher)(nil)

type SuggestedOfferEventPublisher struct {
	config    config.KafkaConfig
	publisher event.Publisher
}

func NewSuggestedOfferEventPublisher(config config.KafkaConfig, publisher event.Publisher) *SuggestedOfferEventPublisher {
	return &SuggestedOfferEventPublisher{
		config:    config,
		publisher: publisher,
	}
}

func (p *SuggestedOfferEventPublisher) NotifySuggestedOfferCreated(
	_ context.Context,
	investorId string,
	config entity.SuggestedOfferConfig,
	createdOffer entity.SuggestedOffer,
) error {
	errorTemplate := "SuggestedOfferEventPublisher NotifySuggestedOfferCreated %w"
	title := ""
	if config.ValueType == entity.ValueTypeInterestRate {
		title = fmt.Sprintf("Phản hồi đề cử Margin %s", config.Value.Mul(decimal.NewFromInt(100)).String()+"%/năm")
	}
	payload := dnse.FinancingOfferSuggestedOfferCreated{
		InvestorId: investorId,
		Title:      title,
		AccountNo:  createdOffer.AccountNo,
		Symbols:    createdOffer.Symbols,
		CreatedAt:  timestamppb.New(createdOffer.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_SUGGESTED_OFFER_CREATED)
	if err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	if err = p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(investorId),
		},
	); err != nil {
		return fmt.Errorf(errorTemplate, err)
	}
	return nil
}
