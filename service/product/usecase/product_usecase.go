package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"th3y3m/e-commerce-microservices/pkg/constant"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/product/model"
	"th3y3m/e-commerce-microservices/service/product/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"
const fallbackPrice = 1000.0 // Fallback price if the calculated price is less than 0

type ProductUsecase struct {
	log         *logrus.Logger
	productRepo repository.IProductRepository
}

type IProductUsecase interface {
	GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.GetProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*model.GetProductResponse, error)
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.GetProductResponse, error)
	UpdateProduct(ctx context.Context, rep *model.UpdateProductRequest) (*model.GetProductResponse, error)
	DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) error
	GetProductList(ctx context.Context, req *model.GetProductsRequest) (*util.PaginatedList[model.GetProductListResponse], error)
	GetProductPriceAfterDiscount(ctx context.Context, req *model.GetProductPriceAfterDiscount) (float64, error)
	UpdateInventory(ctx context.Context, userId, cartId int64) error
}

func NewProductUsecase(productRepo repository.IProductRepository, log *logrus.Logger) IProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
		log:         log,
	}
}

func (pu *ProductUsecase) GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.GetProductResponse, error) {
	pu.log.Infof("Fetching product with ID: %d", req.ProductID)
	product, err := pu.productRepo.Get(ctx, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching product: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched product: %+v", product)
	return &model.GetProductResponse{
		ProductID:   product.ProductID,
		SellerID:    product.SellerID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
		ImageURL:    product.ImageURL,
		CreatedAt:   product.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:   product.UpdatedAt.Format(tsCreateTimeLayout),
		IsDeleted:   product.IsDeleted,
	}, nil
}

func (pu *ProductUsecase) GetAllProducts(ctx context.Context) ([]*model.GetProductResponse, error) {
	pu.log.Info("Fetching all products")
	products, err := pu.productRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all products: %v", err)
		return nil, err
	}

	var productResponses []*model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, &model.GetProductResponse{
			ProductID:   product.ProductID,
			SellerID:    product.SellerID,
			ProductName: product.ProductName,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CategoryID:  product.CategoryID,
			ImageURL:    product.ImageURL,
			CreatedAt:   product.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:   product.UpdatedAt.Format(tsCreateTimeLayout),
			IsDeleted:   product.IsDeleted,
		})
	}

	pu.log.Infof("Fetched %d products", len(productResponses))
	return productResponses, nil
}

func (pu *ProductUsecase) CreateProduct(ctx context.Context, product *model.CreateProductRequest) (*model.GetProductResponse, error) {
	pu.log.Infof("Creating product: %+v", product)
	productData := repository.Product{
		SellerID:    product.SellerID,
		ProductName: product.ProductName,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CategoryID:  product.CategoryID,
		ImageURL:    product.ImageURL,
	}

	createdProduct, err := pu.productRepo.Create(ctx, &productData)
	if err != nil {
		pu.log.Errorf("Error creating product: %v", err)
		return nil, err
	}

	pu.log.Infof("Created product: %+v", createdProduct)
	return &model.GetProductResponse{
		ProductID:   createdProduct.ProductID,
		SellerID:    createdProduct.SellerID,
		ProductName: createdProduct.ProductName,
		Description: createdProduct.Description,
		Price:       createdProduct.Price,
		Quantity:    createdProduct.Quantity,
		CategoryID:  createdProduct.CategoryID,
		ImageURL:    createdProduct.ImageURL,
		CreatedAt:   createdProduct.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:   createdProduct.UpdatedAt.Format(tsCreateTimeLayout),
		IsDeleted:   createdProduct.IsDeleted,
	}, nil
}

func (pu *ProductUsecase) DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) error {
	pu.log.Infof("Deleting product with ID: %d", req.ProductID)
	product, err := pu.productRepo.Get(ctx, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching product for deletion: %v", err)
		return err
	}

	product.IsDeleted = true

	_, err = pu.productRepo.Update(ctx, product)
	if err != nil {
		pu.log.Errorf("Error updating product for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted product with ID: %d", req.ProductID)
	return nil
}

func (pu *ProductUsecase) UpdateProduct(ctx context.Context, rep *model.UpdateProductRequest) (*model.GetProductResponse, error) {
	pu.log.Infof("Updating product with ID: %d", rep.ProductID)
	product, err := pu.productRepo.Get(ctx, rep.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching product for update: %v", err)
		return nil, err
	}

	product.SellerID = rep.SellerID
	product.ProductName = rep.ProductName
	product.Description = rep.Description
	product.Price = rep.Price
	product.Quantity = rep.Quantity
	product.CategoryID = rep.CategoryID
	product.ImageURL = rep.ImageURL
	product.UpdatedAt = time.Now()

	updatedProduct, err := pu.productRepo.Update(ctx, product)
	if err != nil {
		pu.log.Errorf("Error updating product: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated product: %+v", updatedProduct)
	return &model.GetProductResponse{
		ProductID:   updatedProduct.ProductID,
		SellerID:    updatedProduct.SellerID,
		ProductName: updatedProduct.ProductName,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		Quantity:    updatedProduct.Quantity,
		CategoryID:  updatedProduct.CategoryID,
		ImageURL:    updatedProduct.ImageURL,
		CreatedAt:   updatedProduct.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:   updatedProduct.UpdatedAt.Format(tsCreateTimeLayout),
		IsDeleted:   updatedProduct.IsDeleted,
	}, nil
}

func (pu *ProductUsecase) GetProductList(ctx context.Context, req *model.GetProductsRequest) (*util.PaginatedList[model.GetProductListResponse], error) {
	pu.log.Infof("Fetching product list with request: %+v", req)
	products, err := pu.productRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching product list: %v", err)
		return nil, err
	}

	var productResponses []model.GetProductListResponse
	for _, product := range products {
		price, err := pu.GetProductPriceAfterDiscount(ctx, &model.GetProductPriceAfterDiscount{
			ProductID: product.ProductID,
		})
		if err != nil {
			pu.log.Errorf("Error fetching product price after discount: %v", err)
			price = product.Price
		}
		productResponses = append(productResponses, model.GetProductListResponse{
			ProductID:       product.ProductID,
			SellerID:        product.SellerID,
			ProductName:     product.ProductName,
			Description:     product.Description,
			OriginalPrice:   product.Price,
			DiscountedPrice: price,
			Quantity:        product.Quantity,
			CategoryID:      product.CategoryID,
			ImageURL:        product.ImageURL,
			CreatedAt:       product.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:       product.UpdatedAt.Format(tsCreateTimeLayout),
			IsDeleted:       product.IsDeleted,
		})
	}

	list := &util.PaginatedList[model.GetProductListResponse]{
		Items:      productResponses,
		TotalCount: len(productResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d products", len(productResponses))
	return list, nil
}

func (pu *ProductUsecase) GetProductPriceAfterDiscount(ctx context.Context, req *model.GetProductPriceAfterDiscount) (float64, error) {
	pu.log.Infof("Fetching product price after discount with product ID: %d", req.ProductID)

	product, err := pu.productRepo.Get(ctx, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching product: %v", err)
		return 0, err
	}

	productDiscountsRequest := &model.GetProductDiscountsRequest{
		ProductID: &req.ProductID,
	}

	data, err := json.Marshal(productDiscountsRequest)
	if err != nil {
		pu.log.Errorf("Failed to marshal product discounts request: %v", err)
		return product.Price, err
	}

	// Fetch the product discounts
	url := constant.PRODUCT_DISCOUNT_SERVICE
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	request, err := http.NewRequest("GET", url, bytes.NewBuffer(data))
	if err != nil {
		return product.Price, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(request)
	if err != nil {
		return product.Price, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			pu.log.Infof("No discounts found for product with ID: %d", req.ProductID)
			return product.Price, nil
		}
		return product.Price, fmt.Errorf("product discount service returned non-OK status: %d", resp.StatusCode)
	}

	var productDiscounts []*model.ProductDiscount
	if err := json.NewDecoder(resp.Body).Decode(&productDiscounts); err != nil {
		return product.Price, fmt.Errorf("failed to decode product discounts response: %w", err)
	}

	var percentageDiscounts []float64
	var fixedDiscounts []float64

	for _, discount := range productDiscounts {
		url := fmt.Sprintf("%s/%d", constant.DISCOUNT_SERVICE, discount.DiscountID)

		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return product.Price, fmt.Errorf("failed to create discount request: %w", err)
		}

		resp, err := client.Do(request)
		if err != nil {
			return product.Price, fmt.Errorf("failed to execute discount request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return product.Price, fmt.Errorf("discount service returned non-OK status: %d", resp.StatusCode)
		}

		var discountEvent model.GetDiscountResponse
		if err := json.NewDecoder(resp.Body).Decode(&discountEvent); err != nil {
			return product.Price, fmt.Errorf("failed to decode discount response: %w", err)
		}

		startDate, err := time.Parse(tsCreateTimeLayout, discountEvent.StartDate)
		if err != nil {
			pu.log.Errorf("Error parsing start date: %v", err)
			return product.Price, err
		}

		endDate, err := time.Parse(tsCreateTimeLayout, discountEvent.EndDate)
		if err != nil {
			pu.log.Errorf("Error parsing end date: %v", err)
			return product.Price, err
		}

		if !discountEvent.IsDeleted && startDate.Before(time.Now()) && endDate.After(time.Now()) {
			switch discountEvent.DiscountType {
			case "Percentage":
				percentageDiscounts = append(percentageDiscounts, discountEvent.DiscountValue)
			case "Fixed":
				fixedDiscounts = append(fixedDiscounts, discountEvent.DiscountValue)
			}
		}
	}

	// Apply percentage discounts
	for _, discountValue := range percentageDiscounts {
		product.Price -= product.Price * discountValue / 100
	}

	// Apply fixed discounts
	for _, discountValue := range fixedDiscounts {
		product.Price -= discountValue
	}

	if product.Price < 0 {
		product.Price = fallbackPrice
	}

	pu.log.Infof("Fetched product price after discount: %f", product.Price)
	return product.Price, nil
}

func (o *ProductUsecase) UpdateInventory(ctx context.Context, userId, cartId int64) error {
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
		p, err := o.productRepo.Get(ctx, product.ProductID)
		if err != nil {
			return err
		}

		p.Quantity -= product.Quantity
		_, err = o.productRepo.Update(ctx, p)
		if err != nil {
			return err
		}
	}

	return nil
}
