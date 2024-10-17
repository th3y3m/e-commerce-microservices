package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"th3y3m/e-commerce-microservices/pkg/constant"
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
	ProcessOrder(ctx context.Context, userId, cartId, CourierID, VoucherID int64, shipAddress, paymentMethod string, freight float64) (*model.GetOrderResponse, error)
	ProcessPayment(ctx context.Context, order repository.Order, paymentMethod string) (string, error)
	UpdateInventory(ctx context.Context, userId, cartId int64) error
	SendNotification(ctx context.Context, orderID int64) error
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

// func (o *orderUsecase) PlaceOrder(userId, cartId, shipAddress, CourierID, VoucherID, paymentMethod string) (string, error) {
// 	// Step 1: Create the order synchronously
// 	order, err := o.ProcessOrder(userId, cartId, shipAddress, CourierID, VoucherID, paymentMethod)
// 	if err != nil {
// 		return "", err
// 	}

// 	err = o.PublishInventoryUpdateEvent(userId, cartId)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Publish Notification Event
// 	err = o.PublishOrderNotificationEvent(order.OrderID)
// 	if err != nil {
// 		return "", err
// 	}

// 	paymentURL, err := o.ProcessPayment(order)
// 	if err != nil {
// 		return "", err
// 	}

// 	return paymentURL, nil
// }

func (o *orderUsecase) ProcessOrder(ctx context.Context, userId, cartId, CourierID, VoucherID int64, shipAddress, paymentMethod string, freight float64) (*model.GetOrderResponse, error) {

	cartItemReq := model.GetCartItemsRequest{
		CartID: &cartId,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return &model.GetOrderResponse{}, err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var productsList []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&productsList)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	totalAmount := 0.0
	for _, product := range productsList {

		url := constant.PRODUCT_SERVICE + fmt.Sprintf("/%d", product.ProductID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return &model.GetOrderResponse{}, err
		}

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return &model.GetOrderResponse{}, err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			return &model.GetOrderResponse{}, fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
		}

		// Decode the response into cart items
		var p model.GetProductResponse
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			o.log.Errorf("Failed to decode response: %v", err)
			return &model.GetOrderResponse{}, err
		}

		totalAmount += p.Price * float64(product.Quantity)
	}

	url = constant.VOUCHER_SERVICE + fmt.Sprintf("/%d", VoucherID)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Set the context and execute the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var voucher model.GetVoucherResponse
	err = json.NewDecoder(resp.Body).Decode(&voucher)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	checkVoucherUsageRequest := model.CheckVoucherUsageRequest{
		VoucherID: VoucherID,
		Order: model.Order{
			OrderID:               0,
			CustomerID:            userId,
			OrderDate:             time.Now(),
			ShippingAddress:       shipAddress,
			CourierID:             CourierID,
			TotalAmount:           totalAmount,
			OrderStatus:           "Pending",
			FreightPrice:          freight,
			EstimatedDeliveryDate: time.Now(),
			ActualDeliveryDate:    time.Now(),
			VoucherID:             VoucherID,
			IsDeleted:             false,
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		},
	}

	checkVoucherUsageRequestData, err := json.Marshal(checkVoucherUsageRequest)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return &model.GetOrderResponse{}, err
	}

	url = constant.VOUCHER_SERVICE + "/check-usage"
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(checkVoucherUsageRequestData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Set the context and execute the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
	}

	var checkVoucherResponse model.CheckVoucherUsageResponse

	// Decode the response body into the struct
	err = json.NewDecoder(resp.Body).Decode(&checkVoucherResponse)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Log the result and return the valid status
	o.log.Infof("Voucher validity: %v", checkVoucherResponse.Valid)

	if !checkVoucherResponse.Valid {
		return &model.GetOrderResponse{}, fmt.Errorf("Voucher is not valid")
	}

	if voucher.DiscountType == constant.VOUCHER_DISCOUNT_TYPE_PERCENTAGE {
		discountPrice := totalAmount * voucher.DiscountValue / 100
		if discountPrice > voucher.MaxDiscountAmount {
			discountPrice = voucher.MaxDiscountAmount
		}
		totalAmount = totalAmount - discountPrice
	} else if voucher.DiscountType == constant.VOUCHER_DISCOUNT_TYPE_FIXED {
		totalAmount = totalAmount - voucher.DiscountValue
	}

	newOrder := model.CreateOrderRequest{
		CustomerID:            userId,
		CourierID:             CourierID,
		VoucherID:             VoucherID,
		TotalAmount:           totalAmount,
		ShippingAddress:       shipAddress,
		FreightPrice:          freight,
		EstimatedDeliveryDate: time.Now(),
		ActualDeliveryDate:    time.Now(),
		OrderStatus:           "Pending",
	}

	createdOrder, err := o.CreateOrder(ctx, &newOrder)
	if err != nil {
		o.log.Errorf("Failed to create order: %v", err)
		return &model.GetOrderResponse{}, err
	}

	for _, item := range productsList {

		url := constant.PRODUCT_SERVICE + fmt.Sprintf("/%d", item.ProductID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return &model.GetOrderResponse{}, err
		}

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return &model.GetOrderResponse{}, err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			return &model.GetOrderResponse{}, fmt.Errorf("ptoduct service returned non-OK status: %d", resp.StatusCode)
		}

		// Decode the response into cart items
		var product model.GetProductResponse
		err = json.NewDecoder(resp.Body).Decode(&product)
		if err != nil {
			o.log.Errorf("Failed to decode response: %v", err)
			return &model.GetOrderResponse{}, err
		}

		createOrderDetailRequest := model.CreateOrderDetailRequest{
			OrderID:   createdOrder.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		}

		orderDetailsData, err := json.Marshal(createOrderDetailRequest)
		if err != nil {
			o.log.Errorf("Failed to marshal order data: %v", err)
			return &model.GetOrderResponse{}, err
		}

		url = constant.ORDER_DETAILS_SERVICE
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(orderDetailsData))
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return &model.GetOrderResponse{}, err
		}

		// Set the context and execute the request
		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return &model.GetOrderResponse{}, err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
			return &model.GetOrderResponse{}, fmt.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
		}
	}

	return createdOrder, nil
}

func (o *orderUsecase) ProcessPayment(ctx context.Context, order repository.Order, paymentMethod string) (string, error) {
	payment := model.CreatePaymentRequest{
		OrderID:       order.OrderID,
		PaymentAmount: order.TotalAmount,
		PaymentMethod: paymentMethod,
		PaymentStatus: constant.PAYMENT_STATUS_PENDING,
	}

	paymentRequest, err := json.Marshal(payment)
	if err != nil {
		o.log.Errorf("Failed to marshal payment data: %v", err)
		return "", err
	}

	url := constant.PAYMENT_SERVICE
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(paymentRequest))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return "", err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
		return "", fmt.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
	}

	if paymentMethod == constant.PAYMENT_METHOD_MOMO {
		// Create the payment record
		url := constant.MOMO_SERVICE + "?amount=" + fmt.Sprintf("%.2f", order.TotalAmount) + "&orderID=" + fmt.Sprintf("%d", order.OrderID)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return "", err
		}

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return "", err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("momo service returned non-OK status: %d", resp.StatusCode)
			return "", fmt.Errorf("momo service returned non-OK status: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			o.log.Errorf("Failed to read response body: %v", err)
			return "", err
		}
		returnUrl := string(bodyBytes)

		o.log.Infof("Received MoMo URL: %s", returnUrl)

		return returnUrl, nil
	}

	if paymentMethod == constant.PAYMENT_METHOD_VNPAY {

		url := constant.VNPAY_SERVICE + "?amount=" + fmt.Sprintf("%.2f", order.TotalAmount) + "&orderID=" + fmt.Sprintf("%d", order.OrderID)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return "", err
		}

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return "", err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("vnpay service returned non-OK status: %d", resp.StatusCode)
			return "", fmt.Errorf("vnpay service returned non-OK status: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			o.log.Errorf("Failed to read response body: %v", err)
			return "", err
		}
		returnUrl := string(bodyBytes)

		o.log.Infof("Received VnPay URL: %s", returnUrl)

		return returnUrl, nil
	}

	return "", nil
}

func (o *orderUsecase) UpdateInventory(ctx context.Context, userId, cartId int64) error {
	cartItemReq := model.GetCartItemsRequest{
		CartID: &cartId,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var productsList []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&productsList)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return err
	}

	for _, product := range productsList {
		url := constant.PRODUCT_SERVICE + fmt.Sprintf("/%d", product.ProductID)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return err
		}

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			return fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
		}

		// Decode the response into cart items
		var p model.GetProductResponse
		err = json.NewDecoder(resp.Body).Decode(&p)
		if err != nil {
			o.log.Errorf("Failed to decode response: %v", err)
			return err
		}

		p.Quantity -= product.Quantity

		updateProductRequest := model.UpdateProductRequest{
			ProductID:   p.ProductID,
			SellerID:    p.SellerID,
			ProductName: p.ProductName,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			CategoryID:  p.CategoryID,
			ImageURL:    p.ImageURL,
		}

		productData, err := json.Marshal(updateProductRequest)
		if err != nil {
			o.log.Errorf("Failed to marshal order data: %v", err)
			return err
		}

		url = constant.PRODUCT_SERVICE
		req, err = http.NewRequest("PUT", url, bytes.NewBuffer(productData))
		if err != nil {
			o.log.Errorf("Failed to create request: %v", err)
			return err
		}

		// Set the context and execute the request
		client = &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute request: %v", err)
			return err
		}

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			return fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
		}
	}

	return nil
}

func (o *orderUsecase) SendNotification(ctx context.Context, orderID int64) error {
	order, err := o.orderRepo.Get(ctx, orderID)
	if err != nil {
		return err
	}

	getUserRequest := model.GetUserRequest{
		UserID: &order.CustomerID,
	}

	userData, err := json.Marshal(getUserRequest)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url := constant.USER_SERVICE + "/get-user"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(userData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("user service returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("user service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var customer model.GetUserResponse
	err = json.NewDecoder(resp.Body).Decode(&customer)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return err
	}

	orderDetailsRequest := model.GetOrderDetailsRequest{
		OrderID: &orderID,
	}

	orderDetailsData, err := json.Marshal(orderDetailsRequest)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url = constant.ORDER_DETAILS_SERVICE
	req, err = http.NewRequest("GET", url, bytes.NewBuffer(orderDetailsData))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the context and execute the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return err
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var orderDetails []model.GetOrderDetailResponse
	err = json.NewDecoder(resp.Body).Decode(&orderDetails)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return err
	}

	user := model.User{
		UserID:       customer.UserID,
		Email:        customer.Email,
		PasswordHash: customer.PasswordHash,
		FullName:     customer.FullName,
		PhoneNumber:  customer.PhoneNumber,
		Address:      customer.Address,
		Role:         customer.Role,
		ImageURL:     customer.ImageURL,
		CreatedAt:    util.ParseTime(customer.CreatedAt),
		UpdatedAt:    util.ParseTime(customer.UpdatedAt),
		Token:        customer.Token,
		TokenExpires: util.ParseTime(customer.TokenExpires),
		IsDeleted:    customer.IsDeleted,
	}

	orderModel := model.Order{
		OrderID:               order.OrderID,
		CustomerID:            order.CustomerID,
		OrderDate:             order.OrderDate,
		ShippingAddress:       order.ShippingAddress,
		CourierID:             order.CourierID,
		TotalAmount:           order.TotalAmount,
		OrderStatus:           order.OrderStatus,
		FreightPrice:          order.FreightPrice,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate,
		ActualDeliveryDate:    order.ActualDeliveryDate,
		VoucherID:             order.VoucherID,
		IsDeleted:             order.IsDeleted,
		CreatedAt:             order.CreatedAt,
		UpdatedAt:             order.UpdatedAt,
	}

	var orderDetailsModel []model.OrderDetail
	for _, detail := range orderDetails {
		orderDetailsModel = append(orderDetailsModel, model.OrderDetail{
			OrderID:   detail.OrderID,
			ProductID: detail.ProductID,
			Quantity:  detail.Quantity,
			UnitPrice: detail.UnitPrice,
		})
	}

	request := model.SendOrderDetailsRequest{
		Customer:     user,
		Order:        orderModel,
		OrderDetails: orderDetailsModel,
	}

	mailModel, err := json.Marshal(request)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url = constant.MAIL_SERVICE + "/send-order-details"
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(mailModel))
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return err
	}

	// Set the context and execute the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("mail service returned non-OK status: %d", resp.StatusCode)
		return fmt.Errorf("mail service returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
