package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/order/model"
	"th3y3m/e-commerce-microservices/service/order/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type orderUsecase struct {
	log       *logrus.Logger
	orderRepo repository.IOrderRepository
}

type IOrderUsecase interface {
	GetOrder(ctx context.Context, req *model.GetOrderRequest) (*model.GetOrderResponse, error)
	GetAllOrders(ctx context.Context) ([]*model.GetOrderResponse, error)
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.GetOrderResponse, error)
	UpdateOrder(ctx context.Context, rep *model.UpdateOrderRequest) (*model.GetOrderResponse, error)
	DeleteOrder(ctx context.Context, req *model.DeleteOrderRequest) error
	GetOrderList(ctx context.Context, req *model.GetOrdersRequest) (*util.PaginatedList[model.GetOrderResponse], error)
}

func NewOrderUsecase(orderRepo repository.IOrderRepository, log *logrus.Logger) IOrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
		log:       log,
	}
}

func (pu *orderUsecase) GetOrder(ctx context.Context, req *model.GetOrderRequest) (*model.GetOrderResponse, error) {
	pu.log.Infof("Fetching order with ID: %d", req.OrderID)
	order, err := pu.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		pu.log.Errorf("Error fetching order: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched order: %+v", order)
	return &model.GetOrderResponse{
		OrderID:               order.OrderID,
		CustomerID:            order.CustomerID,
		OrderDate:             order.OrderDate.Format(tsCreateTimeLayout),
		ShippingAddress:       order.ShippingAddress,
		CourierID:             order.CourierID,
		TotalAmount:           order.TotalAmount,
		OrderStatus:           order.OrderStatus,
		FreightPrice:          order.FreightPrice,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate.Format(tsCreateTimeLayout),
		ActualDeliveryDate:    order.ActualDeliveryDate.Format(tsCreateTimeLayout),
		VoucherID:             order.VoucherID,
		IsDeleted:             order.IsDeleted,
		CreatedAt:             order.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:             order.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *orderUsecase) GetAllOrders(ctx context.Context) ([]*model.GetOrderResponse, error) {
	pu.log.Infof("Fetching all orders")
	orders, err := pu.orderRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching orders: %v", err)
		return nil, err
	}

	var orderResponses []*model.GetOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, &model.GetOrderResponse{
			OrderID:               order.OrderID,
			CustomerID:            order.CustomerID,
			OrderDate:             order.OrderDate.Format(tsCreateTimeLayout),
			ShippingAddress:       order.ShippingAddress,
			CourierID:             order.CourierID,
			TotalAmount:           order.TotalAmount,
			OrderStatus:           order.OrderStatus,
			FreightPrice:          order.FreightPrice,
			EstimatedDeliveryDate: order.EstimatedDeliveryDate.Format(tsCreateTimeLayout),
			ActualDeliveryDate:    order.ActualDeliveryDate.Format(tsCreateTimeLayout),
			VoucherID:             order.VoucherID,
			IsDeleted:             order.IsDeleted,
			CreatedAt:             order.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:             order.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d orders", len(orderResponses))
	return orderResponses, nil
}

func (pu *orderUsecase) CreateOrder(ctx context.Context, order *model.CreateOrderRequest) (*model.GetOrderResponse, error) {
	pu.log.Infof("Creating order: %+v", order)
	orderData := repository.Order{
		CustomerID:            order.CustomerID,
		OrderDate:             time.Now(),
		ShippingAddress:       order.ShippingAddress,
		CourierID:             order.CourierID,
		TotalAmount:           order.TotalAmount,
		OrderStatus:           order.OrderStatus,
		FreightPrice:          order.FreightPrice,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate,
		ActualDeliveryDate:    order.ActualDeliveryDate,
		VoucherID:             order.VoucherID,
	}

	createdOrder, err := pu.orderRepo.Create(ctx, &orderData)
	if err != nil {
		pu.log.Errorf("Error creating order: %v", err)
		return nil, err
	}

	pu.log.Infof("Created order: %+v", createdOrder)
	return &model.GetOrderResponse{
		OrderID:               createdOrder.OrderID,
		CustomerID:            createdOrder.CustomerID,
		OrderDate:             createdOrder.OrderDate.Format(tsCreateTimeLayout),
		ShippingAddress:       createdOrder.ShippingAddress,
		CourierID:             createdOrder.CourierID,
		TotalAmount:           createdOrder.TotalAmount,
		OrderStatus:           createdOrder.OrderStatus,
		FreightPrice:          createdOrder.FreightPrice,
		EstimatedDeliveryDate: createdOrder.EstimatedDeliveryDate.Format(tsCreateTimeLayout),
		ActualDeliveryDate:    createdOrder.ActualDeliveryDate.Format(tsCreateTimeLayout),
		VoucherID:             createdOrder.VoucherID,
		IsDeleted:             createdOrder.IsDeleted,
		CreatedAt:             createdOrder.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:             createdOrder.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *orderUsecase) DeleteOrder(ctx context.Context, req *model.DeleteOrderRequest) error {
	pu.log.Infof("Deleting order with ID: %d", req.OrderID)
	order, err := pu.orderRepo.Get(ctx, req.OrderID)
	if err != nil {
		pu.log.Errorf("Error fetching order for deletion: %v", err)
		return err
	}

	order.IsDeleted = true

	_, err = pu.orderRepo.Update(ctx, order)
	if err != nil {
		pu.log.Errorf("Error updating order for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted order with ID: %d", req.OrderID)
	return nil
}

func (pu *orderUsecase) UpdateOrder(ctx context.Context, rep *model.UpdateOrderRequest) (*model.GetOrderResponse, error) {
	pu.log.Infof("Updating order with ID: %d", rep.OrderID)
	order, err := pu.orderRepo.Get(ctx, rep.OrderID)
	if err != nil {
		pu.log.Errorf("Error fetching order for update: %v", err)
		return nil, err
	}

	order.CustomerID = rep.CustomerID
	order.OrderDate = rep.OrderDate
	order.ShippingAddress = rep.ShippingAddress
	order.CourierID = rep.CourierID
	order.TotalAmount = rep.TotalAmount
	order.OrderStatus = rep.OrderStatus
	order.FreightPrice = rep.FreightPrice
	order.VoucherID = rep.VoucherID
	order.ActualDeliveryDate = rep.ActualDeliveryDate
	order.EstimatedDeliveryDate = rep.EstimatedDeliveryDate
	order.UpdatedAt = time.Now()
	order.IsDeleted = rep.IsDeleted

	updatedOrder, err := pu.orderRepo.Update(ctx, order)
	if err != nil {
		pu.log.Errorf("Error updating order: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated order: %+v", updatedOrder)
	return &model.GetOrderResponse{
		OrderID:               updatedOrder.OrderID,
		CustomerID:            updatedOrder.CustomerID,
		OrderDate:             updatedOrder.OrderDate.Format(tsCreateTimeLayout),
		ShippingAddress:       updatedOrder.ShippingAddress,
		CourierID:             updatedOrder.CourierID,
		TotalAmount:           updatedOrder.TotalAmount,
		OrderStatus:           updatedOrder.OrderStatus,
		FreightPrice:          updatedOrder.FreightPrice,
		EstimatedDeliveryDate: updatedOrder.EstimatedDeliveryDate.Format(tsCreateTimeLayout),
		ActualDeliveryDate:    updatedOrder.ActualDeliveryDate.Format(tsCreateTimeLayout),
		VoucherID:             updatedOrder.VoucherID,
		IsDeleted:             updatedOrder.IsDeleted,
		CreatedAt:             updatedOrder.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:             updatedOrder.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *orderUsecase) GetOrderList(ctx context.Context, req *model.GetOrdersRequest) (*util.PaginatedList[model.GetOrderResponse], error) {
	pu.log.Infof("Fetching order list with request: %+v", req)
	orders, err := pu.orderRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching order list: %v", err)
		return nil, err
	}

	var orderResponses []model.GetOrderResponse
	for _, order := range orders {
		orderResponses = append(orderResponses, model.GetOrderResponse{
			OrderID:               order.OrderID,
			CustomerID:            order.CustomerID,
			OrderDate:             order.OrderDate.Format(tsCreateTimeLayout),
			ShippingAddress:       order.ShippingAddress,
			CourierID:             order.CourierID,
			TotalAmount:           order.TotalAmount,
			OrderStatus:           order.OrderStatus,
			FreightPrice:          order.FreightPrice,
			EstimatedDeliveryDate: order.EstimatedDeliveryDate.Format(tsCreateTimeLayout),
			ActualDeliveryDate:    order.ActualDeliveryDate.Format(tsCreateTimeLayout),
			VoucherID:             order.VoucherID,
			IsDeleted:             order.IsDeleted,
			CreatedAt:             order.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:             order.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	list := &util.PaginatedList[model.GetOrderResponse]{
		Items:      orderResponses,
		TotalCount: len(orderResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d orders", len(orderResponses))
	return list, nil
}
