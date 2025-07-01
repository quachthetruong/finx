package kafka

import (
	"context"
	"financing-offer/internal/apperrors"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/shopspring/decimal"
	"gitlab.com/enCapital/models"
	"gitlab.com/enCapital/models/dnse"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/timestamppb"

	"financing-offer/internal/config"
	"financing-offer/internal/core/entity"
	"financing-offer/internal/core/loanofferinterest/repository"
	"financing-offer/internal/event"
	"financing-offer/pkg/pb"
)

var _ repository.LoanPackageOfferInterestEventRepository = (*LoanOfferInterestEventPublisher)(nil)

type LoanOfferInterestEventPublisher struct {
	config         config.KafkaConfig
	publisher      event.Publisher
	temporalClient client.Client
}

func (l *LoanOfferInterestEventPublisher) NotifyLoanPackageOfferReady(_ context.Context, data entity.LoanPackageOfferReadyNotify) error {
	payload := dnse.FinancingOfferLoanPackageReady{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		LoanRate:      data.LoanRate.Mul(decimal.NewFromInt(100)).String() + "%",
		LoanType:      data.LoanType.StringNotify(),
		InterestRate:  data.InterestRate.Mul(decimal.NewFromInt(100)).String() + "%/năm",
		LoanPackageId: data.LoanPackageId,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_LOAN_PACKAGE_READY)
	if err != nil {
		return fmt.Errorf("LoanOfferInterestEventPublisher NotifyLoanPackageOfferReady %w", err)
	}
	if err := l.publisher.Publish(
		kafka.Message{
			Topic: l.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanOfferInterestEventPublisher NotifyLoanPackageOfferReady %w", err)
	}
	return nil
}

func (l *LoanOfferInterestEventPublisher) NotifyDerivativeLoanPackageOfferReady(_ context.Context, data entity.DerivativeLoanPackageOfferReadyNotify) error {
	payload := dnse.FinancingOfferDerivativeLoanPackageReady{
		InvestorId:    data.InvestorId,
		RequestName:   data.RequestName,
		AccountNo:     data.AccountNo,
		AccountNoDesc: data.AccountNoDesc,
		Symbol:        data.Symbol,
		LoanPackageId: data.LoanPackageId,
		AssetType:     data.AssetType,
		CreatedAt:     timestamppb.New(data.CreatedAt),
	}
	message, err := pb.MarshalEncapMessage(&payload, models.EncapMessage_FINANCING_OFFER_DERIVATIVE_LOAN_PACKAGE_READY)
	if err != nil {
		return fmt.Errorf("LoanOfferInterestEventPublisher NotifyDerivativeLoanPackageOfferReady %w", err)
	}
	if err := l.publisher.Publish(
		kafka.Message{
			Topic: l.config.NotificationTopic,
			Value: message,
			Key:   []byte(data.InvestorId),
		},
	); err != nil {
		return fmt.Errorf("LoanOfferInterestEventPublisher NotifyDerivativeLoanPackageOfferReady %w", err)
	}
	return nil
}

func (l *LoanOfferInterestEventPublisher) CreateMarginLoanPackage(ctx context.Context, state entity.AssignmentState) error {
	workflowId := fmt.Sprintf("auto-create-loan-package-v3-%d", state.Submission.LoanPackageOfferInterestId)
	workflowOptions := client.StartWorkflowOptions{
		ID:                                       workflowId,
		TaskQueue:                                config.SavingsTaskQueueName,
		WorkflowExecutionErrorWhenAlreadyStarted: false, // do not return error if workflow is already started
	}
	workflowRun, err := l.temporalClient.ExecuteWorkflow(
		ctx, workflowOptions, "CreateAndAssignLoanPackageWorkflow", state,
	)
	if err != nil {
		return fmt.Errorf("SavingsTemporalRepository SettleSavingsAccount: %w", err)
	}
	if err := workflowRun.Get(ctx, nil); err != nil {
		return apperrors.ErrorCreateAndAssignLoanPackageWorkflow
	}
	return nil
}

func NewLoanOfferInterestEventPublisher(config config.KafkaConfig, publisher event.Publisher, temporalClient client.Client) *LoanOfferInterestEventPublisher {
	return &LoanOfferInterestEventPublisher{
		config:         config,
		publisher:      publisher,
		temporalClient: temporalClient,
	}
}
