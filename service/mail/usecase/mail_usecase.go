package usecase

import (
	"bytes"
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
	"th3y3m/e-commerce-microservices/service/mail/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type mailUsecase struct {
	log *logrus.Logger
}

type IMailUsecase interface {
	SendMail(to string, token string) error
	SendOrderDetails(Customer model.User, Order model.Order, OrderDetails []model.OrderDetail) error
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

func (m *mailUsecase) SendOrderDetails(Customer model.User, Order model.Order, OrderDetails []model.OrderDetail) error {
	from, password := viper.GetString("EMAIL"), viper.GetString("PASSWORD")
	smtpHost, smtpPort := viper.GetString("SMTP_HOST"), viper.GetString("SMTP_PORT")

	// Populate the form data
	type OrderDetailWithProduct struct {
		OrderDetail model.OrderDetail
		Product     model.Product
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
		var product model.Product
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
	}{
		Customer:                Customer,
		Order:                   Order,
		OrderDetailsWithProduct: orderDetailsWithProduct,
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
