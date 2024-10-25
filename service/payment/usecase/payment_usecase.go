package usecase

import (
	"context"
	"errors"
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
	GetRevenue(ctx context.Context, day, month, year *int) ([]float64, error)
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

func (pu *paymentUsecase) GetRevenue(ctx context.Context, day, month, year *int) ([]float64, error) {
	pu.log.Infof("Fetching revenue for day: %d, month: %d, year: %d", *day, *month, *year)
	currentTime := time.Now()
	var revenue []float64

	if year != nil {
		for i := currentTime.Year(); i >= *year; i-- {
			model := &model.GetPaymentsRequest{
				FromDate: time.Date(i, 1, 1, 0, 0, 0, 0, time.UTC),
				ToDate:   time.Date(i, 12, 31, 23, 59, 59, 0, time.UTC),
			}
			yearsRevenue, err := pu.paymentRepo.GetList(ctx, model)
			if err != nil {
				pu.log.Errorf("Error fetching revenue for year %d: %v", i, err)
				return nil, err
			}
			yearRevenue := yearsRevenue[0].PaymentAmount
			revenue = append(revenue, yearRevenue)
		}
	} else if month != nil {
		startYear := currentTime.Year()
		startMonth := currentTime.Month()

		for y := startYear; y >= startYear-1; y-- {
			for m := startMonth; m >= time.Month(*month); m-- {
				model := &model.GetPaymentsRequest{
					FromDate: time.Date(y, m, 1, 0, 0, 0, 0, time.UTC),
					ToDate:   time.Date(y, m, daysInMonth(y, m), 23, 59, 59, 0, time.UTC),
				}
				monthsRevenue, err := pu.paymentRepo.GetList(ctx, model)
				if err != nil {
					pu.log.Errorf("Error fetching revenue for month %d-%d: %v", m, y, err)
					return nil, err
				}
				monthRevenue := monthsRevenue[0].PaymentAmount
				revenue = append(revenue, monthRevenue)
			}
			startMonth = 12 // Reset to December for the previous year
		}

	} else if day != nil {
		startYear := currentTime.Year()
		startMonth := currentTime.Month()
		startDay := currentTime.Day()

		for y := startYear; y >= startYear-1; y-- {
			for m := startMonth; m >= time.Month(1); m-- {
				for d := startDay; d >= 1; d-- {
					if y == startYear && m == startMonth && d < *day {
						break
					}

					model := &model.GetPaymentsRequest{
						FromDate: time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
						ToDate:   time.Date(y, m, d, 23, 59, 59, 0, time.UTC),
					}
					daysRevenue, err := pu.paymentRepo.GetList(ctx, model)
					if err != nil {
						pu.log.Errorf("Error fetching revenue for day %d-%d-%d: %v", d, m, y, err)
						return nil, err
					}
					dayRevenue := daysRevenue[0].PaymentAmount
					revenue = append(revenue, dayRevenue)
				}
				startDay = daysInMonth(y, m-1) // Reset to the last day of the previous month
			}
			startMonth = 12 // Reset to December for the previous year
		}
	} else {
		pu.log.Errorf("No valid date parameter provided")
		return nil, errors.New("no valid date parameter provided")
	}

	return revenue, nil
}

func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
