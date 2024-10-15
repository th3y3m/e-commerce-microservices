package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/courier/model"
	"th3y3m/e-commerce-microservices/service/courier/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type courierUsecase struct {
	log         *logrus.Logger
	courierRepo repository.ICourierRepository
}

type ICourierUsecase interface {
	GetCourier(ctx context.Context, req *model.GetCourierRequest) (*model.GetCourierResponse, error)
	GetAllCouriers(ctx context.Context) ([]*model.GetCourierResponse, error)
	CreateCourier(ctx context.Context, req *model.CreateCourierRequest) (*model.GetCourierResponse, error)
	UpdateCourier(ctx context.Context, rep *model.UpdateCourierRequest) (*model.GetCourierResponse, error)
	DeleteCourier(ctx context.Context, req *model.DeleteCourierRequest) error
}

func NewCourierUsecase(courierRepo repository.ICourierRepository, log *logrus.Logger) ICourierUsecase {
	return &courierUsecase{
		courierRepo: courierRepo,
		log:         log,
	}
}

func (pu *courierUsecase) GetCourier(ctx context.Context, req *model.GetCourierRequest) (*model.GetCourierResponse, error) {
	pu.log.Infof("Fetching courier with ID: %d", req.CourierID)
	courier, err := pu.courierRepo.Get(ctx, req.CourierID)
	if err != nil {
		pu.log.Errorf("Error fetching courier: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched courier: %+v", courier)
	return &model.GetCourierResponse{
		CourierID: courier.CourierID,
	}, nil
}

func (pu *courierUsecase) GetAllCouriers(ctx context.Context) ([]*model.GetCourierResponse, error) {
	pu.log.Infof("Fetching all couriers")
	couriers, err := pu.courierRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching couriers: %v", err)
		return nil, err
	}

	var courierResponses []*model.GetCourierResponse
	for _, courier := range couriers {
		courierResponses = append(courierResponses, &model.GetCourierResponse{
			CourierID:   courier.CourierID,
			CourierName: courier.CourierName,
			IsDeleted:   courier.IsDeleted,
			CreatedAt:   courier.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:   courier.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d couriers", len(courierResponses))
	return courierResponses, nil
}

func (pu *courierUsecase) CreateCourier(ctx context.Context, courier *model.CreateCourierRequest) (*model.GetCourierResponse, error) {
	pu.log.Infof("Creating courier: %+v", courier)
	createdCourier, err := pu.courierRepo.Create(ctx, &repository.Courier{
		CourierName: courier.CourierName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		pu.log.Errorf("Error creating courier: %v", err)
		return nil, err
	}

	pu.log.Infof("Created courier: %+v", createdCourier)
	return &model.GetCourierResponse{
		CourierID:   createdCourier.CourierID,
		CourierName: createdCourier.CourierName,
		IsDeleted:   createdCourier.IsDeleted,
		CreatedAt:   createdCourier.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:   createdCourier.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *courierUsecase) DeleteCourier(ctx context.Context, req *model.DeleteCourierRequest) error {
	pu.log.Infof("Deleting courier with ID: %d", req.CourierID)
	courier, err := pu.courierRepo.Get(ctx, req.CourierID)
	if err != nil {
		pu.log.Errorf("Error fetching courier for deletion: %v", err)
		return err
	}

	courier.IsDeleted = true

	_, err = pu.courierRepo.Update(ctx, courier)
	if err != nil {
		pu.log.Errorf("Error updating courier for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted courier with ID: %d", req.CourierID)
	return nil
}

func (pu *courierUsecase) UpdateCourier(ctx context.Context, rep *model.UpdateCourierRequest) (*model.GetCourierResponse, error) {
	pu.log.Infof("Updating courier with ID: %d", rep.CourierID)
	courier, err := pu.courierRepo.Get(ctx, rep.CourierID)
	if err != nil {
		pu.log.Errorf("Error fetching courier: %v", err)
		return nil, err
	}

	courier.CourierName = rep.CourierName
	courier.UpdatedAt = time.Now()

	updatedCourier, err := pu.courierRepo.Update(ctx, courier)
	if err != nil {
		pu.log.Errorf("Error updating courier: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated courier: %+v", updatedCourier)
	return &model.GetCourierResponse{
		CourierID:   updatedCourier.CourierID,
		CourierName: updatedCourier.CourierName,
		IsDeleted:   updatedCourier.IsDeleted,
		CreatedAt:   updatedCourier.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:   updatedCourier.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}
