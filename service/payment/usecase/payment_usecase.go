package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/payment/model"
	"th3y3m/e-commerce-microservices/service/payment/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type paymentUsecase struct {
	log         *logrus.Logger
	paymentRepo repository.IPaymentRepository
}

type IPaymentUsecase interface {
	GetPayment(ctx context.Context, req *model.GetPaymentRequest) (*model.GetPaymentResponse, error)
	GetAllPayments(ctx context.Context) ([]*model.GetPaymentResponse, error)
	CreatePayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.GetPaymentResponse, error)
	UpdatePayment(ctx context.Context, rep *model.UpdatePaymentRequest) (*model.GetPaymentResponse, error)
	GetPaymentList(ctx context.Context, req *model.GetPaymentsRequest) (*util.PaginatedList[model.GetPaymentResponse], error)
}

func NewPaymentUsecase(paymentRepo repository.IPaymentRepository, log *logrus.Logger) IPaymentUsecase {
	return &paymentUsecase{
		paymentRepo: paymentRepo,
		log:         log,
	}
}

func (pu *paymentUsecase) GetPayment(ctx context.Context, req *model.GetPaymentRequest) (*model.GetPaymentResponse, error) {
	pu.log.Infof("Fetching payment with ID: %d", req.PaymentID)
	payment, err := pu.paymentRepo.Get(ctx, req.PaymentID)
	if err != nil {
		pu.log.Errorf("Error fetching payment: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched payment: %+v", payment)
	return &model.GetPaymentResponse{
		PaymentID:        payment.PaymentID,
		OrderID:          payment.OrderID,
		PaymentAmount:    payment.PaymentAmount,
		PaymentDate:      payment.PaymentDate.Format(tsCreateTimeLayout),
		PaymentMethod:    payment.PaymentMethod,
		PaymentStatus:    payment.PaymentStatus,
		PaymentSignature: payment.PaymentSignature,
	}, nil
}

func (pu *paymentUsecase) GetAllPayments(ctx context.Context) ([]*model.GetPaymentResponse, error) {
	pu.log.Info("Fetching all payments")
	payments, err := pu.paymentRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all payments: %v", err)
		return nil, err
	}

	var paymentResponses []*model.GetPaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, &model.GetPaymentResponse{
			PaymentID:        payment.PaymentID,
			OrderID:          payment.OrderID,
			PaymentAmount:    payment.PaymentAmount,
			PaymentDate:      payment.PaymentDate.Format(tsCreateTimeLayout),
			PaymentMethod:    payment.PaymentMethod,
			PaymentStatus:    payment.PaymentStatus,
			PaymentSignature: payment.PaymentSignature,
		})
	}

	pu.log.Infof("Fetched %d payments", len(paymentResponses))
	return paymentResponses, nil
}

func (pu *paymentUsecase) CreatePayment(ctx context.Context, payment *model.CreatePaymentRequest) (*model.GetPaymentResponse, error) {
	pu.log.Infof("Creating payment: %+v", payment)
	paymentData := repository.Payment{
		OrderID:          payment.OrderID,
		PaymentAmount:    payment.PaymentAmount,
		PaymentDate:      time.Now(),
		PaymentMethod:    payment.PaymentMethod,
		PaymentStatus:    payment.PaymentStatus,
		PaymentSignature: payment.PaymentSignature,
	}

	createdPayment, err := pu.paymentRepo.Create(ctx, &paymentData)
	if err != nil {
		pu.log.Errorf("Error creating payment: %v", err)
		return nil, err
	}

	pu.log.Infof("Created payment: %+v", createdPayment)
	return &model.GetPaymentResponse{
		PaymentID:        createdPayment.PaymentID,
		OrderID:          createdPayment.OrderID,
		PaymentAmount:    createdPayment.PaymentAmount,
		PaymentDate:      createdPayment.PaymentDate.Format(tsCreateTimeLayout),
		PaymentMethod:    createdPayment.PaymentMethod,
		PaymentStatus:    createdPayment.PaymentStatus,
		PaymentSignature: createdPayment.PaymentSignature,
	}, nil
}

func (pu *paymentUsecase) UpdatePayment(ctx context.Context, rep *model.UpdatePaymentRequest) (*model.GetPaymentResponse, error) {
	pu.log.Infof("Updating payment with ID: %d", rep.PaymentID)
	payment, err := pu.paymentRepo.Get(ctx, rep.PaymentID)
	if err != nil {
		pu.log.Errorf("Error fetching payment for update: %v", err)
		return nil, err
	}

	payment.PaymentAmount = rep.PaymentAmount
	payment.PaymentMethod = rep.PaymentMethod
	payment.PaymentStatus = rep.PaymentStatus
	payment.PaymentSignature = rep.PaymentSignature

	updatedPayment, err := pu.paymentRepo.Update(ctx, payment)
	if err != nil {
		pu.log.Errorf("Error updating payment: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated payment: %+v", updatedPayment)
	return &model.GetPaymentResponse{
		PaymentID:        updatedPayment.PaymentID,
		OrderID:          updatedPayment.OrderID,
		PaymentAmount:    updatedPayment.PaymentAmount,
		PaymentDate:      updatedPayment.PaymentDate.Format(tsCreateTimeLayout),
		PaymentMethod:    updatedPayment.PaymentMethod,
		PaymentStatus:    updatedPayment.PaymentStatus,
		PaymentSignature: updatedPayment.PaymentSignature,
	}, nil
}

func (pu *paymentUsecase) GetPaymentList(ctx context.Context, req *model.GetPaymentsRequest) (*util.PaginatedList[model.GetPaymentResponse], error) {
	pu.log.Infof("Fetching payment list with request: %+v", req)
	payments, err := pu.paymentRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching payment list: %v", err)
		return nil, err
	}

	var paymentResponses []model.GetPaymentResponse
	for _, payment := range payments {
		paymentResponses = append(paymentResponses, model.GetPaymentResponse{
			PaymentID:        payment.PaymentID,
			OrderID:          payment.OrderID,
			PaymentAmount:    payment.PaymentAmount,
			PaymentDate:      payment.PaymentDate.Format(tsCreateTimeLayout),
			PaymentMethod:    payment.PaymentMethod,
			PaymentStatus:    payment.PaymentStatus,
			PaymentSignature: payment.PaymentSignature,
		})
	}

	list := &util.PaginatedList[model.GetPaymentResponse]{
		Items:      paymentResponses,
		TotalCount: len(paymentResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d payments", len(paymentResponses))
	return list, nil
}
