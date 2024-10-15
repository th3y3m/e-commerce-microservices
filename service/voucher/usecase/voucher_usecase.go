package usecase

import (
	"context"
	"th3y3m/e-commerce-microservices/pkg/util"
	"th3y3m/e-commerce-microservices/service/voucher/model"
	"th3y3m/e-commerce-microservices/service/voucher/repository"
	"time"

	"github.com/sirupsen/logrus"
)

const tsCreateTimeLayout = "2006-01-02 15:04:05 +0700"

type voucherUsecase struct {
	log         *logrus.Logger
	voucherRepo repository.IVoucherRepository
}

type IVoucherUsecase interface {
	GetVoucher(ctx context.Context, req *model.GetVoucherRequest) (*model.GetVoucherResponse, error)
	GetAllVouchers(ctx context.Context) ([]*model.GetVoucherResponse, error)
	CreateVoucher(ctx context.Context, req *model.CreateVoucherRequest) (*model.GetVoucherResponse, error)
	UpdateVoucher(ctx context.Context, rep *model.UpdateVoucherRequest) (*model.GetVoucherResponse, error)
	DeleteVoucher(ctx context.Context, req *model.DeleteVoucherRequest) error
	GetVoucherList(ctx context.Context, req *model.GetVouchersRequest) (*util.PaginatedList[model.GetVoucherResponse], error)
}

func NewVoucherUsecase(voucherRepo repository.IVoucherRepository, log *logrus.Logger) IVoucherUsecase {
	return &voucherUsecase{
		voucherRepo: voucherRepo,
		log:         log,
	}
}

func (pu *voucherUsecase) GetVoucher(ctx context.Context, req *model.GetVoucherRequest) (*model.GetVoucherResponse, error) {
	pu.log.Infof("Fetching voucher with ID: %d", req.VoucherID)
	voucher, err := pu.voucherRepo.Get(ctx, req.VoucherID)
	if err != nil {
		pu.log.Errorf("Error fetching voucher: %v", err)
		return nil, err
	}

	pu.log.Infof("Fetched voucher: %+v", voucher)
	return &model.GetVoucherResponse{
		VoucherID: voucher.VoucherID,
	}, nil
}

func (pu *voucherUsecase) GetAllVouchers(ctx context.Context) ([]*model.GetVoucherResponse, error) {
	pu.log.Info("Fetching all vouchers")
	vouchers, err := pu.voucherRepo.GetAll(ctx)
	if err != nil {
		pu.log.Errorf("Error fetching all vouchers: %v", err)
		return nil, err
	}

	var voucherResponses []*model.GetVoucherResponse
	for _, voucher := range vouchers {
		voucherResponses = append(voucherResponses, &model.GetVoucherResponse{
			VoucherID:          voucher.VoucherID,
			VoucherCode:        voucher.VoucherCode,
			DiscountType:       voucher.DiscountType,
			DiscountValue:      voucher.DiscountValue,
			MinimumOrderAmount: voucher.MinimumOrderAmount,
			MaxDiscountAmount:  voucher.MaxDiscountAmount,
			StartDate:          voucher.StartDate.Format(tsCreateTimeLayout),
			EndDate:            voucher.EndDate.Format(tsCreateTimeLayout),
			UsageLimit:         voucher.UsageLimit,
			UsageCount:         voucher.UsageCount,
			IsDeleted:          voucher.IsDeleted,
			CreatedAt:          voucher.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:          voucher.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	pu.log.Infof("Fetched %d vouchers", len(voucherResponses))
	return voucherResponses, nil
}

func (pu *voucherUsecase) CreateVoucher(ctx context.Context, voucher *model.CreateVoucherRequest) (*model.GetVoucherResponse, error) {
	pu.log.Infof("Creating voucher: %+v", voucher)
	voucherData := repository.Voucher{
		VoucherCode:        voucher.VoucherCode,
		DiscountType:       voucher.DiscountType,
		DiscountValue:      voucher.DiscountValue,
		MinimumOrderAmount: voucher.MinimumOrderAmount,
		MaxDiscountAmount:  voucher.MaxDiscountAmount,
		StartDate:          voucher.StartDate,
		EndDate:            voucher.EndDate,
		UsageLimit:         voucher.UsageLimit,
		UsageCount:         voucher.UsageCount,
	}

	createdVoucher, err := pu.voucherRepo.Create(ctx, &voucherData)
	if err != nil {
		pu.log.Errorf("Error creating voucher: %v", err)
		return nil, err
	}

	pu.log.Infof("Created voucher: %+v", createdVoucher)
	return &model.GetVoucherResponse{
		VoucherID:          createdVoucher.VoucherID,
		VoucherCode:        createdVoucher.VoucherCode,
		DiscountType:       createdVoucher.DiscountType,
		DiscountValue:      createdVoucher.DiscountValue,
		MinimumOrderAmount: createdVoucher.MinimumOrderAmount,
		MaxDiscountAmount:  createdVoucher.MaxDiscountAmount,
		StartDate:          createdVoucher.StartDate.Format(tsCreateTimeLayout),
		EndDate:            createdVoucher.EndDate.Format(tsCreateTimeLayout),
		UsageLimit:         createdVoucher.UsageLimit,
		UsageCount:         createdVoucher.UsageCount,
		IsDeleted:          createdVoucher.IsDeleted,
		CreatedAt:          createdVoucher.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:          createdVoucher.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *voucherUsecase) DeleteVoucher(ctx context.Context, req *model.DeleteVoucherRequest) error {
	pu.log.Infof("Deleting voucher with ID: %d", req.VoucherID)
	voucher, err := pu.voucherRepo.Get(ctx, req.VoucherID)
	if err != nil {
		pu.log.Errorf("Error fetching voucher for deletion: %v", err)
		return err
	}

	voucher.IsDeleted = true

	_, err = pu.voucherRepo.Update(ctx, voucher)
	if err != nil {
		pu.log.Errorf("Error updating voucher for deletion: %v", err)
		return err
	}

	pu.log.Infof("Deleted voucher with ID: %d", req.VoucherID)
	return nil
}

func (pu *voucherUsecase) UpdateVoucher(ctx context.Context, rep *model.UpdateVoucherRequest) (*model.GetVoucherResponse, error) {
	pu.log.Infof("Updating voucher with ID: %d", rep.VoucherID)
	voucher, err := pu.voucherRepo.Get(ctx, rep.VoucherID)
	if err != nil {
		pu.log.Errorf("Error fetching voucher for update: %v", err)
		return nil, err
	}

	voucher.VoucherCode = rep.VoucherCode
	voucher.DiscountType = rep.DiscountType
	voucher.DiscountValue = rep.DiscountValue
	voucher.MinimumOrderAmount = rep.MinimumOrderAmount
	voucher.MaxDiscountAmount = rep.MaxDiscountAmount
	voucher.StartDate = rep.StartDate
	voucher.EndDate = rep.EndDate
	voucher.UsageLimit = rep.UsageLimit
	voucher.UsageCount = rep.UsageCount
	voucher.UpdatedAt = time.Now()
	voucher.IsDeleted = rep.IsDeleted

	updatedVoucher, err := pu.voucherRepo.Update(ctx, voucher)
	if err != nil {
		pu.log.Errorf("Error updating voucher: %v", err)
		return nil, err
	}

	pu.log.Infof("Updated voucher: %+v", updatedVoucher)
	return &model.GetVoucherResponse{
		VoucherID:          updatedVoucher.VoucherID,
		VoucherCode:        updatedVoucher.VoucherCode,
		DiscountType:       updatedVoucher.DiscountType,
		DiscountValue:      updatedVoucher.DiscountValue,
		MinimumOrderAmount: updatedVoucher.MinimumOrderAmount,
		MaxDiscountAmount:  updatedVoucher.MaxDiscountAmount,
		StartDate:          updatedVoucher.StartDate.Format(tsCreateTimeLayout),
		EndDate:            updatedVoucher.EndDate.Format(tsCreateTimeLayout),
		UsageLimit:         updatedVoucher.UsageLimit,
		UsageCount:         updatedVoucher.UsageCount,
		IsDeleted:          updatedVoucher.IsDeleted,
		CreatedAt:          updatedVoucher.CreatedAt.Format(tsCreateTimeLayout),
		UpdatedAt:          updatedVoucher.UpdatedAt.Format(tsCreateTimeLayout),
	}, nil
}

func (pu *voucherUsecase) GetVoucherList(ctx context.Context, req *model.GetVouchersRequest) (*util.PaginatedList[model.GetVoucherResponse], error) {
	pu.log.Infof("Fetching voucher list with request: %+v", req)
	vouchers, err := pu.voucherRepo.GetList(ctx, req)
	if err != nil {
		pu.log.Errorf("Error fetching voucher list: %v", err)
		return nil, err
	}

	var voucherResponses []model.GetVoucherResponse
	for _, voucher := range vouchers {
		voucherResponses = append(voucherResponses, model.GetVoucherResponse{
			VoucherID:          voucher.VoucherID,
			VoucherCode:        voucher.VoucherCode,
			DiscountType:       voucher.DiscountType,
			DiscountValue:      voucher.DiscountValue,
			MinimumOrderAmount: voucher.MinimumOrderAmount,
			MaxDiscountAmount:  voucher.MaxDiscountAmount,
			StartDate:          voucher.StartDate.Format(tsCreateTimeLayout),
			EndDate:            voucher.EndDate.Format(tsCreateTimeLayout),
			UsageLimit:         voucher.UsageLimit,
			UsageCount:         voucher.UsageCount,
			IsDeleted:          voucher.IsDeleted,
			CreatedAt:          voucher.CreatedAt.Format(tsCreateTimeLayout),
			UpdatedAt:          voucher.UpdatedAt.Format(tsCreateTimeLayout),
		})
	}

	list := &util.PaginatedList[model.GetVoucherResponse]{
		Items:      voucherResponses,
		TotalCount: len(voucherResponses),
		PageIndex:  req.Paging.PageIndex,
		PageSize:   req.Paging.PageSize,
		TotalPages: 1,
	}

	list.GetTotalPages()

	pu.log.Infof("Fetched %d vouchers", len(voucherResponses))
	return list, nil
}
