package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/review/model"
	"th3y3m/e-commerce-microservices/service/review/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type reviewUsecase struct {
	log        *logrus.Logger
	reviewRepo repository.IReviewRepository
}

type IReviewUsecase interface {
	GetReview(ctx context.Context, req *model.GetReviewRequest) (*model.GetReviewResponse, error)
	GetAllReviews(ctx context.Context) ([]*model.GetReviewResponse, error)
	CreateReview(ctx context.Context, req *model.CreateReviewRequest) (*model.GetReviewResponse, error)
	UpdateReview(ctx context.Context, rep *model.UpdateReviewRequest) (*model.GetReviewResponse, error)
	DeleteReview(ctx context.Context, req *model.DeleteReviewRequest) error
	GetReviewList(ctx context.Context, req *model.GetReviewsRequest) (*util.PaginatedList[model.GetReviewResponse], error)
}

func NewReviewUsecase(reviewRepo repository.IReviewRepository, log *logrus.Logger) IReviewUsecase {
	return &reviewUsecase{
		reviewRepo: reviewRepo,
		log:        log,
	}
}

func (pu *reviewUsecase) GetReview(ctx context.Context, req *model.GetReviewRequest) (*model.GetReviewResponse, error) {
	pu.log.Infof("Fetching review with ID: %d", req.ReviewID)
	review, err := pu.reviewRepo.Get(ctx, req.ReviewID)
	if err != nil {
		pu.log.Errorf("Error fetching review: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched review: %+v", review)
	return &model.GetReviewResponse{
		ReviewID: review.ReviewID,
	}, nil
}

func (pu *reviewUsecase) GetAllReviews(ctx context.Context) ([]*model.GetReviewResponse, error) {
	pu.log.Info("Fetching all reviews")
	reviews, err := pu.reviewRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all reviews: %v", err)
		return nil, err
	}

	var reviewResponses []*model.GetReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, &model.GetReviewResponse{
			ReviewID:  review.ReviewID,
			ProductID: review.ProductID,
			UserID:    review.UserID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			IsDeleted: review.IsDeleted,
			CreatedAt: review.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt: review.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d reviews", len(reviewResponses))
	return reviewResponses, nil
}

func (pu *reviewUsecase) CreateReview(ctx context.Context, review *model.CreateReviewRequest) (*model.GetReviewResponse, error) {
	pu.log.Infof("Creating review: %+v", review)
	reviewData := repository.Review{
		ProductID: review.ProductID,
		UserID:    review.UserID,
		Rating:    review.Rating,
		Comment:   review.Comment,
	}

	createdReview, err := pu.reviewRepo.Create(ctx, &reviewData)
	if err != nil {
		pu.log.Errorf("Error creating review: %v", err)
		return nil, err
	}

	pu.log.Infof("Created review: %+v", createdReview)
	return &model.GetReviewResponse{
		ReviewID:  createdReview.ReviewID,
		ProductID: createdReview.ProductID,
		UserID:    createdReview.UserID,
		Rating:    createdReview.Rating,
		Comment:   createdReview.Comment,
		IsDeleted: createdReview.IsDeleted,
		CreatedAt: createdReview.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: createdReview.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *reviewUsecase) DeleteReview(ctx context.Context, req *model.DeleteReviewRequest) error {
	pu.log.Infof("Deleting review with ID: %d", req.ReviewID)
	review, err := pu.reviewRepo.Get(ctx, req.ReviewID)
	if err != nil {
		pu.log.Errorf("Error fetching review for deletion: %v", err)
		return err
	}

	review.IsDeleted = true

	_, err = pu.reviewRepo.Update(ctx, review)
	if err != nil {
		pu.log.Errorf("Error updating review for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted review with ID: %d", req.ReviewID)
	return nil
}

func (pu *reviewUsecase) UpdateReview(ctx context.Context, rep *model.UpdateReviewRequest) (*model.GetReviewResponse, error) {
	pu.log.Infof("Updating review with ID: %d", rep.ReviewID)
	review, err := pu.reviewRepo.Get(ctx, rep.ReviewID)
	if err != nil {
		pu.log.Errorf("Error fetching review for update: %v", err)
		return nil, err
	}

	review.Rating = rep.Rating
	review.Comment = rep.Comment
	review.UpdatedAt = time.Now()
	review.IsDeleted = rep.IsDeleted

	updatedReview, err := pu.reviewRepo.Update(ctx, review)
	if err != nil {
		pu.log.Errorf("Error updating review: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated review: %+v", updatedReview)
	return &model.GetReviewResponse{
		ReviewID:  updatedReview.ReviewID,
		ProductID: updatedReview.ProductID,
		UserID:    updatedReview.UserID,
		Rating:    updatedReview.Rating,
		Comment:   updatedReview.Comment,
		IsDeleted: updatedReview.IsDeleted,
		CreatedAt: updatedReview.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt: updatedReview.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *reviewUsecase) GetReviewList(ctx context.Context, req *model.GetReviewsRequest) (*util.PaginatedList[model.GetReviewResponse], error) {
	pu.log.Infof("Fetching review list with request: %+v", req)
	reviews, err := pu.reviewRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching review list: %v", err)
		return nil, err
	}

	var reviewResponses []model.GetReviewResponse
	for _, review := range reviews {
		reviewResponses = append(reviewResponses, model.GetReviewResponse{
			ReviewID:  review.ReviewID,
			ProductID: review.ProductID,
			UserID:    review.UserID,
			Rating:    review.Rating,
			Comment:   review.Comment,
			IsDeleted: review.IsDeleted,
			CreatedAt: review.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt: review.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	list := &util.PaginatedList[model.GetReviewResponse]{
		Items:      reviewResponses,
		TotalCount: len(reviewResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d reviews", len(reviewResponses))
	return list, nil
}
