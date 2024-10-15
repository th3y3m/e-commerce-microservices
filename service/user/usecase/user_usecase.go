package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/user/model"
	"th3y3m/e-commerce-microservices/service/user/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type userUsecase struct {
	log      *logrus.Logger
	userRepo repository.IUserRepository
}

type IUserUsecase interface {
	GetUser(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error)
	GetAllUsers(ctx context.Context) ([]*model.GetUserResponse, error)
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.GetUserResponse, error)
	UpdateUser(ctx context.Context, rep *model.UpdateUserRequest) (*model.GetUserResponse, error)
	DeleteUser(ctx context.Context, req *model.DeleteUserRequest) error
	GetUserList(ctx context.Context, req *model.GetUsersRequest) (*util.PaginatedList[model.GetUserResponse], error)
}

func NewUserUsecase(userRepo repository.IUserRepository, log *logrus.Logger) IUserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		log:      log,
	}
}

func (pu *userUsecase) GetUser(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error) {
	pu.log.Infof("Fetching user with ID: %d", req.UserID)
	user, err := pu.userRepo.Get(ctx, req.UserID, req.Email)
	if err != nil {
		pu.log.Errorf("Error fetching user: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched user: %+v", user)
	return &model.GetUserResponse{
		UserID: user.UserID,
	}, nil
}

func (pu *userUsecase) GetAllUsers(ctx context.Context) ([]*model.GetUserResponse, error) {
	pu.log.Info("Fetching all users")
	users, err := pu.userRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all users: %v", err)
		return nil, err
	}

	var userResponses []*model.GetUserResponse
	for _, user := range users {
		userResponses = append(userResponses, &model.GetUserResponse{
			UserID:       user.UserID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			FullName:     user.FullName,
			PhoneNumber:  user.PhoneNumber,
			Address:      user.Address,
			Role:         user.Role,
			ImageURL:     user.ImageURL,
			Token:        user.Token,
			IsDeleted:    user.IsDeleted,
			TokenExpires: user.TokenExpires.Format(tsCreateTimeLayout),
			CreatedAt:    user.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:    user.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d users", len(userResponses))
	return userResponses, nil
}

func (pu *userUsecase) CreateUser(ctx context.Context, user *model.CreateUserRequest) (*model.GetUserResponse, error) {
	pu.log.Infof("Creating user: %+v", user)
	userData := repository.User{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := pu.userRepo.Create(ctx, &userData)
	if err != nil {
		pu.log.Errorf("Error creating user: %v", err)
		return nil, err
	}

	pu.log.Infof("Created user: %+v", createdUser)
	return &model.GetUserResponse{
		UserID:       createdUser.UserID,
		Email:        createdUser.Email,
		PasswordHash: createdUser.PasswordHash,
		FullName:     createdUser.FullName,
		PhoneNumber:  createdUser.PhoneNumber,
		Address:      createdUser.Address,
		Role:         createdUser.Role,
		ImageURL:     createdUser.ImageURL,
		Token:        createdUser.Token,
		IsDeleted:    createdUser.IsDeleted,
		TokenExpires: createdUser.TokenExpires.Format(tsCreateTimeLayout),
		CreatedAt:    createdUser.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:    createdUser.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *userUsecase) DeleteUser(ctx context.Context, req *model.DeleteUserRequest) error {
	pu.log.Infof("Deleting user with ID: %d", req.UserID)
	user, err := pu.userRepo.Get(ctx, &req.UserID, "")
	if err != nil {
		pu.log.Errorf("Error fetching user for deletion: %v", err)
		return err
	}

	user.IsDeleted = true

	_, err = pu.userRepo.Update(ctx, user)
	if err != nil {
		pu.log.Errorf("Error updating user for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted user with ID: %d", req.UserID)
	return nil
}

func (pu *userUsecase) UpdateUser(ctx context.Context, rep *model.UpdateUserRequest) (*model.GetUserResponse, error) {
	pu.log.Infof("Updating user with ID: %d", rep.UserID)
	user, err := pu.userRepo.Get(ctx, &rep.UserID, rep.Email)
	if err != nil {
		pu.log.Errorf("Error fetching user for update: %v", err)
		return nil, err
	}

	user.Email = rep.Email
	user.PasswordHash = rep.PasswordHash
	user.FullName = rep.FullName
	user.PhoneNumber = rep.PhoneNumber
	user.Address = rep.Address
	user.Role = rep.Role
	user.Token = rep.Token
	user.TokenExpires = rep.TokenExpires
	user.ImageURL = rep.ImageURL
	user.UpdatedAt = time.Now()
	user.IsDeleted = rep.IsDeleted

	updatedUser, err := pu.userRepo.Update(ctx, user)
	if err != nil {
		pu.log.Errorf("Error updating user: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated user: %+v", updatedUser)
	return &model.GetUserResponse{
		UserID:       updatedUser.UserID,
		Email:        updatedUser.Email,
		PasswordHash: updatedUser.PasswordHash,
		FullName:     updatedUser.FullName,
		PhoneNumber:  updatedUser.PhoneNumber,
		Address:      updatedUser.Address,
		Role:         updatedUser.Role,
		ImageURL:     updatedUser.ImageURL,
		Token:        updatedUser.Token,
		IsDeleted:    updatedUser.IsDeleted,
		TokenExpires: updatedUser.TokenExpires.Format(tsCreateTimeLayout),
		CreatedAt:    updatedUser.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:    updatedUser.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *userUsecase) GetUserList(ctx context.Context, req *model.GetUsersRequest) (*util.PaginatedList[model.GetUserResponse], error) {
	pu.log.Infof("Fetching user list with request: %+v", req)
	users, err := pu.userRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching user list: %v", err)
		return nil, err
	}

	var userResponses []model.GetUserResponse
	for _, user := range users {
		userResponses = append(userResponses, model.GetUserResponse{
			UserID:       user.UserID,
			Email:        user.Email,
			PasswordHash: user.PasswordHash,
			FullName:     user.FullName,
			PhoneNumber:  user.PhoneNumber,
			Address:      user.Address,
			Role:         user.Role,
			ImageURL:     user.ImageURL,
			Token:        user.Token,
			IsDeleted:    user.IsDeleted,
			TokenExpires: user.TokenExpires.Format(tsCreateTimeLayout),
			CreatedAt:    user.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:    user.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	list := &util.PaginatedList[model.GetUserResponse]{
		Items:      userResponses,
		TotalCount: len(userResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d users", len(userResponses))
	return list, nil
}
