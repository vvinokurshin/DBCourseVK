package repository

import "github.com/vvinokurshin/DBCourseVK/internal/models"

type RepositoryI interface {
	InsertPosts(posts []models.Post) error
	SelectPostByID(ID int64) (*models.Post, error)
	SelectPostsByThread(threadID int64, limit, since int, reverse bool, sort string) ([]models.Post, error)
	UpdatePost(post *models.Post) error
}
