package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/constant"
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
	VerifyToken(ctx context.Context, token string, userID int64) (bool, error)
}

func NewUserUsecase(userRepo repository.IUserRepository, log *logrus.Logger) IUserUsecase {
	return &userUsecase{
		userRepo: userRepo,
		log:      log,
	}
}

func (pu *userUsecase) GetUser(ctx context.Context, req *model.GetUserRequest) (*model.GetUserResponse, error) {
	pu.log.Infof("Fetching user with ID: %d or email: %s", req.UserID, req.Email)
	user, err := pu.userRepo.Get(ctx, req.UserID, req.Email)
	if err != nil {
		pu.log.Errorf("Error fetching user: %v", err)
		return nil, err
	}

	if user == nil {
		pu.log.Errorf("User not found")
		return &model.GetUserResponse{}, nil
	}

	pu.log.Infof("Fetched user: %+v", user)
	return &model.GetUserResponse{
		UserID:       user.UserID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
		Address:      user.Address,
		Role:         user.Role,
		ImageURL:     user.ImageURL,
		Provider:     user.Provider,
		Token:        user.Token,
		IsDeleted:    user.IsDeleted,
		IsVerified:   user.IsVerified,
		TokenExpires: user.TokenExpires.Format(tsCreateTimeLayout),
		CreatedAt:    user.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:    user.UpdatedAt.Format(tsCreateTimeLayout),
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
			Provider:     user.Provider,
			Token:        user.Token,
			IsDeleted:    user.IsDeleted,
			IsVerified:   user.IsVerified,
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
		PasswordHash: user.Password,
		Role:         user.Role,
		ImageURL:     user.ImageURL,
		Provider:     user.Provider,
		IsVerified:   *user.IsVerified,
	}

	if user.Role == "" {
		userData.Role = "Customer"
	}

	if user.ImageURL == "" {
		userData.ImageURL = constant.DEFAULT_USER_IMAGE
	}

	if user.Provider == "" {
		userData.Provider = "System"
	}

	if user.IsVerified == nil {
		userData.IsVerified = false
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
		Provider:     createdUser.Provider,
		Token:        createdUser.Token,
		IsDeleted:    createdUser.IsDeleted,
		IsVerified:   createdUser.IsVerified,
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
	user.IsVerified = rep.IsVerified

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
		Provider:     updatedUser.Provider,
		IsDeleted:    updatedUser.IsDeleted,
		IsVerified:   updatedUser.IsVerified,
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
			Provider:     user.Provider,
			Token:        user.Token,
			IsDeleted:    user.IsDeleted,
			IsVerified:   user.IsVerified,
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

func (pu *userUsecase) VerifyToken(ctx context.Context, token string, userID int64) (bool, error) {
	pu.log.Infof("Verifying token: %s for user: %d", token, userID)
	isValid, err := pu.userRepo.VerifyToken(ctx, token, userID)
	if err != nil {
		pu.log.Errorf("Error verifying token: %v", err)
		return false, err
	}

	user, err := pu.userRepo.Get(ctx, &userID, "")
	if err != nil {
		pu.log.Errorf("Error fetching user for token verification: %v", err)
		return false, err
	}

	if !isValid {
		user.Token = ""
		user.TokenExpires = time.Time{}
		user.IsVerified = false

		_, err = pu.userRepo.Update(ctx, user)
		if err != nil {
			pu.log.Errorf("Error updating user after token verification: %v", err)
			return false, err
		}
	} else {
		user.IsVerified = true

		_, err = pu.userRepo.Update(ctx, user)
		if err != nil {
			pu.log.Errorf("Error updating user after token verification: %v", err)
			return false, err
		}
	}

	pu.log.Infof("Token is valid: %t", isValid)
	return isValid, nil
}
