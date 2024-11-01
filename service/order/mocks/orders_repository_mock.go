// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "th3y3m/e-commerce-microservices/service/order/model"

	repository "th3y3m/e-commerce-microservices/service/order/repository"
)

// IOrderRepository is an autogenerated mock type for the IOrderRepository type
type IOrderRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, order
func (_m *IOrderRepository) Create(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Order) (*repository.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Order) *repository.Order); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, orderID
func (_m *IOrderRepository) Delete(ctx context.Context, orderID int64) error {
	ret := _m.Called(ctx, orderID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, orderID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, orderID
func (_m *IOrderRepository) Get(ctx context.Context, orderID int64) (*repository.Order, error) {
	ret := _m.Called(ctx, orderID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*repository.Order, error)); ok {
		return rf(ctx, orderID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *repository.Order); ok {
		r0 = rf(ctx, orderID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, orderID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *IOrderRepository) GetAll(ctx context.Context) ([]*repository.Order, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*repository.Order, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.Order); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetList provides a mock function with given fields: ctx, req
func (_m *IOrderRepository) GetList(ctx context.Context, req *model.GetOrdersRequest) ([]*repository.Order, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetList")
	}

	var r0 []*repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrdersRequest) ([]*repository.Order, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrdersRequest) []*repository.Order); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetOrdersRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, order
func (_m *IOrderRepository) Update(ctx context.Context, order *repository.Order) (*repository.Order, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *repository.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Order) (*repository.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Order) *repository.Order); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// getQuerySearch provides a mock function with given fields: db, req
func (_m *IOrderRepository) getQuerySearch(db *gorm.DB, req *model.GetOrdersRequest) *gorm.DB {
	ret := _m.Called(db, req)

	if len(ret) == 0 {
		panic("no return value specified for getQuerySearch")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.GetOrdersRequest) *gorm.DB); ok {
		r0 = rf(db, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// NewIOrderRepository creates a new instance of IOrderRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIOrderRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IOrderRepository {
	mock := &IOrderRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
