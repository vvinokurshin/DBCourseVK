package repository

import "github.com/vvinokurshin/DBCourseVK/internal/models"

type RepositoryI interface {
	DeleteAll() error
	SelectStatus() (*models.ServiceStatus, error)
}
