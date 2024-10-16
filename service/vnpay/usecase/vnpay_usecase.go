package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/vnpay/model"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewVnpayUsecase(log *logrus.Logger) IVnpayUsecase {

	return &VnpayUsecase{
		url:        viper.GetString("VNPAY_URL"),
		returnUrl:  viper.GetString("VNPAY_RETURN_URL"),
		tmnCode:    viper.GetString("VNPAY_TMNCODE"),
		hashSecret: viper.GetString("VNPAY_HASH_SECRET"),
		log:        log,
	}
}

type IVnpayUsecase interface {
	CreateVNPayUrl(amount float64, orderinfor string) (string, error)
	ValidateVNPayResponse(queryString url.Values) (*model.PaymentResponse, error)
}

type VnpayUsecase struct {
	url        string
	returnUrl  string
	tmnCode    string
	hashSecret string
	log        *logrus.Logger
}

func (s *VnpayUsecase) CreateVNPayUrl(amount float64, orderID string) (string, error) {
	hostName, err := os.Hostname()
	if err != nil {
		return "", err
	}

	ipAddrs, err := net.LookupIP(hostName)
	if err != nil || len(ipAddrs) == 0 {
		return "", err
	}
	clientIPAddress := ipAddrs[0].String()

	pay := util.NewPayLib()
	vnpAmount := amount
	pay.AddRequestData("vnp_Version", "2.1.0")
	pay.AddRequestData("vnp_Command", "pay")
	pay.AddRequestData("vnp_TmnCode", s.tmnCode)
	pay.AddRequestData("vnp_Amount", fmt.Sprintf("%.0f", vnpAmount))
	pay.AddRequestData("vnp_BankCode", "")
	pay.AddRequestData("vnp_CreateDate", time.Now().Format("20060102150405"))
	pay.AddRequestData("vnp_CurrCode", "VND")
	pay.AddRequestData("vnp_IpAddr", clientIPAddress)
	pay.AddRequestData("vnp_Locale", "vn")
	pay.AddRequestData("vnp_OrderInfo", "Customer")
	pay.AddRequestData("vnp_OrderType", "other")
	pay.AddRequestData("vnp_ReturnUrl", s.returnUrl)
	pay.AddRequestData("vnp_TxnRef", orderID)

	TransactionUrl := pay.CreateRequestUrl(s.url, s.hashSecret)
	return TransactionUrl, nil
}

func (s *VnpayUsecase) ValidateVNPayResponse(queryString url.Values) (*model.PaymentResponse, error) {

	vnpSecureHash := queryString.Get("vnp_SecureHash")
	vnpAmount := queryString.Get("vnp_Amount")
	queryString.Del("vnp_SecureHash")
	queryString.Del("vnp_SecureHashType")

	rawData := make([]string, 0, len(queryString))
	for key, val := range queryString {
		rawData = append(rawData, key+"="+strings.Join(val, ""))
	}
	sort.Strings(rawData)
	rawQueryString := strings.Join(rawData, "&")

	if !s.ValidateSignature(rawQueryString, vnpSecureHash, s.hashSecret) {
		return &model.PaymentResponse{IsSuccessful: false, RedirectUrl: "LINK_INVALID"}, nil
	}

	res, err := http.Get(fmt.Sprintf("%s/%s", constant.ORDER_SERVICE, queryString.Get("vnp_TxnRef")))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d from order service", res.StatusCode)
	}

	var order model.Order
	if err := json.NewDecoder(res.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("error decoding order response: %v", err)
	}

	if order.OrderStatus == constant.ORDER_STATUS_COMPLETED {
		return &model.PaymentResponse{
			IsSuccessful: false,
			RedirectUrl:  "LINK_INVALID",
		}, nil
	}

	vnpResponseCode := queryString.Get("vnp_ResponseCode")
	if vnpResponseCode == "00" && queryString.Get("vnp_TransactionStatus") == "00" {
		order.OrderStatus = constant.ORDER_STATUS_COMPLETED
		updateModel := model.UpdateOrderRequest{
			OrderID:               order.OrderID,
			CustomerID:            order.CustomerID,
			OrderDate:             order.OrderDate,
			OrderStatus:           order.OrderStatus,
			ActualDeliveryDate:    order.ActualDeliveryDate,
			EstimatedDeliveryDate: order.EstimatedDeliveryDate,
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

		paymentAmount, err := strconv.ParseFloat(vnpAmount, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid payment amount: %v", err)
		}

		paymentCreateModel := &model.CreatePaymentRequest{
			OrderID:          order.OrderID,
			PaymentAmount:    paymentAmount,
			PaymentStatus:    constant.ORDER_STATUS_COMPLETED,
			PaymentSignature: queryString.Get("vnp_BankTranNo"),
			PaymentMethod:    constant.PAYMENT_METHOD_VNPAY,
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

		// cart, err := s.shoppingCartService.GetUserShoppingCart(order.CustomerID)
		// if err != nil {
		// 	return nil, err
		// }

		// if err := s.shoppingCartService.UpdateShoppingCartStatus(cart.CartID, false); err != nil {
		// 	return nil, err
		// }

		return &model.PaymentResponse{
			IsSuccessful: true,
			RedirectUrl:  constant.PAYMENT_RESPONSE_CONFIRM_URL + "?orderId=" + strconv.FormatInt(order.OrderID, 10),
		}, nil
	}

	order.OrderStatus = constant.ORDER_STATUS_FAILED
	updateModel := model.UpdateOrderRequest{
		OrderID:               order.OrderID,
		CustomerID:            order.CustomerID,
		OrderDate:             order.OrderDate,
		OrderStatus:           order.OrderStatus,
		ActualDeliveryDate:    order.ActualDeliveryDate,
		EstimatedDeliveryDate: order.EstimatedDeliveryDate,
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
		RedirectUrl:  constant.PAYMENT_RESPONSE_REJECT_URL + "?orderId=" + queryString.Get("vnp_TxnRef"),
	}, nil
}

func (s *VnpayUsecase) ValidateSignature(rspraw, inputHash, secretKey string) bool {
	return util.HmacSHA512(secretKey, rspraw) == inputHash
}
