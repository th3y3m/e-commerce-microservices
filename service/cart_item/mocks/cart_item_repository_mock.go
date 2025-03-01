// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	repository "th3y3m/e-commerce-microservices/service/cart_item/repository"

	mock "github.com/stretchr/testify/mock"
)

// ICartItemRepository is an autogenerated mock type for the ICartItemRepository type
type ICartItemRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, cartItem
func (_m *ICartItemRepository) Create(ctx context.Context, cartItem *repository.CartItem) (*repository.CartItem, error) {
	ret := _m.Called(ctx, cartItem)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *repository.CartItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.CartItem) (*repository.CartItem, error)); ok {
		return rf(ctx, cartItem)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.CartItem) *repository.CartItem); ok {
		r0 = rf(ctx, cartItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.CartItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.CartItem) error); ok {
		r1 = rf(ctx, cartItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, cartID, productID
func (_m *ICartItemRepository) Delete(ctx context.Context, cartID int64, productID int64) error {
	ret := _m.Called(ctx, cartID, productID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) error); ok {
		r0 = rf(ctx, cartID, productID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, cartID, productID
func (_m *ICartItemRepository) Get(ctx context.Context, cartID int64, productID int64) (*repository.CartItem, error) {
	ret := _m.Called(ctx, cartID, productID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *repository.CartItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) (*repository.CartItem, error)); ok {
		return rf(ctx, cartID, productID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64) *repository.CartItem); ok {
		r0 = rf(ctx, cartID, productID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.CartItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64) error); ok {
		r1 = rf(ctx, cartID, productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetList provides a mock function with given fields: ctx, cartID, productID
func (_m *ICartItemRepository) GetList(ctx context.Context, cartID *int64, productID *int64) ([]*repository.CartItem, error) {
	ret := _m.Called(ctx, cartID, productID)

	if len(ret) == 0 {
		panic("no return value specified for GetList")
	}

	var r0 []*repository.CartItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *int64, *int64) ([]*repository.CartItem, error)); ok {
		return rf(ctx, cartID, productID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *int64, *int64) []*repository.CartItem); ok {
		r0 = rf(ctx, cartID, productID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.CartItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *int64, *int64) error); ok {
		r1 = rf(ctx, cartID, productID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, cartItem
func (_m *ICartItemRepository) Update(ctx context.Context, cartItem *repository.CartItem) (*repository.CartItem, error) {
	ret := _m.Called(ctx, cartItem)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *repository.CartItem
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.CartItem) (*repository.CartItem, error)); ok {
		return rf(ctx, cartItem)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.CartItem) *repository.CartItem); ok {
		r0 = rf(ctx, cartItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.CartItem)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.CartItem) error); ok {
		r1 = rf(ctx, cartItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrCreate provides a mock function with given fields: ctx, cartItem
func (_m *ICartItemRepository) UpdateOrCreate(ctx context.Context, cartItem repository.CartItem) error {
	ret := _m.Called(ctx, cartItem)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrCreate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.CartItem) error); ok {
		r0 = rf(ctx, cartItem)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewICartItemRepository creates a new instance of ICartItemRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICartItemRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICartItemRepository {
	mock := &ICartItemRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
