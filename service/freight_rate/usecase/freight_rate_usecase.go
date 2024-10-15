package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/service/freight_rate/model"
	"th3y3m/e-commerce-microservices/service/freight_rate/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type freightRateUsecase struct {
	log             *logrus.Logger
	freightRateRepo repository.IFreightRateRepository
}

type IFreightRateUsecase interface {
	GetFreightRate(ctx context.Context, req *model.GetFreightRateRequest) (*model.GetFreightRateResponse, error)
	GetAllFreightRates(ctx context.Context) ([]*model.GetFreightRateResponse, error)
	CreateFreightRate(ctx context.Context, req *model.CreateFreightRateRequest) (*model.GetFreightRateResponse, error)
	UpdateFreightRate(ctx context.Context, rep *model.UpdateFreightRateRequest) (*model.GetFreightRateResponse, error)
	DeleteFreightRate(ctx context.Context, req *model.DeleteFreightRateRequest) error
}

func NewFreightRateUsecase(freightRateRepo repository.IFreightRateRepository, log *logrus.Logger) IFreightRateUsecase {
	return &freightRateUsecase{
		freightRateRepo: freightRateRepo,
		log:             log,
	}
}

func (pu *freightRateUsecase) GetFreightRate(ctx context.Context, req *model.GetFreightRateRequest) (*model.GetFreightRateResponse, error) {
	pu.log.Infof("Fetching freightRate with ID: %d", req.FreightRateID)
	freightRate, err := pu.freightRateRepo.Get(ctx, req.FreightRateID)
	if err != nil {
		pu.log.Errorf("Error fetching freightRate: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched freightRate: %+v", freightRate)
	return &model.GetFreightRateResponse{
		FreightRateID: freightRate.FreightRateID,
		CourierID:     freightRate.CourierID,
		DistanceMinKM: freightRate.DistanceMinKM,
		DistanceMaxKM: freightRate.DistanceMaxKM,
		CostPerKM:     freightRate.CostPerKM,
		IsDeleted:     freightRate.IsDeleted,
		CreatedAt:     freightRate.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     freightRate.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *freightRateUsecase) GetAllFreightRates(ctx context.Context) ([]*model.GetFreightRateResponse, error) {

	pu.log.Infof("Fetching all freightRates")
	freightRates, err := pu.freightRateRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching freightRates: %v", err)
		return nil, err
	}

	var freightRateResponses []*model.GetFreightRateResponse
	for _, freightRate := range freightRates {
		freightRateResponses = append(freightRateResponses, &model.GetFreightRateResponse{
			FreightRateID: freightRate.FreightRateID,
			CourierID:     freightRate.CourierID,
			DistanceMinKM: freightRate.DistanceMinKM,
			DistanceMaxKM: freightRate.DistanceMaxKM,
			CostPerKM:     freightRate.CostPerKM,
			IsDeleted:     freightRate.IsDeleted,
			CreatedAt:     freightRate.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:     freightRate.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d freightRates", len(freightRateResponses))
	return freightRateResponses, nil
}

func (pu *freightRateUsecase) CreateFreightRate(ctx context.Context, freightRate *model.CreateFreightRateRequest) (*model.GetFreightRateResponse, error) {
	pu.log.Infof("Creating freightRate: %+v", freightRate)
	newFreightRate := &repository.FreightRate{
		CourierID:     freightRate.CourierID,
		DistanceMinKM: freightRate.DistanceMinKM,
		DistanceMaxKM: freightRate.DistanceMaxKM,
		CostPerKM:     freightRate.CostPerKM,
		IsDeleted:     false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdFreightRate, err := pu.freightRateRepo.Create(ctx, newFreightRate)
	if err != nil {
		pu.log.Errorf("Error creating freightRate: %v", err)
		return nil, err
	}

	pu.log.Infof("Created freightRate: %+v", createdFreightRate)
	return &model.GetFreightRateResponse{
		FreightRateID: createdFreightRate.FreightRateID,
		CourierID:     createdFreightRate.CourierID,
		DistanceMinKM: createdFreightRate.DistanceMinKM,
		DistanceMaxKM: createdFreightRate.DistanceMaxKM,
		CostPerKM:     createdFreightRate.CostPerKM,
		IsDeleted:     createdFreightRate.IsDeleted,
		CreatedAt:     createdFreightRate.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     createdFreightRate.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *freightRateUsecase) DeleteFreightRate(ctx context.Context, req *model.DeleteFreightRateRequest) error {
	pu.log.Infof("Deleting freightRate with ID: %d", req.FreightRateID)
	freightRate, err := pu.freightRateRepo.Get(ctx, req.FreightRateID)
	if err != nil {
		pu.log.Errorf("Error fetching freightRate for deletion: %v", err)
		return err
	}

	freightRate.IsDeleted = true

	_, err = pu.freightRateRepo.Update(ctx, freightRate)
	if err != nil {
		pu.log.Errorf("Error updating freightRate for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted freightRate with ID: %d", req.FreightRateID)
	return nil
}

func (pu *freightRateUsecase) UpdateFreightRate(ctx context.Context, rep *model.UpdateFreightRateRequest) (*model.GetFreightRateResponse, error) {
	pu.log.Infof("Updating freightRate: %+v", rep)
	freightRate, err := pu.freightRateRepo.Get(ctx, rep.FreightRateID)
	if err != nil {
		pu.log.Errorf("Error fetching freightRate for update: %v", err)
		return nil, err
	}

	freightRate.CourierID = rep.CourierID
	freightRate.DistanceMinKM = rep.DistanceMinKM
	freightRate.DistanceMaxKM = rep.DistanceMaxKM
	freightRate.CostPerKM = rep.CostPerKM
	freightRate.IsDeleted = rep.IsDeleted
	freightRate.UpdatedAt = time.Now()

	updatedFreightRate, err := pu.freightRateRepo.Update(ctx, freightRate)
	if err != nil {
		pu.log.Errorf("Error updating freightRate: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated freightRate: %+v", updatedFreightRate)
	return &model.GetFreightRateResponse{
		FreightRateID: updatedFreightRate.FreightRateID,
		CourierID:     updatedFreightRate.CourierID,
		DistanceMinKM: updatedFreightRate.DistanceMinKM,
		DistanceMaxKM: updatedFreightRate.DistanceMaxKM,
		CostPerKM:     updatedFreightRate.CostPerKM,
		IsDeleted:     updatedFreightRate.IsDeleted,
		CreatedAt:     updatedFreightRate.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:     updatedFreightRate.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}
