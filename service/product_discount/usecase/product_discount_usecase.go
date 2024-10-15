package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/product_discount/model"
	"th3y3m/e-commerce-microservices/service/product_discount/repository"

	"github.com/sirupsen/logrus"
)

type productDiscountUsecase struct {
	log                 *logrus.Logger
	productDiscountRepo repository.IProductDiscountRepository
}

type IProductDiscountUsecase interface {
	CreateProductDiscount(ctx context.Context, req *model.CreateProductDiscountRequest) (*model.GetProductDiscountResponse, error)
	DeleteProductDiscount(ctx context.Context, req *model.DeleteProductDiscountRequest) error
	GetProductDiscountList(ctx context.Context, req *model.GetProductDiscountsRequest) ([]*model.GetProductDiscountResponse, error)
}

func NewProductDiscountUsecase(productDiscountRepo repository.IProductDiscountRepository, log *logrus.Logger) IProductDiscountUsecase {
	return &productDiscountUsecase{
		productDiscountRepo: productDiscountRepo,
		log:                 log,
	}
}

func (pu *productDiscountUsecase) CreateProductDiscount(ctx context.Context, productDiscount *model.CreateProductDiscountRequest) (*model.GetProductDiscountResponse, error) {
	pu.log.Infof("Creating productDiscount: %+v", productDiscount)
	productDiscountData := repository.ProductDiscount{
		ProductID:  productDiscount.ProductID,
		DiscountID: productDiscount.DiscountID,
	}

	createdProductDiscount, err := pu.productDiscountRepo.Create(ctx, &productDiscountData)
	if err != nil {
		pu.log.Errorf("Error creating productDiscount: %v", err)
		return nil, err
	}

	pu.log.Infof("Created productDiscount: %+v", createdProductDiscount)
	return &model.GetProductDiscountResponse{
		ProductID:  createdProductDiscount.ProductID,
		DiscountID: createdProductDiscount.DiscountID,
	}, nil
}

func (pu *productDiscountUsecase) DeleteProductDiscount(ctx context.Context, req *model.DeleteProductDiscountRequest) error {
	pu.log.Infof("Deleting productDiscount with ProductID: %d and DiscountID: %d", req.ProductID, req.DiscountID)

	request := &model.GetProductDiscountsRequest{
		ProductID:  &req.ProductID,
		DiscountID: &req.DiscountID,
	}
	productDiscounts, err := pu.productDiscountRepo.Get(ctx, request)
	if err != nil {
		pu.log.Errorf("Error fetching productDiscount for deletion: %v", err)
		return err
	}

	err = pu.productDiscountRepo.Delete(ctx, productDiscounts[0].ProductID, productDiscounts[0].DiscountID)
	if err != nil {
		pu.log.Errorf("Error deleting productDiscount: %v", err)
		return err
	}

	pu.log.Infof("Deleted productDiscount with ProductID: %d and DisctountID: %d", req.ProductID, req.DiscountID)
	return nil
}

func (pu *productDiscountUsecase) GetProductDiscountList(ctx context.Context, req *model.GetProductDiscountsRequest) ([]*model.GetProductDiscountResponse, error) {
	pu.log.Infof("Fetching productDiscount list with request: %+v", req)

	productDiscounts, err := pu.productDiscountRepo.Get(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching productDiscount list: %v", err)
		return nil, err
	}

	productDiscountResponses := make([]*model.GetProductDiscountResponse, 0)
	for _, productDiscount := range productDiscounts {
		productDiscountResponses = append(productDiscountResponses, &model.GetProductDiscountResponse{
			ProductID:  productDiscount.ProductID,
			DiscountID: productDiscount.DiscountID,
		})
	}

	pu.log.Infof("Fetched productDiscount list: %+v", productDiscountResponses)
	return productDiscountResponses, nil
}
