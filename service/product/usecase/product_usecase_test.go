package usecase

import (
	"context"
	"testing"
	"th3y3m/e-commerce-microservices/service/product/mocks"
	"th3y3m/e-commerce-microservices/service/product/model"
	"th3y3m/e-commerce-microservices/service/product/repository"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct(t *testing.T) {
	// Create a new mock instance
	log := logrus.New()
	mockRepo := mocks.NewIProductRepository(t)
	productUsecase := NewProductUsecase(
		mockRepo,
		log,
	)

	// Define the expected behavior
	expectedProduct := &repository.Product{
		ProductID:   1,
		ProductName: "Product 1",
	}
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*repository.Product")).Return(expectedProduct, nil)

	// Call the method
	ctx := context.Background()
	product, err := productUsecase.CreateProduct(ctx, &model.CreateProductRequest{ProductName: "Product 1"})

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ProductName, product.ProductName)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetProduct(t *testing.T) {
	// Create a new mock instance
	log := logrus.New()
	mockRepo := mocks.NewIProductRepository(t)
	productUsecase := NewProductUsecase(
		mockRepo,
		log,
	)

	// Define the expected behavior
	expectedProduct := &repository.Product{
		ProductID:   1,
		ProductName: "Product 1",
	}
	mockRepo.On("Get", mock.Anything, mock.AnythingOfType("int64")).Return(expectedProduct, nil)

	// Call the method
	ctx := context.Background()
	product, err := productUsecase.GetProduct(ctx, &model.GetProductRequest{ProductID: 1})

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ProductID, product.ProductID)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	// Create a new mock instance
	log := logrus.New()
	mockRepo := mocks.NewIProductRepository(t)
	productUsecase := NewProductUsecase(
		mockRepo,
		log,
	)

	// Define the expected behavior
	expectedProducts := []*repository.Product{
		{
			ProductID:   1,
			ProductName: "Product 1",
		},
		{
			ProductID:   2,
			ProductName: "Product 2",
		},
	}
	mockRepo.On("GetAll", mock.Anything).Return(expectedProducts, nil)

	// Call the method
	ctx := context.Background()
	products, err := productUsecase.GetAllProducts(ctx)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, len(expectedProducts), len(products))

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	// Create a new mock instance
	log := logrus.New()
	mockRepo := mocks.NewIProductRepository(t)
	productUsecase := NewProductUsecase(
		mockRepo,
		log,
	)

	product := &repository.Product{
		ProductID:   1,
		ProductName: "Product 1",
	}

	// Define the expected behavior
	expectedProduct := &repository.Product{
		ProductID:   1,
		ProductName: "Product 10",
	}
	mockRepo.On("Get", mock.Anything, expectedProduct.ProductID).Return(product, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*repository.Product")).Return(expectedProduct, nil)

	// Call the method
	ctx := context.Background()
	productRes, err := productUsecase.UpdateProduct(ctx, &model.UpdateProductRequest{ProductID: 1, ProductName: "Product 10"})

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ProductName, productRes.ProductName)

	// Assert that the expectations were met
	mockRepo.AssertExpectations(t)
}
