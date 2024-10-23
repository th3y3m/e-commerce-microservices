package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/cart/model"
	"th3y3m/e-commerce-microservices/service/cart/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type cartUsecase struct {
	log      *logrus.Logger
	cartRepo repository.ICartRepository
}

type ICartUsecase interface {
	GetCart(ctx context.Context, req *model.GetCartRequest) (*model.GetCartResponse, error)
	CreateCart(ctx context.Context, req *model.CreateCartRequest) (*model.GetCartResponse, error)
	UpdateCart(ctx context.Context, rep *model.UpdateCartRequest) (*model.GetCartResponse, error)
	DeleteCart(ctx context.Context, req *model.DeleteCartRequest) error
	GetUserCart(ctx context.Context, userID int64) (*model.GetCartResponse, error)
	AddProductToShoppingCart(ctx context.Context, userID, productID int64, quantity int) error
	RemoveProductFromShoppingCart(ctx context.Context, userID, productID int64, quantity int) error
	ClearShoppingCart(ctx context.Context, userID int64) error
	NumberOfItemsInCart(ctx context.Context, userID int64) (int, error)

	DeleteUnitItem(w http.ResponseWriter, r *http.Request, productId int64) error
	RemoveFromCart(w http.ResponseWriter, r *http.Request, productId int64) error
	GetCartFromCookie(r *http.Request) ([]util.Item, error)
	DeleteCartInCookie(w http.ResponseWriter) error
	NumberOfItemsInCartCookie(r *http.Request) (int, error)
	SaveCartToCookieHandler(w http.ResponseWriter, r *http.Request, productId int64) error
}

func NewCartUsecase(cartRepo repository.ICartRepository, log *logrus.Logger) ICartUsecase {
	return &cartUsecase{
		cartRepo: cartRepo,
		log:      log,
	}
}

func (pu *cartUsecase) GetCart(ctx context.Context, req *model.GetCartRequest) (*model.GetCartResponse, error) {
	pu.log.Infof("Fetching cart with ID: %d", req.CartID)
	cart, err := pu.cartRepo.Get(ctx, &req.CartID)
	if err != nil {
		pu.log.Errorf("Error fetching cart: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched cart: %+v", cart)
	return &model.GetCartResponse{
		CartID:    cart.CartID,
		UserID:    cart.UserID,
		IsDeleted: cart.IsDeleted,
		CreatedAt: cart.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: cart.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *cartUsecase) CreateCart(ctx context.Context, cart *model.CreateCartRequest) (*model.GetCartResponse, error) {
	pu.log.Infof("Creating cart: %+v", cart)
	createdCart, err := pu.cartRepo.Create(ctx, &repository.Cart{
		UserID:    cart.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		pu.log.Errorf("Error creating cart: %v", err)
		return nil, err
	}

	pu.log.Infof("Created cart: %+v", createdCart)
	return &model.GetCartResponse{
		CartID:    createdCart.CartID,
		UserID:    createdCart.UserID,
		IsDeleted: createdCart.IsDeleted,
		CreatedAt: createdCart.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: createdCart.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *cartUsecase) DeleteCart(ctx context.Context, req *model.DeleteCartRequest) error {
	pu.log.Infof("Deleting cart with ID: %d", req.CartID)
	err := pu.cartRepo.Delete(ctx, req.CartID)
	if err != nil {
		pu.log.Errorf("Error deleting cart: %v", err)
		return err
	}

	pu.log.Infof("Deleted cart with ID: %d", req.CartID)
	return nil
}

func (pu *cartUsecase) UpdateCart(ctx context.Context, rep *model.UpdateCartRequest) (*model.GetCartResponse, error) {
	pu.log.Infof("Updating cart with ID: %d", rep.CartID)
	cart, err := pu.cartRepo.Get(ctx, &rep.CartID)
	if err != nil {
		pu.log.Errorf("Error fetching cart: %v", err)
		return nil, err
	}

	cart.IsDeleted = rep.IsDeleted
	cart.UpdatedAt = time.Now()

	updatedCart, err := pu.cartRepo.Update(ctx, cart)
	if err != nil {
		pu.log.Errorf("Error updating cart: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated cart: %+v", updatedCart)
	return &model.GetCartResponse{
		CartID:    updatedCart.CartID,
		UserID:    updatedCart.UserID,
		IsDeleted: updatedCart.IsDeleted,
		CreatedAt: updatedCart.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: updatedCart.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *cartUsecase) GetUserCart(ctx context.Context, userID int64) (*model.GetCartResponse, error) {
	pu.log.Infof("Fetching cart for user with ID: %d", userID)
	cart, err := pu.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		pu.log.Errorf("Error fetching cart: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched cart: %+v", cart)
	return &model.GetCartResponse{
		CartID:    cart.CartID,
		UserID:    cart.UserID,
		IsDeleted: cart.IsDeleted,
		CreatedAt: cart.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: cart.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *cartUsecase) AddProductToShoppingCart(ctx context.Context, userID, productID int64, quantity int) error {
	// Retrieve or create the shopping cart
	cart, err := pu.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		return err
	}

	cartItemReq := model.GetCartItemsRequest{
		CartID: &cart.CartID,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		pu.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		return err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var cartItems []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&cartItems)
	if err != nil {
		return err
	}

	// Create a map to track product quantities
	productList := make(map[int64]int)
	for _, item := range cartItems {
		productList[item.ProductID] = item.Quantity
	}

	// Update the quantity if the product exists, otherwise add it
	if val, ok := productList[productID]; ok {
		productList[productID] = val + quantity
	} else {
		productList[productID] = quantity
	}

	// Update or create the cart item
	cartItem := model.CartItem{
		CartID:    cart.CartID,
		ProductID: productID,
		Quantity:  productList[productID],
	}

	data, err := json.Marshal(cartItem)
	if err != nil {
		pu.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url = constant.CART_ITEM_SERVICE + "/UpdateOrCreateCartItem"
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	// Set the context and execute the request
	client = &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}

func (pu *cartUsecase) RemoveProductFromShoppingCart(ctx context.Context, userID, productID int64, quantity int) error {
	// Retrieve the shopping cart
	cart, err := pu.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		return err
	}

	// Retrieve the cart items
	cartItemReq := model.GetCartItemsRequest{
		CartID: &cart.CartID,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		pu.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		return err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var cartItems []*model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&cartItems)
	if err != nil {
		return err
	}

	// Create a map to track product quantities
	productList := make(map[int64]int)
	for _, item := range cartItems {
		productList[item.ProductID] = item.Quantity
	}

	// Remove the product if it exists
	if val, ok := productList[productID]; ok {
		if val > quantity {
			productList[productID] = val - quantity
		} else {
			delete(productList, productID)
		}
	}

	// Update the cart items
	for _, item := range cartItems {
		if _, ok := productList[item.ProductID]; ok {
			cartItem := model.CartItem{
				CartID:    cart.CartID,
				ProductID: item.ProductID,
				Quantity:  productList[item.ProductID],
			}

			data, err := json.Marshal(cartItem)
			if err != nil {
				pu.log.Errorf("Failed to marshal order data: %v", err)
				return err
			}

			url = constant.CART_ITEM_SERVICE + "/UpdateOrCreateCartItem"
			req, err = http.NewRequest("PUT", url, bytes.NewBuffer(data))
			if err != nil {
				return err
			}

			// Set the context and execute the request
			client = &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Check if the request was successful
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
			}
		} else {
			deleteCartItemRequest := model.DeleteCartItemRequest{
				CartID:    cart.CartID,
				ProductID: item.ProductID,
			}

			data, err := json.Marshal(deleteCartItemRequest)
			if err != nil {
				pu.log.Errorf("Failed to marshal order data: %v", err)
				return err
			}

			url = constant.CART_ITEM_SERVICE
			req, err = http.NewRequest("DELETE", url, bytes.NewBuffer(data))
			if err != nil {
				return err
			}

			// Set the context and execute the request
			client = &http.Client{}
			resp, err = client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// Check if the request was successful
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
			}
		}
	}

	return nil
}

func (pu *cartUsecase) ClearShoppingCart(ctx context.Context, userID int64) error {
	// Retrieve the shopping cart
	cart, err := pu.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		return err
	}

	// Retrieve the cart items
	cartItemReq := model.GetCartItemsRequest{
		CartID: &cart.CartID,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		pu.log.Errorf("Failed to marshal order data: %v", err)
		return err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		return err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var cartItems []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&cartItems)
	if err != nil {
		return err
	}

	// Delete all cart items
	for _, item := range cartItems {
		deleteCartItemRequest := model.DeleteCartItemRequest{
			CartID:    cart.CartID,
			ProductID: item.ProductID,
		}

		data, err := json.Marshal(deleteCartItemRequest)
		if err != nil {
			pu.log.Errorf("Failed to marshal delete cart item request: %v", err)
			return err
		}

		url := constant.CART_ITEM_SERVICE
		req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		// Set the context and execute the request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
		}
	}

	return nil
}

func (pu *cartUsecase) NumberOfItemsInCart(ctx context.Context, userID int64) (int, error) {
	// Retrieve the shopping cart
	cart, err := pu.cartRepo.GetUserCart(ctx, userID)
	if err != nil {
		return 0, err
	}

	// Retrieve the cart items
	cartItemReq := model.GetCartItemsRequest{
		CartID: &cart.CartID,
	}
	cartItemData, err := json.Marshal(cartItemReq)
	if err != nil {
		pu.log.Errorf("Failed to marshal order data: %v", err)
		return 0, err
	}

	url := constant.CART_ITEM_SERVICE
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(cartItemData))
	if err != nil {
		return 0, err
	}

	// Set the context and execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("cart item service returned non-OK status: %d", resp.StatusCode)
	}

	// Decode the response into cart items
	var cartItems []model.GetCartItemResponse
	err = json.NewDecoder(resp.Body).Decode(&cartItems)
	if err != nil {
		return 0, err
	}

	// Calculate the total number of items
	count := 0
	for _, item := range cartItems {
		count += item.Quantity
	}

	return count, nil
}

func (s *cartUsecase) DeleteUnitItem(w http.ResponseWriter, r *http.Request, productId int64) error {
	cookieName := "Cart"
	savedCart, err := r.Cookie(cookieName)
	if err == nil && savedCart != nil {
		// Retrieve cart items from the cookie (map[string]CartItem)
		cartItems, err := util.GetCartFromCookie(savedCart.Value)
		if err != nil {
			return fmt.Errorf("error removing item from cart: %w", err)
		}

		// Check if the item exists in the cart
		if item, exists := cartItems[productId]; exists {
			item.Quantity-- // Decrease the quantity

			// If quantity is zero or less, remove the item from the cart
			if item.Quantity <= 0 {
				delete(cartItems, productId)
			} else {
				// Otherwise, reassign the updated item back to the map
				cartItems[productId] = item
			}
		}

		// Convert the map of CartItems to a slice of CartItem
		var cartItemSlice []util.Item
		for _, item := range cartItems {
			cartItemSlice = append(cartItemSlice, item)
		}

		// Convert the updated cart to a string and save it to the cookie
		strItemsInCart, err := util.ConvertCartToString(cartItemSlice)
		if err != nil {
			return fmt.Errorf("error removing item from cart: %w", err)
		}

		err = util.SaveCartToCookie(w, strItemsInCart)
		if err != nil {
			return fmt.Errorf("error removing item from cart: %w", err)
		}
	} else {
		return fmt.Errorf("error removing item from cart: cookie not found or empty")
	}
	return nil
}

// RemoveFromCart removes a product from the cart.
func (s *cartUsecase) RemoveFromCart(w http.ResponseWriter, r *http.Request, productId int64) error {
	cookieName := "Cart"
	savedCart, err := r.Cookie(cookieName)
	if err == nil && savedCart != nil {
		cartItems, err := util.GetCartFromCookie(savedCart.Value)
		if err != nil {
			return fmt.Errorf("error removing item from cart: %w", err)
		}
		delete(cartItems, productId)

		var cartItemSlice []util.Item
		for _, item := range cartItems {
			cartItemSlice = append(cartItemSlice, item)
		}

		strItemsInCart, err := util.ConvertCartToString(cartItemSlice)
		if err != nil {
			return fmt.Errorf("error removing item from cart: %w", err)
		}
		util.SaveCartToCookie(w, strItemsInCart)
	} else {
		fmt.Println("Error removing item from cart: cookie not found or empty")
	}
	return nil
}

// GetCart retrieves the cart items for a user.
func (s *cartUsecase) GetCartFromCookie(r *http.Request) ([]util.Item, error) {
	var savedCart string
	cartCookie, err := r.Cookie("Cart")

	if err == nil {
		savedCart = cartCookie.Value
	}

	if savedCart != "" {
		cart, err := util.GetCartFromCookie(savedCart)
		if err != nil {
			return []util.Item{}, fmt.Errorf("error getting cart: %w", err)
		}
		var cartItemSlice []util.Item
		for _, item := range cart {
			cartItemSlice = append(cartItemSlice, item)
		}
		return cartItemSlice, nil
	}

	return []util.Item{}, nil
}

// DeleteCartInCookie removes the cart cookie for the user.
func (s *cartUsecase) DeleteCartInCookie(w http.ResponseWriter) error {
	err := util.DeleteCartToCookie(w)
	if err != nil {
		return fmt.Errorf("error deleting cart in cookie: %w", err)
	}
	return nil
}

// NumberOfItemsInCart returns the number of items in the user's cart.
func (s *cartUsecase) NumberOfItemsInCartCookie(r *http.Request) (int, error) {
	var savedCart string
	cartCookie, err := r.Cookie("Cart")
	if err == nil {
		savedCart = cartCookie.Value
	}

	if savedCart != "" {
		cartItems, err := util.GetCartFromCookie(savedCart)
		if err != nil {
			return 0, fmt.Errorf("error getting number of items in cart: %w", err)
		}
		count := 0
		for _, item := range cartItems {
			count += item.Quantity
		}
		return count, nil
	}

	return 0, nil
}

// SaveCartToCookie adds or updates a product in the cart, then saves it to a cookie.
func (s *cartUsecase) SaveCartToCookieHandler(w http.ResponseWriter, r *http.Request, productId int64) error {

	cartItems := make(map[int64]util.Item)

	savedCart, err := r.Cookie("Cart")
	if err == nil && savedCart != nil {
		cartItems, err = util.GetCartFromCookie(savedCart.Value)
		if err != nil {
			return fmt.Errorf("error saving cart to cookie: %w", err)
		}
	} else {
		fmt.Println("Error saving cart to cookie: cookie not found or empty")
	}

	item, exists := cartItems[productId]
	if !exists {
		item = util.Item{
			ProductID: productId,
			Quantity:  1,
		}
	} else {
		item.Quantity++
	}
	cartItems[productId] = item

	var cartItemSlice []util.Item
	for _, item := range cartItems {
		cartItemSlice = append(cartItemSlice, item)
	}

	strItemsInCart, err := util.ConvertCartToString(cartItemSlice)
	if err != nil {
		return fmt.Errorf("error saving cart to cookie: %w", err)
	}
	err = util.SaveCartToCookie(w, strItemsInCart)
	if err != nil {
		return fmt.Errorf("error saving cart to cookie: %w", err)
	}
	return nil
}
