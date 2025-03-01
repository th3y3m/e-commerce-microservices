// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "th3y3m/e-commerce-microservices/service/order/model"

	mock "github.com/stretchr/testify/mock"

	util "th3y3m/e-commerce-microservices/pkg/util"
)

// IOrderUsecase is an autogenerated mock type for the IOrderUsecase type
type IOrderUsecase struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, req
func (_m *IOrderUsecase) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.GetOrderResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 *model.GetOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateOrderRequest) (*model.GetOrderResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateOrderRequest) *model.GetOrderResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.CreateOrderRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteOrder provides a mock function with given fields: ctx, req
func (_m *IOrderUsecase) DeleteOrder(ctx context.Context, req *model.DeleteOrderRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOrder")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DeleteOrderRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllOrders provides a mock function with given fields: ctx
func (_m *IOrderUsecase) GetAllOrders(ctx context.Context) ([]*model.GetOrderResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllOrders")
	}

	var r0 []*model.GetOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.GetOrderResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.GetOrderResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GetOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrder provides a mock function with given fields: ctx, req
func (_m *IOrderUsecase) GetOrder(ctx context.Context, req *model.GetOrderRequest) (*model.GetOrderResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetOrder")
	}

	var r0 *model.GetOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrderRequest) (*model.GetOrderResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrderRequest) *model.GetOrderResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetOrderRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderList provides a mock function with given fields: ctx, req
func (_m *IOrderUsecase) GetOrderList(ctx context.Context, req *model.GetOrdersRequest) (*util.PaginatedList[model.GetOrderResponse], error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetOrderList")
	}

	var r0 *util.PaginatedList[model.GetOrderResponse]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrdersRequest) (*util.PaginatedList[model.GetOrderResponse], error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrdersRequest) *util.PaginatedList[model.GetOrderResponse]); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*util.PaginatedList[model.GetOrderResponse])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetOrdersRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PlaceOrder provides a mock function with given fields: ctx, userId, cartId, CourierID, VoucherID, paymentMethod, shipAddress, freight
func (_m *IOrderUsecase) PlaceOrder(ctx context.Context, userId int64, cartId int64, CourierID int64, VoucherID int64, paymentMethod string, shipAddress string, freight float64) (string, error) {
	ret := _m.Called(ctx, userId, cartId, CourierID, VoucherID, paymentMethod, shipAddress, freight)

	if len(ret) == 0 {
		panic("no return value specified for PlaceOrder")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64, int64, string, string, float64) (string, error)); ok {
		return rf(ctx, userId, cartId, CourierID, VoucherID, paymentMethod, shipAddress, freight)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64, int64, string, string, float64) string); ok {
		r0 = rf(ctx, userId, cartId, CourierID, VoucherID, paymentMethod, shipAddress, freight)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, int64, int64, string, string, float64) error); ok {
		r1 = rf(ctx, userId, cartId, CourierID, VoucherID, paymentMethod, shipAddress, freight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessOrder provides a mock function with given fields: ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight
func (_m *IOrderUsecase) ProcessOrder(ctx context.Context, userId int64, cartId int64, CourierID int64, VoucherID int64, shipAddress string, paymentMethod string, freight float64) (*model.GetOrderResponse, error) {
	ret := _m.Called(ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight)

	if len(ret) == 0 {
		panic("no return value specified for ProcessOrder")
	}

	var r0 *model.GetOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64, int64, string, string, float64) (*model.GetOrderResponse, error)); ok {
		return rf(ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64, int64, int64, int64, string, string, float64) *model.GetOrderResponse); ok {
		r0 = rf(ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64, int64, int64, int64, string, string, float64) error); ok {
		r1 = rf(ctx, userId, cartId, CourierID, VoucherID, shipAddress, paymentMethod, freight)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessPayment provides a mock function with given fields: ctx, order, paymentMethod
func (_m *IOrderUsecase) ProcessPayment(ctx context.Context, order *model.GetOrderResponse, paymentMethod string) (string, error) {
	ret := _m.Called(ctx, order, paymentMethod)

	if len(ret) == 0 {
		panic("no return value specified for ProcessPayment")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrderResponse, string) (string, error)); ok {
		return rf(ctx, order, paymentMethod)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetOrderResponse, string) string); ok {
		r0 = rf(ctx, order, paymentMethod)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetOrderResponse, string) error); ok {
		r1 = rf(ctx, order, paymentMethod)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrder provides a mock function with given fields: ctx, rep
func (_m *IOrderUsecase) UpdateOrder(ctx context.Context, rep *model.UpdateOrderRequest) (*model.GetOrderResponse, error) {
	ret := _m.Called(ctx, rep)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrder")
	}

	var r0 *model.GetOrderResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateOrderRequest) (*model.GetOrderResponse, error)); ok {
		return rf(ctx, rep)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateOrderRequest) *model.GetOrderResponse); ok {
		r0 = rf(ctx, rep)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetOrderResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UpdateOrderRequest) error); ok {
		r1 = rf(ctx, rep)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIOrderUsecase creates a new instance of IOrderUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIOrderUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IOrderUsecase {
	mock := &IOrderUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
