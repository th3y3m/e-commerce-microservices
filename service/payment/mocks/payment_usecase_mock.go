// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "th3y3m/e-commerce-microservices/service/payment/model"

	mock "github.com/stretchr/testify/mock"

	util "th3y3m/e-commerce-microservices/pkg/util"
)

// IPaymentUsecase is an autogenerated mock type for the IPaymentUsecase type
type IPaymentUsecase struct {
	mock.Mock
}

// CreatePayment provides a mock function with given fields: ctx, req
func (_m *IPaymentUsecase) CreatePayment(ctx context.Context, req *model.CreatePaymentRequest) (*model.GetPaymentResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreatePayment")
	}

	var r0 *model.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreatePaymentRequest) (*model.GetPaymentResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreatePaymentRequest) *model.GetPaymentResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.CreatePaymentRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllPayments provides a mock function with given fields: ctx
func (_m *IPaymentUsecase) GetAllPayments(ctx context.Context) ([]*model.GetPaymentResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllPayments")
	}

	var r0 []*model.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.GetPaymentResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.GetPaymentResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPayment provides a mock function with given fields: ctx, req
func (_m *IPaymentUsecase) GetPayment(ctx context.Context, req *model.GetPaymentRequest) (*model.GetPaymentResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetPayment")
	}

	var r0 *model.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetPaymentRequest) (*model.GetPaymentResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetPaymentRequest) *model.GetPaymentResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetPaymentRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPaymentList provides a mock function with given fields: ctx, req
func (_m *IPaymentUsecase) GetPaymentList(ctx context.Context, req *model.GetPaymentsRequest) (*util.PaginatedList[model.GetPaymentResponse], error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetPaymentList")
	}

	var r0 *util.PaginatedList[model.GetPaymentResponse]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetPaymentsRequest) (*util.PaginatedList[model.GetPaymentResponse], error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetPaymentsRequest) *util.PaginatedList[model.GetPaymentResponse]); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*util.PaginatedList[model.GetPaymentResponse])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetPaymentsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdatePayment provides a mock function with given fields: ctx, rep
func (_m *IPaymentUsecase) UpdatePayment(ctx context.Context, rep *model.UpdatePaymentRequest) (*model.GetPaymentResponse, error) {
	ret := _m.Called(ctx, rep)

	if len(ret) == 0 {
		panic("no return value specified for UpdatePayment")
	}

	var r0 *model.GetPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdatePaymentRequest) (*model.GetPaymentResponse, error)); ok {
		return rf(ctx, rep)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdatePaymentRequest) *model.GetPaymentResponse); ok {
		r0 = rf(ctx, rep)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetPaymentResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UpdatePaymentRequest) error); ok {
		r1 = rf(ctx, rep)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIPaymentUsecase creates a new instance of IPaymentUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIPaymentUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IPaymentUsecase {
	mock := &IPaymentUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
