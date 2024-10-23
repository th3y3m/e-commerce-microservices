// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	model "th3y3m/e-commerce-microservices/service/review/model"

	repository "th3y3m/e-commerce-microservices/service/review/repository"
)

// IReviewRepository is an autogenerated mock type for the IReviewRepository type
type IReviewRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, review
func (_m *IReviewRepository) Create(ctx context.Context, review *repository.Review) (*repository.Review, error) {
	ret := _m.Called(ctx, review)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 *repository.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Review) (*repository.Review, error)); ok {
		return rf(ctx, review)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Review) *repository.Review); ok {
		r0 = rf(ctx, review)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Review)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Review) error); ok {
		r1 = rf(ctx, review)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, reviewID
func (_m *IReviewRepository) Delete(ctx context.Context, reviewID int64) error {
	ret := _m.Called(ctx, reviewID)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, reviewID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, reviewID
func (_m *IReviewRepository) Get(ctx context.Context, reviewID int64) (*repository.Review, error) {
	ret := _m.Called(ctx, reviewID)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *repository.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (*repository.Review, error)); ok {
		return rf(ctx, reviewID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) *repository.Review); ok {
		r0 = rf(ctx, reviewID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Review)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, reviewID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: ctx
func (_m *IReviewRepository) GetAll(ctx context.Context) ([]*repository.Review, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAll")
	}

	var r0 []*repository.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*repository.Review, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*repository.Review); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Review)
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
func (_m *IReviewRepository) GetList(ctx context.Context, req *model.GetReviewsRequest) ([]*repository.Review, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetList")
	}

	var r0 []*repository.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetReviewsRequest) ([]*repository.Review, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *model.GetReviewsRequest) []*repository.Review); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*repository.Review)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *model.GetReviewsRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, review
func (_m *IReviewRepository) Update(ctx context.Context, review *repository.Review) (*repository.Review, error) {
	ret := _m.Called(ctx, review)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 *repository.Review
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Review) (*repository.Review, error)); ok {
		return rf(ctx, review)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *repository.Review) *repository.Review); ok {
		r0 = rf(ctx, review)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*repository.Review)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *repository.Review) error); ok {
		r1 = rf(ctx, review)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// getQuerySearch provides a mock function with given fields: db, req
func (_m *IReviewRepository) getQuerySearch(db *gorm.DB, req *model.GetReviewsRequest) *gorm.DB {
	ret := _m.Called(db, req)

	if len(ret) == 0 {
		panic("no return value specified for getQuerySearch")
	}

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func(*gorm.DB, *model.GetReviewsRequest) *gorm.DB); ok {
		r0 = rf(db, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// NewIReviewRepository creates a new instance of IReviewRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewIReviewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *IReviewRepository {
	mock := &IReviewRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
