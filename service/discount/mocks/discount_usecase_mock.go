// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	model "th3y3m/e-commerce-microservices/service/discount/model"

	mock "github.com/stretchr/testify/mock"
)

// IDiscountUsecase is an autogenerated mock type for the IDiscountUsecase type
type IDiscountUsecase struct {
	mock.Mock
}

// CreateDiscount provides a mock function with given fields: ctx, req
func (_m *IDiscountUsecase) CreateDiscount(ctx context.Context, req *model.CreateDiscountRequest) (*model.GetDiscountResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateDiscount")
	}

	var r0 *model.GetDiscountResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateDiscountRequest) (*model.GetDiscountResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.CreateDiscountRequest) *model.GetDiscountResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetDiscountResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.CreateDiscountRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteDiscount provides a mock function with given fields: ctx, req
func (_m *IDiscountUsecase) DeleteDiscount(ctx context.Context, req *model.DeleteDiscountRequest) error {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for DeleteDiscount")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.DeleteDiscountRequest) error); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllDiscounts provides a mock function with given fields: ctx
func (_m *IDiscountUsecase) GetAllDiscounts(ctx context.Context) ([]*model.GetDiscountResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllDiscounts")
	}

	var r0 []*model.GetDiscountResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*model.GetDiscountResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*model.GetDiscountResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.GetDiscountResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDiscount provides a mock function with given fields: ctx, req
func (_m *IDiscountUsecase) GetDiscount(ctx context.Context, req *model.GetDiscountRequest) (*model.GetDiscountResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetDiscount")
	}

	var r0 *model.GetDiscountResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetDiscountRequest) (*model.GetDiscountResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetDiscountRequest) *model.GetDiscountResponse); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetDiscountResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetDiscountRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateDiscount provides a mock function with given fields: ctx, rep
func (_m *IDiscountUsecase) UpdateDiscount(ctx context.Context, rep *model.UpdateDiscountRequest) (*model.GetDiscountResponse, error) {
	ret := _m.Called(ctx, rep)

	if len(ret) == 0 {
		panic("no return value specified for UpdateDiscount")
	}

	var r0 *model.GetDiscountResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateDiscountRequest) (*model.GetDiscountResponse, error)); ok {
		return rf(ctx, rep)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.UpdateDiscountRequest) *model.GetDiscountResponse); ok {
		r0 = rf(ctx, rep)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.GetDiscountResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.UpdateDiscountRequest) error); ok {
		r1 = rf(ctx, rep)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewIDiscountUsecase creates a new instance of IDiscountUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIDiscountUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *IDiscountUsecase {
	mock := &IDiscountUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
