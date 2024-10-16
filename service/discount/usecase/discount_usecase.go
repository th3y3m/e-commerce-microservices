package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/discount/model"
	"th3y3m/e-commerce-microservices/service/discount/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type discountUsecase struct {
	log          *logrus.Logger
	discountRepo repository.IDiscountRepository
}

type IDiscountUsecase interface {
	GetDiscount(ctx context.Context, req *model.GetDiscountRequest) (*model.GetDiscountResponse, error)
	GetAllDiscounts(ctx context.Context) ([]*model.GetDiscountResponse, error)
	CreateDiscount(ctx context.Context, req *model.CreateDiscountRequest) (*model.GetDiscountResponse, error)
	UpdateDiscount(ctx context.Context, rep *model.UpdateDiscountRequest) (*model.GetDiscountResponse, error)
	DeleteDiscount(ctx context.Context, req *model.DeleteDiscountRequest) error
}

func NewDiscountUsecase(discountRepo repository.IDiscountRepository, log *logrus.Logger) IDiscountUsecase {
	return &discountUsecase{
		discountRepo: discountRepo,
		log:          log,
	}
}

func (pu *discountUsecase) GetDiscount(ctx context.Context, req *model.GetDiscountRequest) (*model.GetDiscountResponse, error) {
	pu.log.Infof("Fetching discount with ID: %d", req.DiscountID)
	discount, err := pu.discountRepo.Get(ctx, req.DiscountID)
	if err != nil {
		pu.log.Errorf("Error fetching discount: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched discount: %+v", discount)
	return &model.GetDiscountResponse{
		DiscountID:    discount.DiscountID,
		DiscountType:  discount.DiscountType,
		DiscountValue: discount.DiscountValue,
		StartDate:     discount.StartDate.Format(tsCreateTimeLayout),
		EndDate:       discount.EndDate.Format(tsCreateTimeLayout),
		IsDeleted:     discount.IsDeleted,
		CreatedAt:     discount.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     discount.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *discountUsecase) GetAllDiscounts(ctx context.Context) ([]*model.GetDiscountResponse, error) {
	pu.log.Info("Fetching all discounts")
	discounts, err := pu.discountRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching discounts: %v", err)
		return nil, err
	}

	var discountResponses []*model.GetDiscountResponse
	for _, discount := range discounts {
		discountResponses = append(discountResponses, &model.GetDiscountResponse{
			DiscountID:    discount.DiscountID,
			DiscountType:  discount.DiscountType,
			DiscountValue: discount.DiscountValue,
			StartDate:     discount.StartDate.Format(tsCreateTimeLayout),
			EndDate:       discount.EndDate.Format(tsCreateTimeLayout),
			IsDeleted:     discount.IsDeleted,
			CreatedAt:     discount.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:     discount.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d discounts", len(discounts))
	return discountResponses, nil
}

func (pu *discountUsecase) CreateDiscount(ctx context.Context, discount *model.CreateDiscountRequest) (*model.GetDiscountResponse, error) {
	pu.log.Infof("Creating discount: %+v", discount)
	createdDiscount, err := pu.discountRepo.Create(ctx, &repository.Discount{
		DiscountType:  discount.DiscountType,
		DiscountValue: discount.DiscountValue,
		StartDate:     discount.StartDate,
		EndDate:       discount.EndDate,
	})
	if err != nil {
		pu.log.Errorf("Error creating discount: %v", err)
		return nil, err
	}

	pu.log.Infof("Created discount: %+v", createdDiscount)
	return &model.GetDiscountResponse{
		DiscountID:    createdDiscount.DiscountID,
		DiscountType:  createdDiscount.DiscountType,
		DiscountValue: createdDiscount.DiscountValue,
		StartDate:     createdDiscount.StartDate.Format(tsCreateTimeLayout),
		EndDate:       createdDiscount.EndDate.Format(tsCreateTimeLayout),
		IsDeleted:     createdDiscount.IsDeleted,
		CreatedAt:     createdDiscount.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     createdDiscount.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *discountUsecase) DeleteDiscount(ctx context.Context, req *model.DeleteDiscountRequest) error {
	pu.log.Infof("Deleting discount with ID: %d", req.DiscountID)
	discount, err := pu.discountRepo.Get(ctx, req.DiscountID)
	if err != nil {
		pu.log.Errorf("Error fetching discount for deletion: %v", err)
		return err
	}

	discount.IsDeleted = true

	_, err = pu.discountRepo.Update(ctx, discount)
	if err != nil {
		pu.log.Errorf("Error updating discount for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted discount with ID: %d", req.DiscountID)
	return nil
}

func (pu *discountUsecase) UpdateDiscount(ctx context.Context, rep *model.UpdateDiscountRequest) (*model.GetDiscountResponse, error) {
	pu.log.Infof("Updating discount with ID: %d", rep.DiscountID)
	discount, err := pu.discountRepo.Get(ctx, rep.DiscountID)
	if err != nil {
		pu.log.Errorf("Error fetching discount for update: %v", err)
		return nil, err
	}

	discount.DiscountType = rep.DiscountType
	discount.DiscountValue = rep.DiscountValue
	discount.StartDate = rep.StartDate
	discount.EndDate = rep.EndDate
	discount.UpdatedAt = time.Now()

	updatedDiscount, err := pu.discountRepo.Update(ctx, discount)
	if err != nil {
		pu.log.Errorf("Error updating discount: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated discount: %+v", updatedDiscount)
	return &model.GetDiscountResponse{
		DiscountID:    updatedDiscount.DiscountID,
		DiscountType:  updatedDiscount.DiscountType,
		DiscountValue: updatedDiscount.DiscountValue,
		StartDate:     updatedDiscount.StartDate.Format(tsCreateTimeLayout),
		EndDate:       updatedDiscount.EndDate.Format(tsCreateTimeLayout),
		IsDeleted:     updatedDiscount.IsDeleted,
		CreatedAt:     updatedDiscount.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     updatedDiscount.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}
