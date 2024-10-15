package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/product/model"
	"th3y3m/e-commerce-microservices/service/product/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type productUsecase struct {
	log         *logrus.Logger
	productRepo repository.IProductRepository
}

type IProductUsecase interface {
	GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.GetProductResponse, error)
	GetAllProducts(ctx context.Context) ([]*model.GetProductResponse, error)
	CreateProduct(ctx context.Context, req *model.CreateProductRequest) (*model.GetProductResponse, error)
	UpdateProduct(ctx context.Context, rep *model.UpdateProductRequest) (*model.GetProductResponse, error)
	DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) error
	GetProductList(ctx context.Context, req *model.GetProductsRequest) (*util.PaginatedList[model.GetProductResponse], error)
}

func NewProductUsecase(productRepo repository.IProductRepository, log *logrus.Logger) IProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		log:         log,
	}
}

func (pu *productUsecase) GetProduct(ctx context.Context, req *model.GetProductRequest) (*model.GetProductResponse, error) {
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

func (pu *productUsecase) GetAllProducts(ctx context.Context) ([]*model.GetProductResponse, error) {
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

func (pu *productUsecase) CreateProduct(ctx context.Context, product *model.CreateProductRequest) (*model.GetProductResponse, error) {
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

func (pu *productUsecase) DeleteProduct(ctx context.Context, req *model.DeleteProductRequest) error {
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

func (pu *productUsecase) UpdateProduct(ctx context.Context, rep *model.UpdateProductRequest) (*model.GetProductResponse, error) {
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

func (pu *productUsecase) GetProductList(ctx context.Context, req *model.GetProductsRequest) (*util.PaginatedList[model.GetProductResponse], error) {
	pu.log.Infof("Fetching product list with request: %+v", req)
	products, err := pu.productRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching product list: %v", err)
		return nil, err
	}

	var productResponses []model.GetProductResponse
	for _, product := range products {
		productResponses = append(productResponses, model.GetProductResponse{
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

	list := &util.PaginatedList[model.GetProductResponse]{
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
