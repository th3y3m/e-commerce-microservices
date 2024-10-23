// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	http "net/http"

	mock "github.com/stretchr/testify/mock"

	model "th3y3m/e-commerce-microservices/service/cart/model"

	util "th3y3m/e-commerce-microservices/pkg/util"
)

// ICartUsecase is an autogenerated mock type for the ICartUsecase type
type ICartUsecase struct {
	mock.Mock
}

// AddProductToShoppingCart provides a mock function with given fields: ctx, userID, productID, quantity
func (_m *ICartUsecase) AddProductToShoppingCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	ret := _m.Called(ctx, userID, productID, quantity)

	if len(ret) == 0 {
		panic("no return value specified for AddProductToShoppingCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int) error); ok {
		r0 = rf(ctx, userID, productID, quantity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ClearShoppingCart provides a mock function with given fields: ctx, userID
func (_m *ICartUsecase) ClearShoppingCart(ctx context.Context, userID int64) error {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for ClearShoppingCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateCart provides a mock function with given fields: ctx, req
func (_m *ICartUsecase) CreateCart(ctx context.Context, req *model.CreateCartRequest) (*model.GetCartResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateCart")
	}

	var r0 *model.GetCartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateCartRequest) (*model.GetCartResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateCartRequest) *model.GetCartResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetCartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.CreateCartRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteCart provides a mock function with given fields: ctx, req
func (_m *ICartUsecase) DeleteCart(ctx context.Context, req *model.DeleteCartRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DeleteCartRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCartInCookie provides a mock function with given fields: w, userId
func (_m *ICartUsecase) DeleteCartInCookie(w http.ResponseWriter, userId int64) error {
	ret := _m.Called(w, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCartInCookie")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, int64) error); ok {
		r0 = rf(w, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteUnitItem provides a mock function with given fields: w, r, productId, userId
func (_m *ICartUsecase) DeleteUnitItem(w http.ResponseWriter, r *http.Request, productId int64, userId int64) error {
	ret := _m.Called(w, r, productId, userId)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUnitItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request, int64, int64) error); ok {
		r0 = rf(w, r, productId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCart provides a mock function with given fields: ctx, req
func (_m *ICartUsecase) GetCart(ctx context.Context, req *model.GetCartRequest) (*model.GetCartResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetCart")
	}

	var r0 *model.GetCartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetCartRequest) (*model.GetCartResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetCartRequest) *model.GetCartResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetCartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetCartRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCartFromCookie provides a mock function with given fields: r, userId
func (_m *ICartUsecase) GetCartFromCookie(r *http.Request, userId int64) ([]util.Item, error) {
	ret := _m.Called(r, userId)

	if len(ret) == 0 {
		panic("no return value specified for GetCartFromCookie")
	}

	var r0 []util.Item
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, int64) ([]util.Item, error)); ok {
		return rf(r, userId)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, int64) []util.Item); ok {
		r0 = rf(r, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]util.Item)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request, int64) error); ok {
		r1 = rf(r, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserCart provides a mock function with given fields: ctx, userID
func (_m *ICartUsecase) GetUserCart(ctx context.Context, userID int64) (*model.GetCartResponse, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetUserCart")
	}

	var r0 *model.GetCartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*model.GetCartResponse, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *model.GetCartResponse); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetCartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NumberOfItemsInCart provides a mock function with given fields: ctx, userID
func (_m *ICartUsecase) NumberOfItemsInCart(ctx context.Context, userID int64) (int, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for NumberOfItemsInCart")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (int, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) int); ok {
		r0 = rf(ctx, userID)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NumberOfItemsInCartCookie provides a mock function with given fields: r, userId
func (_m *ICartUsecase) NumberOfItemsInCartCookie(r *http.Request, userId string) (int, error) {
	ret := _m.Called(r, userId)

	if len(ret) == 0 {
		panic("no return value specified for NumberOfItemsInCartCookie")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (int, error)); ok {
		return rf(r, userId)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) int); ok {
		r0 = rf(r, userId)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveFromCart provides a mock function with given fields: w, r, productId, userId
func (_m *ICartUsecase) RemoveFromCart(w http.ResponseWriter, r *http.Request, productId int64, userId int64) error {
	ret := _m.Called(w, r, productId, userId)

	if len(ret) == 0 {
		panic("no return value specified for RemoveFromCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request, int64, int64) error); ok {
		r0 = rf(w, r, productId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RemoveProductFromShoppingCart provides a mock function with given fields: ctx, userID, productID, quantity
func (_m *ICartUsecase) RemoveProductFromShoppingCart(ctx context.Context, userID int64, productID int64, quantity int) error {
	ret := _m.Called(ctx, userID, productID, quantity)

	if len(ret) == 0 {
		panic("no return value specified for RemoveProductFromShoppingCart")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int) error); ok {
		r0 = rf(ctx, userID, productID, quantity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveCartToCookieHandler provides a mock function with given fields: w, r, productId, userId
func (_m *ICartUsecase) SaveCartToCookieHandler(w http.ResponseWriter, r *http.Request, productId int64, userId int64) error {
	ret := _m.Called(w, r, productId, userId)

	if len(ret) == 0 {
		panic("no return value specified for SaveCartToCookieHandler")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request, int64, int64) error); ok {
		r0 = rf(w, r, productId, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateCart provides a mock function with given fields: ctx, rep
func (_m *ICartUsecase) UpdateCart(ctx context.Context, rep *model.UpdateCartRequest) (*model.GetCartResponse, error) {
	ret := _m.Called(ctx, rep)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCart")
	}

	var r0 *model.GetCartResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateCartRequest) (*model.GetCartResponse, error)); ok {
		return rf(ctx, rep)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateCartRequest) *model.GetCartResponse); ok {
		r0 = rf(ctx, rep)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetCartResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UpdateCartRequest) error); ok {
		r1 = rf(ctx, rep)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICartUsecase creates a new instance of ICartUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICartUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICartUsecase {
	mock := &ICartUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}