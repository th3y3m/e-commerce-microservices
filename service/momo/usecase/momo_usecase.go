package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/momo/model"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// IMoMoUsecase is the interface that defines the MoMo usecase methods.
type IMoMoUsecase interface {
	CreateMoMoUrl(amount float64, orderId string) (string, error)
	ValidateMoMoResponse(queryString url.Values) (*model.PaymentResponse, error)
}

func NewMoMoUsecase(log *logrus.Logger) IMoMoUsecase {

	return &MoMoService{
		endpoint:    viper.GetString("MOMO_ENDPOINT"),
		secretKey:   viper.GetString("MOMO_SECRET_KEY"),
		accessKey:   viper.GetString("MOMO_ACCESS_KEY"),
		returnUrl:   viper.GetString("MOMO_RETURN_URL"),
		notifyUrl:   viper.GetString("MOMO_NOTIFY_URL"),
		partnerCode: viper.GetString("MOMO_PARTNER_CODE"),
		requestType: viper.GetString("MOMO_REQUEST_TYPE"),
		extraData:   viper.GetString("MOMO_EXTRA_DATA"),
		log:         log,
	}
}

type MoMoService struct {
	endpoint    string
	secretKey   string
	accessKey   string
	returnUrl   string
	notifyUrl   string
	partnerCode string
	requestType string
	extraData   string
	log         *logrus.Logger
}

// CreatePaymentUrl generates a payment URL for the given amount and order details.
func (s *MoMoService) CreateMoMoUrl(amount float64, orderId string) (string, error) {
	requestId := uuid.New().String()
	orderInfo := "Customer"
	formattedAmount := int64(amount) // Convert to VND
	orderID := fmt.Sprintf("%s-%s", orderId, uuid.New().String())

	// Create raw signature string
	rawHash := fmt.Sprintf("accessKey=%s&amount=%d&extraData=%s&ipnUrl=%s&orderId=%s&orderInfo=%s&partnerCode=%s&redirectUrl=%s&requestId=%s&requestType=%s",
		s.accessKey, formattedAmount, s.extraData, s.notifyUrl, orderID, orderInfo, s.partnerCode, s.returnUrl, requestId, s.requestType)
	signature := util.HmacSHA256(s.secretKey, rawHash)

	// Build request payload
	paymentRequest := map[string]interface{}{
		"partnerCode": s.partnerCode,
		"partnerName": "MoMo",
		"storeId":     "MoMoStore",
		"requestId":   requestId,
		"amount":      strconv.FormatInt(formattedAmount, 10),
		"orderId":     orderID,
		"orderInfo":   orderInfo,
		"redirectUrl": s.returnUrl,
		"ipnUrl":      s.notifyUrl,
		"extraData":   s.extraData,
		"requestType": s.requestType,
		"signature":   signature,
		"lang":        "en",
	}

	// Send POST request to MoMo API
	response, err := util.SendHttpRequest(s.endpoint, paymentRequest)
	if err != nil {
		return "", err
	}

	// Parse response and extract payment URL
	var jsonResponse map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jsonResponse); err != nil {
		return "", err
	}

	if payUrl, ok := jsonResponse["payUrl"].(string); ok {
		return payUrl, nil
	}

	if message, ok := jsonResponse["message"].(string); ok {
		return "", fmt.Errorf("error creating payment URL: %s", message)
	}

	return "", errors.New("unexpected response from MoMo API")
}

func (s *MoMoService) ValidateMoMoResponse(queryString url.Values) (*model.PaymentResponse, error) {
	combinedOrderId := queryString.Get("orderId")
	resultCode := queryString.Get("resultCode")
	amount := queryString.Get("amount")
	signature := queryString.Get("signature")

	orderIdParts := strings.Split(combinedOrderId, "-")
	if len(orderIdParts) < 2 {
		return nil, fmt.Errorf("invalid orderId format")
	}
	orderId := orderIdParts[0]

	// Fetch the order details
	res, err := http.Get(fmt.Sprintf("%s/%s", constant.ORDER_SERVICE, orderId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d from order service", res.StatusCode)
	}

	var order model.GetOrderResponse
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("error decoding order response: %v", err)
	}

	if order.OrderID == 0 || order.OrderStatus == constant.ORDER_STATUS_COMPLETED {
		return &model.PaymentResponse{
			IsSuccessful: false,
			RedirectUrl:  constant.PAYMENT_RESPONSE_REJECT_URL,
		}, nil
	}

	if resultCode == "0" {
		order.OrderStatus = constant.ORDER_STATUS_COMPLETED
		updateModel := model.UpdateOrderRequest{
			OrderID:               order.OrderID,
			CustomerID:            order.CustomerID,
			OrderDate:             util.ParseTime(order.OrderDate),
			OrderStatus:           order.OrderStatus,
			ActualDeliveryDate:    util.ParseTime(order.ActualDeliveryDate),
			EstimatedDeliveryDate: util.ParseTime(order.EstimatedDeliveryDate),
			ShippingAddress:       order.ShippingAddress,
			CourierID:             order.CourierID,
			FreightPrice:          order.FreightPrice,
			TotalAmount:           order.TotalAmount,
			VoucherID:             order.VoucherID,
			IsDeleted:             order.IsDeleted,
		}

		// Update the order status
		url := constant.ORDER_SERVICE
		orderData, err := json.Marshal(updateModel)
		if err != nil {
			s.log.Errorf("Failed to marshal order data: %v", err)
			return nil, err
		}

		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(orderData))
		if err != nil {
			s.log.Errorf("Failed to create request: %v", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			s.log.Errorf("Failed to update order in order service: %v", err)
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			s.log.Errorf("Error updating order: received status %v", res.StatusCode)
			return nil, errors.New("error updating order")
		}

		// Create the payment record
		paymentAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid payment amount: %v", err)
		}
		paymentCreateModel := &model.CreatePaymentRequest{
			OrderID:          order.OrderID,
			PaymentAmount:    paymentAmount,
			PaymentStatus:    constant.PAYMENT_STATUS_COMPLETED,
			PaymentSignature: signature,
			PaymentMethod:    constant.PAYMENT_METHOD_MOMO,
		}

		url = constant.PAYMENT_SERVICE
		paymentData, err := json.Marshal(paymentCreateModel)
		if err != nil {
			s.log.Errorf("Failed to marshal payment data: %v", err)
			return nil, err
		}

		req, err = http.NewRequest("POST", url, bytes.NewBuffer(paymentData))
		if err != nil {
			s.log.Errorf("Failed to create request: %v", err)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

		client = &http.Client{}
		res, err = client.Do(req)
		if err != nil {
			s.log.Errorf("Failed to create payment in payment service: %v", err)
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			s.log.Errorf("Error creating payment: received status %v", res.StatusCode)
			return nil, errors.New("error creating payment")
		}

		return &model.PaymentResponse{
			IsSuccessful: true,
			RedirectUrl:  constant.PAYMENT_RESPONSE_CONFIRM_URL + "?orderId=" + orderId,
		}, nil
	}

	// Handle payment failure
	order.OrderStatus = constant.ORDER_STATUS_FAILED
	updateModel := model.UpdateOrderRequest{
		OrderID:               order.OrderID,
		CustomerID:            order.CustomerID,
		OrderDate:             util.ParseTime(order.OrderDate),
		OrderStatus:           order.OrderStatus,
		ActualDeliveryDate:    util.ParseTime(order.ActualDeliveryDate),
		EstimatedDeliveryDate: util.ParseTime(order.EstimatedDeliveryDate),
		ShippingAddress:       order.ShippingAddress,
		CourierID:             order.CourierID,
		FreightPrice:          order.FreightPrice,
		TotalAmount:           order.TotalAmount,
		VoucherID:             order.VoucherID,
		IsDeleted:             order.IsDeleted,
	}

	url := constant.ORDER_SERVICE
	orderData, err := json.Marshal(updateModel)
	if err != nil {
		s.log.Errorf("Failed to marshal order data: %v", err)
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(orderData))
	if err != nil {
		s.log.Errorf("Failed to create request: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		s.log.Errorf("Failed to update order in order service: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		s.log.Errorf("Error updating order: received status %v", res.StatusCode)
		return nil, errors.New("error updating order")
	}

	return &model.PaymentResponse{
		IsSuccessful: false,
		RedirectUrl:  constant.PAYMENT_RESPONSE_REJECT_URL + "?orderId=" + orderId,
	}, nil
}
