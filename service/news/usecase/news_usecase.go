package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/news/model"
	"th3y3m/e-commerce-microservices/service/news/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type newUsecase struct {
	log     *logrus.Logger
	newRepo repository.INewRepository
}

type INewUsecase interface {
	GetNews(ctx context.Context, req *model.GetNewRequest) (*model.GetNewsResponse, error)
	GetAllNews(ctx context.Context) ([]*model.GetNewsResponse, error)
	CreateNews(ctx context.Context, req *model.CreateNewsRequest) (*model.GetNewsResponse, error)
	UpdateNews(ctx context.Context, rep *model.UpdateNewsRequest) (*model.GetNewsResponse, error)
	DeleteNews(ctx context.Context, req *model.DeleteNewsRequest) error
	GetNewsList(ctx context.Context, req *model.GetNewsRequest) (*util.PaginatedList[model.GetNewsResponse], error)
}

func NewNewsUsecase(newRepo repository.INewRepository, log *logrus.Logger) INewUsecase {
	return &newUsecase{
		newRepo: newRepo,
		log:     log,
	}
}

func (pu *newUsecase) GetNews(ctx context.Context, req *model.GetNewRequest) (*model.GetNewsResponse, error) {
	pu.log.Infof("Fetching new with ID: %d", req.NewsID)
	new, err := pu.newRepo.Get(ctx, req.NewsID)
	if err != nil {
		pu.log.Errorf("Error fetching new: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched new: %+v", new)
	return &model.GetNewsResponse{
		NewsID: new.NewsID,
	}, nil
}

func (pu *newUsecase) GetAllNews(ctx context.Context) ([]*model.GetNewsResponse, error) {
	pu.log.Info("Fetching all news")
	news, err := pu.newRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all news: %v", err)
		return nil, err
	}

	var newResponses []*model.GetNewsResponse
	for _, new := range news {
		newResponses = append(newResponses, &model.GetNewsResponse{
			NewsID:        new.NewsID,
			Title:         new.Title,
			Content:       new.Content,
			PublishedDate: new.PublishedDate.Format(tsCreateTimeLayout),
			AuthorID:      new.AuthorID,
			ImageURL:      new.ImageURL,
			Category:      new.Category,
			IsDeleted:     new.IsDeleted,
			CreatedAt:     new.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:     new.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d news", len(newResponses))
	return newResponses, nil
}

func (pu *newUsecase) CreateNews(ctx context.Context, new *model.CreateNewsRequest) (*model.GetNewsResponse, error) {
	pu.log.Infof("Creating new: %+v", new)
	newData := repository.News{
		Title:         new.Title,
		Content:       new.Content,
		AuthorID:      new.AuthorID,
		ImageURL:      new.ImageURL,
		Category:      new.Category,
		IsDeleted:     new.IsDeleted,
		PublishedDate: time.Now(),
	}

	createdNew, err := pu.newRepo.Create(ctx, &newData)
	if err != nil {
		pu.log.Errorf("Error creating new: %v", err)
		return nil, err
	}

	pu.log.Infof("Created new: %+v", createdNew)
	return &model.GetNewsResponse{
		NewsID:        createdNew.NewsID,
		Title:         createdNew.Title,
		Content:       createdNew.Content,
		PublishedDate: createdNew.PublishedDate.Format(tsCreateTimeLayout),
		AuthorID:      createdNew.AuthorID,
		ImageURL:      createdNew.ImageURL,
		Category:      createdNew.Category,
		IsDeleted:     createdNew.IsDeleted,
		CreatedAt:     createdNew.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     createdNew.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *newUsecase) DeleteNews(ctx context.Context, req *model.DeleteNewsRequest) error {
	pu.log.Infof("Deleting new with ID: %d", req.NewsID)
	new, err := pu.newRepo.Get(ctx, req.NewsID)
	if err != nil {
		pu.log.Errorf("Error fetching new for deletion: %v", err)
		return err
	}

	new.IsDeleted = true

	_, err = pu.newRepo.Update(ctx, new)
	if err != nil {
		pu.log.Errorf("Error updating new for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted new with ID: %d", req.NewsID)
	return nil
}

func (pu *newUsecase) UpdateNews(ctx context.Context, rep *model.UpdateNewsRequest) (*model.GetNewsResponse, error) {
	pu.log.Infof("Updating new with ID: %d", rep.NewsID)
	new, err := pu.newRepo.Get(ctx, rep.NewsID)
	if err != nil {
		pu.log.Errorf("Error fetching new for update: %v", err)
		return nil, err
	}

	new.Title = rep.Title
	new.Content = rep.Content
	new.AuthorID = rep.AuthorID
	new.Category = rep.Category
	new.ImageURL = rep.ImageURL
	new.UpdatedAt = time.Now()
	new.IsDeleted = rep.IsDeleted

	updatedNew, err := pu.newRepo.Update(ctx, new)
	if err != nil {
		pu.log.Errorf("Error updating new: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated new: %+v", updatedNew)
	return &model.GetNewsResponse{
		NewsID:        updatedNew.NewsID,
		Title:         updatedNew.Title,
		Content:       updatedNew.Content,
		PublishedDate: updatedNew.PublishedDate.Format(tsCreateTimeLayout),
		AuthorID:      updatedNew.AuthorID,
		ImageURL:      updatedNew.ImageURL,
		Category:      updatedNew.Category,
		IsDeleted:     updatedNew.IsDeleted,
		CreatedAt:     updatedNew.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     updatedNew.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *newUsecase) GetNewsList(ctx context.Context, req *model.GetNewsRequest) (*util.PaginatedList[model.GetNewsResponse], error) {
	pu.log.Infof("Fetching new list with request: %+v", req)
	news, err := pu.newRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching new list: %v", err)
		return nil, err
	}

	var newResponses []model.GetNewsResponse
	for _, new := range news {
		newResponses = append(newResponses, model.GetNewsResponse{
			NewsID:        new.NewsID,
			Title:         new.Title,
			Content:       new.Content,
			PublishedDate: new.PublishedDate.Format(tsCreateTimeLayout),
			AuthorID:      new.AuthorID,
			ImageURL:      new.ImageURL,
			Category:      new.Category,
			IsDeleted:     new.IsDeleted,
			CreatedAt:     new.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:     new.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	list := &util.PaginatedList[model.GetNewsResponse]{
		Items:      newResponses,
		TotalCount: len(newResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d news", len(newResponses))
	return list, nil
}
