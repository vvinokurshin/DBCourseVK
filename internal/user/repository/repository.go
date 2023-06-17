package repository

import (
	"github.com/vvinokurshin/DBCourseVK/internal/models"
)

type RepositoryI interface {
	InsertUser(user *models.User) error
	SelectUsersByNicknameOrEmail(nickname, email string) ([]models.User, error)
	SelectUserByNickname(nickname string) (*models.User, error)
	SelectUserByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
}
