package postgres

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	"github.com/vvinokurshin/DBCourseVK/internal/service/repository"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.RepositoryI {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) DeleteAll() error {
	_, err := repo.DB.Exec(`TRUNCATE posts,threads,forums,users,forum_user CASCADE`)
	if err != nil {
		return errors.Wrap(err, "database error (table: all, method: DeleteAll)")
	}

	return nil
}

func (repo *Repository) selectTableStatus(tableName string) (int64, error) {
	var count int64

	err := repo.DB.Get(&count, `SELECT COUNT(*) FROM `+tableName)
	if err != nil {
		return 0, errors.Wrap(err, "database error (table: "+tableName+", method: selectTableStatus)")
	}

	return count, nil
}

func (repo *Repository) SelectStatus() (*models.ServiceStatus, error) {
	var status models.ServiceStatus
	var err error

	status.UserCount, err = repo.selectTableStatus("users")
	if err != nil {
		return &models.ServiceStatus{}, err
	}

	status.ForumCount, err = repo.selectTableStatus("forums")
	if err != nil {
		return &models.ServiceStatus{}, err
	}

	status.ThreadCount, err = repo.selectTableStatus("threads")
	if err != nil {
		return &models.ServiceStatus{}, err
	}

	status.PostCount, err = repo.selectTableStatus("posts")
	if err != nil {
		return &models.ServiceStatus{}, err
	}

	return &status, nil
}
