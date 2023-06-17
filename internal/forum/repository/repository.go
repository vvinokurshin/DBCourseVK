package repository

import "github.com/vvinokurshin/DBCourseVK/internal/models"

type RepositoryI interface {
	InsertForum(forum *models.Forum) error
	SelectForum(slug string) (*models.Forum, error)
	InsertForumUser(slug string, nickname string) error
	SelectUsersByForum(slug string, limit int, since string, reverse bool) ([]models.User, error)
}
