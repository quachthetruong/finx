package kafka

import (
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	testifyMock "github.com/stretchr/testify/mock"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/test/mock"
)

func TestSuggestedOfferEvent_NotifySuggestedOfferCreated(t *testing.T) {
	t.Parallel()

	t.Run("Publish success", func(t *testing.T) {
		kafkaPublisher := mock.NewMockPublisher(t)
		publisher := NewSuggestedOfferEventPublisher(
			config.KafkaConfig{
				NotificationTopic: "notification",
			}, kafkaPublisher,
		)
		kafkaPublisher.EXPECT().Publish(testifyMock.Anything).Return(nil)
		err := publisher.NotifySuggestedOfferCreated(context.Background(), "investorId", entity.SuggestedOfferConfig{
			ValueType: entity.ValueTypeInterestRate,
			Value:     decimal.NewFromFloat(0.12),
		}, entity.SuggestedOffer{})

		assert.Nil(t, err)
	})

	t.Run("Publish fail", func(t *testing.T) {
		kafkaPublisher := mock.NewMockPublisher(t)
		publisher := NewSuggestedOfferEventPublisher(
			config.KafkaConfig{
				NotificationTopic: "notification",
			}, kafkaPublisher,
		)
		kafkaPublisher.EXPECT().Publish(testifyMock.Anything).Return(errors.New("test error"))
		err := publisher.NotifySuggestedOfferCreated(context.Background(), "investorId", entity.SuggestedOfferConfig{
			ValueType: entity.ValueTypeInterestRate,
			Value:     decimal.NewFromFloat(0.12),
		}, entity.SuggestedOffer{})

		assert.Equal(t, "SuggestedOfferEventPublisher NotifySuggestedOfferCreated test error", err.Error())
	})
}
