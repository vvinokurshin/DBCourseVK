package repository

import "github.com/vvinokurshin/DBCourseVK/internal/models"

type RepositoryI interface {
	InsertThread(thread *models.Thread) error
	SelectThreadBySlug(slug string) (*models.Thread, error)
	SelectThreadsByForum(forumSlug string, limit int, since string, reverse bool) ([]models.Thread, error)
	SelectThreadByID(ID int64) (*models.Thread, error)
	InsertVote(vote *models.Vote) error
	UpdateThread(thread *models.Thread) error
}
