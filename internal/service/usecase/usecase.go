package usecase

import (
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	serviceRepo "github.com/vvinokurshin/DBCourseVK/internal/service/repository"
)

type UseCaseI interface {
	ClearAll() error
	GetStatus() (*models.ServiceStatus, error)
}

type UseCase struct {
	servRepo serviceRepo.RepositoryI
}

func NewUseCase(servRepo serviceRepo.RepositoryI) UseCaseI {
	return &UseCase{
		servRepo: servRepo,
	}
}

func (uc *UseCase) ClearAll() error {
	return uc.servRepo.DeleteAll()
}

func (uc *UseCase) GetStatus() (*models.ServiceStatus, error) {
	return uc.servRepo.SelectStatus()
}
