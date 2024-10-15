package usecase

import (
	"context"
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
