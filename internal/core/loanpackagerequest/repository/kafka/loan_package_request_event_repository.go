package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
	"gitlab.com/enCapital/models"
	"gitlab.com/enCapital/models/dnse"
	"google.golang.org/protobuf/types/known/timestamppb"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanpackagerequest/repository"
	"financing-offer/internal/event"
	"financing-offer/pkg/pb"
)

var _ repository.LoanPackageRequestEventRepository = (*LoanPackageRequestEventPublisher)(nil)

type LoanPackageRequestEventPublisher struct {
	config    config.KafkaConfig
	publisher event.Publisher
}

func NewLoanPackageRequestEventPublisher(config config.KafkaConfig, publisher event.Publisher) *LoanPackageRequestEventPublisher {
	return &LoanPackageRequestEventPublisher{
		config:    config,
		publisher: publisher,
	}
}

func (p *LoanPackageRequestEventPublisher) NotifyOfflineConfirmation(_ context.Context, data entity.RequestOfflineConfirmation) error {
	payload := dnse.FinancingOfferOfflineConfirmation{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_OFFLINE_CONFIRMATION)
	if err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyOfflineConfirmation %w", err)
	}
	if err := p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyOfflineConfirmation %w", err)
	}
	return nil
}

func (p *LoanPackageRequestEventPublisher) NotifyOnlineConfirmation(_ context.Context, data entity.RequestOnlineConfirmationNotify) error {
	payload := dnse.FinancingOfferOnlineConfirmation{
		InvestorId:      data.InvestorId,
		RequestName:     data.RequestName,
		AccountNo:       data.AccountNo,
		AccountNoDesc:   data.AccountNoDesc,
		OfferId:         data.OfferId,
		OfferInterestId: data.OfferInterestId,
		Symbol:          data.Symbol,
		CreatedAt:       timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_ONLINE_CONFIRMATION)
	fmt.Printf("message %+v\n", message)

	if err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyOnlineConfirmation %w", err)
	}
	if err := p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyOnlineConfirmation %w", err)
	}
	return nil
}

func (p *LoanPackageRequestEventPublisher) NotifyRequestDeclined(_ context.Context, data entity.LoanPackageRequestDeclinedNotify) error {
	payload := dnse.FinancingOfferRequestDeclined{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_REQUEST_DECLINED)
	if err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyRequestDeclined %w", err)
	}
	if err := p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyRequestConfirmed %w", err)
	}
	return nil
}

func (p *LoanPackageRequestEventPublisher) NotifyDerivativeRequestDeclined(_ context.Context, data entity.LoanPackageDerivativeRequestDeclinedNotify) error {
	payload := dnse.FinancingOfferDerivativeRequestDeclined{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		AssetType:     data.AssetType,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_DERIVATIVE_REQUEST_DECLINED)
	if err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyDerivativeRequestDeclined %w", err)
	}
	if err := p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyDerivativeRequestDeclined %w", err)
	}
	return nil
}

func (p *LoanPackageRequestEventPublisher) NotifyDerivativeOfflineConfirmation(_ context.Context, data entity.DerivativeRequestOfflineConfirmation) error {
	payload := dnse.FinancingOfferDerivativeOfflineConfirmation{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		AssetType:     data.AssetType,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_DERIVATIVE_OFFLINE_CONFIRMATION)
	if err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyDerivativeOfflineConfirmation %w", err)
	}
	if err := p.publisher.Publish(
		kafka.Message{
			Topic: p.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanPackageRequestEventPublisher NotifyDerivativeOfflineConfirmation %w", err)
	}
	return nil
}
