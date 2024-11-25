package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/mail/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type mailUsecase struct {
	log *logrus.Logger
}

type IMailUsecase interface {
	SendMail(to string, token string) error
	SendOrderDetails(Customer model.User, Order model.Order, OrderDetails []model.OrderDetail, urlPayment string) error
	SendNotification(ctx context.Context, orderID int64, url string) error
}

func NewMailUsecase(log *logrus.Logger) IMailUsecase {
	return &mailUsecase{
		log: log,
	}
}

// SendMail sends the email to the user
func (m *mailUsecase) SendMail(to string, token string) error {

	// Retrieve environment variables
	from, password := viper.GetString("EMAIL"), viper.GetString("PASSWORD")
	smtpHost, smtpPort := viper.GetString("SMTP_HOST"), viper.GetString("SMTP_PORT")

	// Construct the file path for the HTML template
	htmlPath := filepath.Join("templates", "Confirmation.html") // Adjust path as needed

	// Read the HTML template
	htmlTemplate, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Printf("Failed to read HTML template: %v", err)
		return err
	}

	// Replace the {{TOKEN}} placeholder with the actual token
	htmlContent := strings.Replace(string(htmlTemplate), "{{TOKEN}}", token, 1)

	// Set up the email headers and body
	subject := "Subject: Verify your email\n"
	msg := []byte(subject + "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" + htmlContent)

	// Set up SMTP authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", to, err)
		return err
	}

	return nil
}

func (m *mailUsecase) SendOrderDetails(Customer model.User, Order model.Order, OrderDetails []model.OrderDetail, urlPayment string) error {
	from, password := viper.GetString("EMAIL"), viper.GetString("PASSWORD")
	smtpHost, smtpPort := viper.GetString("SMTP_HOST"), viper.GetString("SMTP_PORT")

	// Populate the form data
	type OrderDetailWithProduct struct {
		OrderDetail model.OrderDetail
		Product     model.GetProductResponse
	}

	var orderDetailsWithProduct []OrderDetailWithProduct
	for _, od := range OrderDetails {
		// Fetch the product details from the product service
		url := constant.PRODUCT_SERVICE + "/" + strconv.FormatInt(od.ProductID, 10)

		res, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to get product details for product ID %d: %v", od.ProductID, err)
			return err
		}
		defer res.Body.Close()

		// Check if the request was successful
		if res.StatusCode != http.StatusOK {
			log.Printf("Error: Received status code %d from product service", res.StatusCode)
			return fmt.Errorf("product service returned status code %d", res.StatusCode)
		}

		// Unmarshal the response into the Product struct
		var product model.GetProductResponse
		err = json.NewDecoder(res.Body).Decode(&product)
		if err != nil {
			log.Printf("Failed to decode product details for product ID %d: %v", od.ProductID, err)
			return err
		}

		// Append the OrderDetail and corresponding Product to the slice
		orderDetailsWithProduct = append(orderDetailsWithProduct, OrderDetailWithProduct{
			OrderDetail: od,
			Product:     product,
		})
	}

	// Form data for the email template
	form := struct {
		Customer                model.User
		Order                   model.Order
		OrderDetailsWithProduct []OrderDetailWithProduct
		UrlPayment              string
	}{
		Customer:                Customer,
		Order:                   Order,
		OrderDetailsWithProduct: orderDetailsWithProduct,
		UrlPayment:              urlPayment,
	}

	// Load and parse the email HTML template
	htmlPath := filepath.Join("templates", "Notifycation.html")
	htmlTemplate, err := os.ReadFile(htmlPath)
	if err != nil {
		log.Printf("Failed to read HTML template: %v", err)
		return err
	}

	tmpl, err := template.New("email").Funcs(template.FuncMap{
		"multiply": func(a, b interface{}) float64 {
			var af, bf float64
			switch v := a.(type) {
			case int:
				af = float64(v)
			case float64:
				af = v
			}
			switch v := b.(type) {
			case int:
				bf = float64(v)
			case float64:
				bf = v
			}
			return af * bf
		},
		"toFloat64": func(i interface{}) float64 {
			switch v := i.(type) {
			case int:
				return float64(v)
			case float64:
				return v
			default:
				return 0
			}
		},
		"formatCurrency": func(amount float64) string {
			return fmt.Sprintf("%.0f VND", amount)
		},
		"formatWithSpaces": func(amount float64) string {
			s := fmt.Sprintf("%.0f", amount)
			n := len(s)
			if n <= 3 {
				return s
			}
			var result strings.Builder
			for i, c := range s {
				if i > 0 && (n-i)%3 == 0 {
					result.WriteRune(' ')
				}
				result.WriteRune(c)
			}
			return result.String() + " VND"
		},
	}).Parse(string(htmlTemplate))
	if err != nil {
		log.Printf("Failed to parse HTML template: %v", err)
		return err
	}

	var htmlContent bytes.Buffer
	err = tmpl.Execute(&htmlContent, form)
	if err != nil {
		log.Printf("Failed to execute HTML template: %v", err)
		return err
	}

	// Prepare the email content
	subject := "Subject: Order Details\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	msg := []byte(subject + mime + htmlContent.String())

	// SMTP authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send the email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{form.Customer.Email}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", form.Customer.Email, err)
		return err
	}

	return nil
}

func (o *mailUsecase) SendNotification(ctx context.Context, orderID int64, urlPayment string) error {
	url := constant.ORDER_SERVICE + "/" + strconv.FormatInt(orderID, 10)

	// Fetch the order details from the order service
	res, err := http.Get(url)
	if err != nil {
		o.log.Errorf("Failed to get order details: %v", err)
		return err
	}

	defer res.Body.Close()

	// Check if the request was successful
	if res.StatusCode != http.StatusOK {
		o.log.Errorf("order service returned non-OK status: %d", res.StatusCode)
		return fmt.Errorf("order service returned non-OK status: %d", res.StatusCode)
	}

	// Decode the response into the GetOrderResponse struct
	var order model.GetOrderResponse
	err = json.NewDecoder(res.Body).Decode(&order)
	if err != nil {
		o.log.Errorf("Failed to decode response: %v", err)
		return err
	}

	// Fetch the user details from the user service

	getUserRequest := model.GetUserRequest{
		UserID: &order.CustomerID,
	}

	userData, err := json.Marshal(getUserRequest)
	if err != nil {
		o.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url = constant.USER_SERVICE + "/get-user"
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
		OrderID: &order.OrderID,
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
		OrderDate:             util.ParseTime(order.OrderDate),
		ShippingAddress:       order.ShippingAddress,
		CourierID:             order.CourierID,
		TotalAmount:           order.TotalAmount,
		OrderStatus:           order.OrderStatus,
		FreightPrice:          order.FreightPrice,
		EstimatedDeliveryDate: util.ParseTime(order.EstimatedDeliveryDate),
		ActualDeliveryDate:    util.ParseTime(order.ActualDeliveryDate),
		VoucherID:             order.VoucherID,
		IsDeleted:             order.IsDeleted,
		CreatedAt:             util.ParseTime(order.CreatedAt),
		UpdatedAt:             util.ParseTime(order.UpdatedAt),
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

	err = o.SendOrderDetails(user, orderModel, orderDetailsModel, urlPayment)
	if err != nil {
		o.log.Errorf("Failed to send email: %v", err)
		return err
	}

	return nil
}
