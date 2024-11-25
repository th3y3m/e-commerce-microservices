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
	"th3y3m/e-commerce-microservices/service/order/rabbitmq"
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
	ProcessPayment(ctx context.Context, order *model.GetOrderResponse, paymentMethod string) (string, error)
	PlaceOrder(ctx context.Context, userId, cartId, CourierID, VoucherID int64, paymentMethod, shipAddress string, freight float64) (string, error)
	CancelOrder(ctx context.Context, orderID int64) error
}

func NewOrderUsecase(orderRepo repository.IOrderRepository, log *logrus.Logger) IOrderUsecase {
	return &orderUsecase{
		orderRepo: orderRepo,
		log:       log,
	}
}

func (pu *orderUsecase) CancelOrder(ctx context.Context, orderID int64) error {
	pu.log.Infof("Cancelling order with ID: %d", orderID)
	order, err := pu.orderRepo.Get(ctx, orderID)
	if err != nil {
		pu.log.Errorf("Error fetching order for cancellation: %v", err)
		return err
	}

	order.OrderStatus = constant.ORDER_STATUS_CANCELLED

	_, err = pu.orderRepo.Update(ctx, order)
	if err != nil {
		pu.log.Errorf("Error updating order for cancellation: %v", err)
		return err
	}

	pu.log.Infof("Cancelled order with ID: %d", orderID)
	return nil
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

func (o *orderUsecase) PlaceOrder(ctx context.Context, userId, cartId, CourierID, VoucherID int64, paymentMethod, shipAddress string, freight float64) (string, error) {
	order, err := o.ProcessOrder(ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight)
	if err != nil {
		return "", err
	}
	paymentURL, err := o.ProcessPayment(ctx, order, paymentMethod)
	if err != nil {
		return "", err
	}

	// Use a channel to capture errors from goroutines
	errChan := make(chan error, 2)

	// Publish Inventory Update Event in a goroutine
	go func() {
		errChan <- rabbitmq.PublishInventoryUpdateEvent(userId, cartId)
	}()

	// Publish Order Notification Event in a goroutine
	go func() {
		errChan <- rabbitmq.PublishOrderNotificationEvent(order.OrderID, paymentURL)
	}()

	// Wait for both goroutines to finish and check for errors
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			return "", err
		}
	}

	return paymentURL, nil
}

// func AutomaticFailedOrder(ctx context.Context, orderID int64, duration time.Duration) {
// 	select {
// 	case <-time.After(duration):
// 		isPaid, err := CheckPaymentStatus(orderID)
// 		if err != nil {
// 			// Log the error if necessary
// 			return
// 		}

// 		if !isPaid {
// 			// If the payment failed, cancel the order and return products to inventory
// 			err = CancelOrder(ctx, orderID)
// 			if err != nil {
// 				// Handle error (log it or notify)
// 				return
// 			}
// 		}
// 	}
// }

func (o *orderUsecase) ProcessOrder(ctx context.Context, userId, cartId, CourierID, VoucherID int64, shipAddress, paymentMethod string, freight float64) (*model.GetOrderResponse, error) {
	client := &http.Client{}

	// Fetch cart items
	cartItemReq := model.GetCartItemsRequest{
		CartID: &cartId,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		o.log.Errorf("Failed to marshal cart item request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequestWithContext(ctx, "GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		o.log.Errorf("Failed to create cart item request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute cart item request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("Cart item service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	var productsList []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&productsList)
	if err != nil {
		o.log.Errorf("Failed to decode cart item response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Fetch product details and calculate total amount
	totalAmount := 0.0
	productDetails := make(map[int64]model.GetProductResponse)
	for _, product := range productsList {
		if _, exists := productDetails[product.ProductID]; !exists {
			url := constant.PRODUCT_SERVICE + fmt.Sprintf("/%d", product.ProductID)
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				o.log.Errorf("Failed to create product request: %v", err)
				return &model.GetOrderResponse{}, err
			}

			resp, err := client.Do(req)
			if err != nil {
				o.log.Errorf("Failed to execute product request: %v", err)
				return &model.GetOrderResponse{}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				o.log.Errorf("Product service returned non-OK status: %d", resp.StatusCode)
				return &model.GetOrderResponse{}, fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			}

			var p model.GetProductResponse
			err = json.NewDecoder(resp.Body).Decode(&p)
			if err != nil {
				o.log.Errorf("Failed to decode product response: %v", err)
				return &model.GetOrderResponse{}, err
			}

			url = constant.PRODUCT_SERVICE + "/discount-price" + fmt.Sprintf("/%d", product.ProductID)

			req, err = http.NewRequestWithContext(ctx, "GET", url, nil)

			if err != nil {
				o.log.Errorf("Failed to create product discount request: %v", err)
				return &model.GetOrderResponse{}, err
			}

			resp, err = client.Do(req)
			if err != nil {
				o.log.Errorf("Failed to execute product discount request: %v", err)
				return &model.GetOrderResponse{}, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				o.log.Errorf("Product service returned non-OK status: %d", resp.StatusCode)
				return &model.GetOrderResponse{}, fmt.Errorf("product service returned non-OK status: %d", resp.StatusCode)
			}

			var discountPrice float64
			err = json.NewDecoder(resp.Body).Decode(&discountPrice)
			if err != nil {
				o.log.Errorf("Failed to decode product discount response: %v", err)
				return &model.GetOrderResponse{}, err
			}

			p.Price = discountPrice

			productDetails[product.ProductID] = p
		}

		totalAmount += productDetails[product.ProductID].Price * float64(product.Quantity)
	}

	// Fetch voucher details
	url = constant.VOUCHER_SERVICE + fmt.Sprintf("/%d", VoucherID)
	req, err = http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		o.log.Errorf("Failed to create voucher request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute voucher request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("Voucher service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
	}

	var voucher model.GetVoucherResponse
	err = json.NewDecoder(resp.Body).Decode(&voucher)
	if err != nil {
		o.log.Errorf("Failed to decode voucher response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	// Check voucher usage
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
		o.log.Errorf("Failed to marshal check voucher usage request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	url = constant.VOUCHER_SERVICE + "/check-usage"
	req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(checkVoucherUsageRequestData))
	if err != nil {
		o.log.Errorf("Failed to create check voucher usage request: %v", err)
		return &model.GetOrderResponse{}, err
	}

	resp, err = client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute check voucher usage request: %v", err)
		return &model.GetOrderResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("Voucher service returned non-OK status: %d", resp.StatusCode)
		return &model.GetOrderResponse{}, fmt.Errorf("voucher service returned non-OK status: %d", resp.StatusCode)
	}

	var checkVoucherResponse model.CheckVoucherUsageResponse
	err = json.NewDecoder(resp.Body).Decode(&checkVoucherResponse)
	if err != nil {
		o.log.Errorf("Failed to decode check voucher usage response: %v", err)
		return &model.GetOrderResponse{}, err
	}

	o.log.Infof("Voucher validity: %v", checkVoucherResponse.Valid)
	if !checkVoucherResponse.Valid {
		return &model.GetOrderResponse{}, fmt.Errorf("Voucher is not valid")
	}

	// Apply voucher discount
	if voucher.DiscountType == constant.VOUCHER_DISCOUNT_TYPE_PERCENTAGE {
		discountPrice := totalAmount * voucher.DiscountValue / 100
		if discountPrice > voucher.MaxDiscountAmount {
			discountPrice = voucher.MaxDiscountAmount
		}
		totalAmount -= discountPrice
	} else if voucher.DiscountType == constant.VOUCHER_DISCOUNT_TYPE_FIXED {
		totalAmount -= voucher.DiscountValue
	}

	// Create order
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

	// Create order details
	for _, item := range productsList {
		product := productDetails[item.ProductID]

		createOrderDetailRequest := model.CreateOrderDetailRequest{
			OrderID:   createdOrder.OrderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		}

		orderDetailsData, err := json.Marshal(createOrderDetailRequest)
		if err != nil {
			o.log.Errorf("Failed to marshal order detail request: %v", err)
			return &model.GetOrderResponse{}, err
		}

		url = constant.ORDER_DETAILS_SERVICE
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(orderDetailsData))
		if err != nil {
			o.log.Errorf("Failed to create order detail request: %v", err)
			return &model.GetOrderResponse{}, err
		}

		resp, err = client.Do(req)
		if err != nil {
			o.log.Errorf("Failed to execute order detail request: %v", err)
			return &model.GetOrderResponse{}, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			o.log.Errorf("Order detail service returned non-OK status: %d", resp.StatusCode)
			return &model.GetOrderResponse{}, fmt.Errorf("order detail service returned non-OK status: %d", resp.StatusCode)
		}
	}

	return createdOrder, nil
}

func (o *orderUsecase) ProcessPayment(ctx context.Context, order *model.GetOrderResponse, paymentMethod string) (string, error) {
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
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(paymentRequest))
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
		o.log.Errorf("Payment service returned non-OK status: %d", resp.StatusCode)
		return "", fmt.Errorf("Payment service returned non-OK status: %d", resp.StatusCode)
	}

	if paymentMethod == constant.PAYMENT_METHOD_MOMO {
		return o.processMomoPayment(ctx, order)
	}

	if paymentMethod == constant.PAYMENT_METHOD_VNPAY {
		return o.processVnPayPayment(ctx, order)
	}

	return "", nil
}

func (o *orderUsecase) processMomoPayment(ctx context.Context, order *model.GetOrderResponse) (string, error) {
	url := constant.MOMO_SERVICE + "?amount=" + fmt.Sprintf("%.2f", order.TotalAmount) + "&orderID=" + fmt.Sprintf("%d", order.OrderID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("Momo service returned non-OK status: %d", resp.StatusCode)
		return "", fmt.Errorf("Momo service returned non-OK status: %d", resp.StatusCode)
	}

	var paymentUrl model.MoMoResponse
	err = json.NewDecoder(resp.Body).Decode(&paymentUrl)
	if err != nil {
		o.log.Errorf("Failed to decode cart item response: %v", err)
		return paymentUrl.PaymentURL, err
	}

	o.log.Infof("Received MoMo URL: %s", paymentUrl.PaymentURL)

	return paymentUrl.PaymentURL, nil
}

func (o *orderUsecase) processVnPayPayment(ctx context.Context, order *model.GetOrderResponse) (string, error) {
	url := constant.VNPAY_SERVICE + "?amount=" + fmt.Sprintf("%.2f", order.TotalAmount) + "&orderID=" + fmt.Sprintf("%d", order.OrderID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		o.log.Errorf("Failed to create request: %v", err)
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		o.log.Errorf("Failed to execute request: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		o.log.Errorf("VnPay service returned non-OK status: %d", resp.StatusCode)
		return "", fmt.Errorf("VnPay service returned non-OK status: %d", resp.StatusCode)
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
