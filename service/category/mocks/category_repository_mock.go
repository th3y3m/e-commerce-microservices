// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"
	repository "th3y3m/e-commerce-microservices/service/category/repository"

	mock "github.com/stretchr/testify/mock"
)

// ICategoryRepository is an autogenerated mock type for the ICategoryRepository type
type ICategoryRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, category
func (_m *ICategoryRepository) Create(ctx context.Context, category *repository.Category) (*repository.Category, error) {
	ret := _m.Called(ctx, category)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *repository.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Category) (*repository.Category, error)); ok {
		return rf(ctx, category)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Category) *repository.Category); ok {
		r0 = rf(ctx, category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Category) error); ok {
		r1 = rf(ctx, category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, categoryID
func (_m *ICategoryRepository) Delete(ctx context.Context, categoryID int64) error {
	ret := _m.Called(ctx, categoryID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, categoryID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, categoryID
func (_m *ICategoryRepository) Get(ctx context.Context, categoryID int64) (*repository.Category, error) {
	ret := _m.Called(ctx, categoryID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *repository.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*repository.Category, error)); ok {
		return rf(ctx, categoryID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *repository.Category); ok {
		r0 = rf(ctx, categoryID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, categoryID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *ICategoryRepository) GetAll(ctx context.Context) ([]*repository.Category, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*repository.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*repository.Category, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.Category); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, category
func (_m *ICategoryRepository) Update(ctx context.Context, category *repository.Category) (*repository.Category, error) {
	ret := _m.Called(ctx, category)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *repository.Category
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Category) (*repository.Category, error)); ok {
		return rf(ctx, category)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Category) *repository.Category); ok {
		r0 = rf(ctx, category)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Category)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Category) error); ok {
		r1 = rf(ctx, category)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewICategoryRepository creates a new instance of ICategoryRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICategoryRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICategoryRepository {
	mock := &ICategoryRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
