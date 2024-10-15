package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/cart_item/model"
	"th3y3m/e-commerce-microservices/service/cart_item/repository"

	"github.com/sirupsen/logrus"
)

type cartItemUsecase struct {
	log          *logrus.Logger
	cartItemRepo repository.ICartItemRepository
}

type ICartItemUsecase interface {
	GetCartItem(ctx context.Context, req *model.GetCartItemRequest) (*model.GetCartItemResponse, error)
	CreateCartItem(ctx context.Context, req *model.CreateCartItemRequest) (*model.GetCartItemResponse, error)
	UpdateCartItem(ctx context.Context, rep *model.UpdateCartItemRequest) (*model.GetCartItemResponse, error)
	DeleteCartItem(ctx context.Context, req *model.DeleteCartItemRequest) error
	GetCartItemList(ctx context.Context, req *model.GetCartItemsRequest) ([]*model.GetCartItemResponse, error)
}

func NewCartItemUsecase(cartItemRepo repository.ICartItemRepository, log *logrus.Logger) ICartItemUsecase {
	return &cartItemUsecase{
		cartItemRepo: cartItemRepo,
		log:          log,
	}
}

func (pu *cartItemUsecase) GetCartItem(ctx context.Context, req *model.GetCartItemRequest) (*model.GetCartItemResponse, error) {
	pu.log.Infof("Fetching cartItem with cartID: %d and productID: %d", req.CartID, req.ProductID)
	cartItem, err := pu.cartItemRepo.Get(ctx, req.CartID, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching cartItem: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched cartItem: %+v", cartItem)
	return &model.GetCartItemResponse{
		CartID:    cartItem.CartID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	}, nil
}

func (pu *cartItemUsecase) CreateCartItem(ctx context.Context, cartItem *model.CreateCartItemRequest) (*model.GetCartItemResponse, error) {
	pu.log.Infof("Creating cartItem: %+v", cartItem)
	createdCartItem, err := pu.cartItemRepo.Create(ctx, &repository.CartItem{
		CartID:    cartItem.CartID,
		ProductID: cartItem.ProductID,
		Quantity:  cartItem.Quantity,
	})
	if err != nil {
		pu.log.Errorf("Error creating cartItem: %v", err)
		return nil, err
	}

	pu.log.Infof("Created cartItem: %+v", createdCartItem)
	return &model.GetCartItemResponse{
		CartID:    createdCartItem.CartID,
		ProductID: createdCartItem.ProductID,
		Quantity:  createdCartItem.Quantity,
	}, nil
}

func (pu *cartItemUsecase) UpdateCartItem(ctx context.Context, req *model.UpdateCartItemRequest) (*model.GetCartItemResponse, error) {
	pu.log.Infof("Updating cartItem with cartID: %d and productID: %d", req.CartID, req.ProductID)
	updatedCartItem, err := pu.cartItemRepo.Update(ctx, &repository.CartItem{
		CartID:    req.CartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		pu.log.Errorf("Error updating cartItem: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated cartItem: %+v", updatedCartItem)
	return &model.GetCartItemResponse{
		CartID:    updatedCartItem.CartID,
		ProductID: updatedCartItem.ProductID,
		Quantity:  updatedCartItem.Quantity,
	}, nil
}

func (pu *cartItemUsecase) DeleteCartItem(ctx context.Context, req *model.DeleteCartItemRequest) error {
	pu.log.Infof("Deleting cartItem with cartID: %d and productID: %d", req.CartID, req.ProductID)
	err := pu.cartItemRepo.Delete(ctx, req.CartID, req.ProductID)
	if err != nil {
		pu.log.Errorf("Error deleting cartItem: %v", err)
		return err
	}

	pu.log.Infof("Deleted cartItem with cartID: %d and productID: %d", req.CartID, req.ProductID)
	return nil
}

func (pu *cartItemUsecase) GetCartItemList(ctx context.Context, req *model.GetCartItemsRequest) ([]*model.GetCartItemResponse, error) {
	pu.log.Infof("Fetching cartItems with cartID: %d and productID: %d", req.CartID, req.ProductID)
	cartItems, err := pu.cartItemRepo.GetList(ctx, &req.CartID, &req.ProductID)
	if err != nil {
		pu.log.Errorf("Error fetching cartItems: %v", err)
		return nil, err
	}

	var cartItemsResp []*model.GetCartItemResponse
	for _, cartItem := range cartItems {
		cartItemsResp = append(cartItemsResp, &model.GetCartItemResponse{
			CartID:    cartItem.CartID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
		})
	}

	pu.log.Infof("Fetched cartItems: %+v", cartItemsResp)
	return cartItemsResp, nil
}
