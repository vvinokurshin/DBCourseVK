package usecase

import (
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	userRepo "github.com/vvinokurshin/DBCourseVK/internal/user/repository"
	"github.com/vvinokurshin/DBCourseVK/pkg"
)

type UseCaseI interface {
	CreateUser(user *models.User) ([]models.User, error)
	GetUserByNickname(nickname string) (*models.User, error)
	UpdateUser(user *models.User) error
}

type UseCase struct {
	userRepo userRepo.RepositoryI
}

func NewUseCase(userRepo userRepo.RepositoryI) UseCaseI {
	return &UseCase{
		userRepo: userRepo,
	}
}

func (uc *UseCase) CreateUser(user *models.User) ([]models.User, error) {
	existUsers, err := uc.userRepo.SelectUsersByNicknameOrEmail(user.Nickname, user.Email)
	if err != nil {
		return nil, err
	} else if len(existUsers) > 0 {
		return existUsers, pkg.ErrConflict
	}

	return nil, uc.userRepo.InsertUser(user)
}

func (uc *UseCase) GetUserByNickname(nickname string) (*models.User, error) {
	user, err := uc.userRepo.SelectUserByNickname(nickname)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UseCase) UpdateUser(user *models.User) error {
	oldUser, err := uc.userRepo.SelectUserByNickname(user.Nickname)
	if err != nil {
		return err
	}

	if user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if user.Email == "" {
		user.Email = oldUser.Email
	}
	if user.About == "" {
		user.About = oldUser.About
	}

	userByEmail, err := uc.userRepo.SelectUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, pkg.ErrNotFound) {
			return err
		}
	} else if userByEmail.Email != oldUser.Email {
		return pkg.ErrConflict
	}

	return uc.userRepo.UpdateUser(user)
}
