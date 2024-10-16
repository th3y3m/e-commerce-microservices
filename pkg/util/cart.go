package util

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// GetCartFromCookie decodes the cookie value and returns a map of CartItems, with error handling
func GetCartFromCookie(cookieValue string) (map[int64]Item, error) {
	cart := make(map[int64]Item)
	decodedBytes, err := base64.StdEncoding.DecodeString(cookieValue)
	if err != nil {
		return nil, fmt.Errorf("failed to decode cookie value: %w", err)
	}
	decodedString := string(decodedBytes)
	itemsList := strings.Split(decodedString, "|")

	for _, strItem := range itemsList {
		if strItem != "" {
			arrItemDetail := strings.Split(strItem, ",")
			if len(arrItemDetail) < 2 {
				return nil, errors.New("invalid cart item format")
			}
			productID, err := strconv.ParseInt(strings.TrimSpace(arrItemDetail[0]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse product ID: %w", err)
			}
			quantity, err := strconv.Atoi(strings.TrimSpace(arrItemDetail[1]))
			if err != nil {
				return nil, fmt.Errorf("failed to parse quantity: %w", err)
			}

			item := Item{
				ProductID: productID,
				Quantity:  quantity,
			}
			cart[productID] = item
		}
	}

	return cart, nil
}

// GetCookieByName retrieves a cookie by name from the request and returns an error if not found
func GetCookieByName(r *http.Request, cookieName string) (*http.Cookie, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, fmt.Errorf("cookie %s not found: %w", cookieName, err)
	}
	return cookie, nil
}

// SaveCartToCookie saves the encoded cart string into a cookie
func SaveCartToCookie(w http.ResponseWriter, cartString string, userId int64) error {
	if cartString == "" {
		return errors.New("cart string is empty")
	}
	if userId == 0 {
		return errors.New("user ID is empty")
	}
	cookieName := "Cart_" + strconv.FormatInt(userId, 10)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    cartString,
		MaxAge:   7200, // 120 minutes
		HttpOnly: true,
		Path:     "/",
	})
	return nil
}

// DeleteCartToCookie deletes the cart cookie
func DeleteCartToCookie(w http.ResponseWriter, userId int64) error {
	if userId == 0 {
		return errors.New("user ID is empty")
	}
	cookieName := "Cart_" + strconv.FormatInt(userId, 10)
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Expires:  time.Now().Add(-24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	return nil
}

// ConvertCartToString converts the list of CartItems to an encoded string
func ConvertCartToString(itemsList []Item) (string, error) {
	if len(itemsList) == 0 {
		return "", errors.New("cart is empty")
	}

	var strItemsInCart strings.Builder
	for _, item := range itemsList {
		if item.ProductID == 0 || item.Quantity < 0 {
			return "", errors.New("invalid product ID or quantity in cart item")
		}
		strItemsInCart.WriteString(fmt.Sprintf("%d,%d|", item.ProductID, item.Quantity))
	}

	encodedString := base64.StdEncoding.EncodeToString([]byte(strItemsInCart.String()))
	return encodedString, nil
}

// CookieNames returns a list of cookie names from the request
func CookieNames(r *http.Request) ([]string, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	var names []string
	for _, cookie := range r.Cookies() {
		names = append(names, cookie.Name)
	}
	if len(names) == 0 {
		return nil, errors.New("no cookies found")
	}
	return names, nil
}
