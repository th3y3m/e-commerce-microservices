package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/order_detail/model"
	"th3y3m/e-commerce-microservices/service/order_detail/repository"

	"github.com/sirupsen/logrus"
)

type orderDetailUsecase struct {
	log             *logrus.Logger
	orderDetailRepo repository.IOrderDetailRepository
}

type IOrderDetailUsecase interface {
	GetOrderDetail(ctx context.Context, req *model.GetOrderDetailRequest) (*model.GetOrderDetailResponse, error)
	CreateOrderDetail(ctx context.Context, req *model.CreateOrderDetailRequest) (*model.GetOrderDetailResponse, error)
	UpdateOrderDetail(ctx context.Context, rep *model.UpdateOrderDetailRequest) (*model.GetOrderDetailResponse, error)
	DeleteOrderDetail(ctx context.Context, req *model.DeleteOrderDetailRequest) error
	GetOrderDetailList(ctx context.Context, req *model.GetOrderDetailsRequest) ([]*model.GetOrderDetailResponse, error)
}

func NewOrderDetailUsecase(orderDetailRepo repository.IOrderDetailRepository, log *logrus.Logger) IOrderDetailUsecase {
	return &orderDetailUsecase{
		orderDetailRepo: orderDetailRepo,
		log:             log,
	}
}

func (pu *orderDetailUsecase) GetOrderDetail(ctx context.Context, req *model.GetOrderDetailRequest) (*model.GetOrderDetailResponse, error) {
	pu.log.Infof("Fetching orderDetail with OrderID: %d and ProductID: %d", req.OrderID, req.ProductID)
	orderDetail, err := pu.orderDetailRepo.Get(ctx, req.OrderID, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching orderDetail: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched orderDetail: %+v", orderDetail)
	return &model.GetOrderDetailResponse{
		OrderID:   orderDetail.OrderID,
		ProductID: orderDetail.ProductID,
		Quantity:  orderDetail.Quantity,
		UnitPrice: orderDetail.UnitPrice,
	}, nil
}

func (pu *orderDetailUsecase) CreateOrderDetail(ctx context.Context, orderDetail *model.CreateOrderDetailRequest) (*model.GetOrderDetailResponse, error) {
	pu.log.Infof("Creating orderDetail: %+v", orderDetail)
	orderDetailData := repository.OrderDetail{
		OrderID:   orderDetail.OrderID,
		ProductID: orderDetail.ProductID,
		Quantity:  orderDetail.Quantity,
		UnitPrice: orderDetail.UnitPrice,
	}

	createdOrderDetail, err := pu.orderDetailRepo.Create(ctx, &orderDetailData)
	if err != nil {
		pu.log.Errorf("Error creating orderDetail: %v", err)
		return nil, err
	}

	pu.log.Infof("Created orderDetail: %+v", createdOrderDetail)
	return &model.GetOrderDetailResponse{
		OrderID:   createdOrderDetail.OrderID,
		ProductID: createdOrderDetail.ProductID,
		Quantity:  createdOrderDetail.Quantity,
		UnitPrice: createdOrderDetail.UnitPrice,
	}, nil
}

func (pu *orderDetailUsecase) DeleteOrderDetail(ctx context.Context, req *model.DeleteOrderDetailRequest) error {
	pu.log.Infof("Deleting orderDetail with OrderID: %d and ProductID: %d", req.OrderID, req.ProductID)
	orderDetail, err := pu.orderDetailRepo.Get(ctx, req.OrderID, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching orderDetail for deletion: %v", err)
		return err
	}

	err = pu.orderDetailRepo.Delete(ctx, orderDetail.OrderID, orderDetail.ProductID)
	if err != nil {
		pu.log.Errorf("Error deleting orderDetail: %v", err)
		return err
	}

	pu.log.Infof("Deleted orderDetail with OrderID: %d and ProductID: %d", req.OrderID, req.ProductID)
	return nil
}

func (pu *orderDetailUsecase) UpdateOrderDetail(ctx context.Context, rep *model.UpdateOrderDetailRequest) (*model.GetOrderDetailResponse, error) {
	pu.log.Infof("Updating orderDetail: %+v", rep)
	orderDetail, err := pu.orderDetailRepo.Get(ctx, rep.OrderID, rep.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching orderDetail for update: %v", err)
		return nil, err
	}

	orderDetail.Quantity = rep.Quantity
	orderDetail.UnitPrice = rep.UnitPrice

	updatedOrderDetail, err := pu.orderDetailRepo.Update(ctx, orderDetail)
	if err != nil {
		pu.log.Errorf("Error updating orderDetail: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated orderDetail: %+v", updatedOrderDetail)
	return &model.GetOrderDetailResponse{
		OrderID:   updatedOrderDetail.OrderID,
		ProductID: updatedOrderDetail.ProductID,
		Quantity:  updatedOrderDetail.Quantity,
		UnitPrice: updatedOrderDetail.UnitPrice,
	}, nil
}

func (pu *orderDetailUsecase) GetOrderDetailList(ctx context.Context, req *model.GetOrderDetailsRequest) ([]*model.GetOrderDetailResponse, error) {
	pu.log.Infof("Fetching orderDetails with OrderID: %d", req.OrderID)
	orderDetails, err := pu.orderDetailRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching orderDetails: %v", err)
		return nil, err
	}

	var orderDetailResponses []*model.GetOrderDetailResponse
	for _, orderDetail := range orderDetails {
		orderDetailResponses = append(orderDetailResponses, &model.GetOrderDetailResponse{
			OrderID:   orderDetail.OrderID,
			ProductID: orderDetail.ProductID,
			Quantity:  orderDetail.Quantity,
			UnitPrice: orderDetail.UnitPrice,
		})
	}

	pu.log.Infof("Fetched orderDetails: %+v", orderDetailResponses)
	return orderDetailResponses, nil
}
