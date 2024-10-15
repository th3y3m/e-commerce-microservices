package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/category/model"
	"th3y3m/e-commerce-microservices/service/category/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type categoryUsecase struct {
	log          *logrus.Logger
	categoryRepo repository.ICategoryRepository
}

type ICategoryUsecase interface {
	GetCategory(ctx context.Context, req *model.GetCategoryRequest) (*model.GetCategoryResponse, error)
	GetAllCategorys(ctx context.Context) ([]*model.GetCategoryResponse, error)
	CreateCategory(ctx context.Context, req *model.CreateCategoryRequest) (*model.GetCategoryResponse, error)
	UpdateCategory(ctx context.Context, rep *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error)
	DeleteCategory(ctx context.Context, req *model.DeleteCategoryRequest) error
}

func NewCategoryUsecase(categoryRepo repository.ICategoryRepository, log *logrus.Logger) ICategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
		log:          log,
	}
}

func (pu *categoryUsecase) GetCategory(ctx context.Context, req *model.GetCategoryRequest) (*model.GetCategoryResponse, error) {
	pu.log.Infof("Fetching category with ID: %d", req.CategoryID)
	category, err := pu.categoryRepo.Get(ctx, req.CategoryID)
	if err != nil {
		pu.log.Errorf("Error fetching category: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched category: %+v", category)
	return &model.GetCategoryResponse{
		CategoryID: category.CategoryID,
	}, nil
}

func (pu *categoryUsecase) GetAllCategorys(ctx context.Context) ([]*model.GetCategoryResponse, error) {
	pu.log.Infof("Fetching all categorys")
	categorys, err := pu.categoryRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching categorys: %v", err)
		return nil, err
	}

	var categoryResponses []*model.GetCategoryResponse

	for _, category := range categorys {
		categoryResponses = append(categoryResponses, &model.GetCategoryResponse{
			CategoryID:   category.CategoryID,
			CategoryName: category.CategoryName,
			IsDeleted:    category.IsDeleted,
			CreatedAt:    category.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:    category.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d categorys", len(categoryResponses))
	return categoryResponses, nil
}

func (pu *categoryUsecase) CreateCategory(ctx context.Context, category *model.CreateCategoryRequest) (*model.GetCategoryResponse, error) {
	pu.log.Infof("Creating category: %+v", category)
	createdCategory, err := pu.categoryRepo.Create(ctx, &repository.Category{
		CategoryName: category.CategoryName,
	})
	if err != nil {
		pu.log.Errorf("Error creating category: %v", err)
		return nil, err
	}

	pu.log.Infof("Created category: %+v", createdCategory)
	return &model.GetCategoryResponse{
		CategoryID:   createdCategory.CategoryID,
		CategoryName: createdCategory.CategoryName,
		IsDeleted:    createdCategory.IsDeleted,
		CreatedAt:    createdCategory.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:    createdCategory.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *categoryUsecase) DeleteCategory(ctx context.Context, req *model.DeleteCategoryRequest) error {
	pu.log.Infof("Deleting category with ID: %d", req.CategoryID)
	category, err := pu.categoryRepo.Get(ctx, req.CategoryID)
	if err != nil {
		pu.log.Errorf("Error fetching category for deletion: %v", err)
		return err
	}

	category.IsDeleted = true

	_, err = pu.categoryRepo.Update(ctx, category)
	if err != nil {
		pu.log.Errorf("Error updating category for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted category with ID: %d", req.CategoryID)
	return nil
}

func (pu *categoryUsecase) UpdateCategory(ctx context.Context, rep *model.UpdateCategoryRequest) (*model.GetCategoryResponse, error) {
	pu.log.Infof("Updating category with ID: %d", rep.CategoryID)
	category, err := pu.categoryRepo.Get(ctx, rep.CategoryID)
	if err != nil {
		pu.log.Errorf("Error fetching category: %v", err)
		return nil, err
	}

	category.CategoryName = rep.CategoryName
	category.UpdatedAt = time.Now()

	updatedCategory, err := pu.categoryRepo.Update(ctx, category)
	if err != nil {
		pu.log.Errorf("Error updating category: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated category: %+v", updatedCategory)
	return &model.GetCategoryResponse{
		CategoryID:   updatedCategory.CategoryID,
		CategoryName: updatedCategory.CategoryName,
		IsDeleted:    updatedCategory.IsDeleted,
		CreatedAt:    updatedCategory.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:    updatedCategory.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}
