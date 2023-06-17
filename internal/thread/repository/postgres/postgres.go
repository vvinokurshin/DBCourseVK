package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	"github.com/vvinokurshin/DBCourseVK/internal/thread/repository"
	"github.com/vvinokurshin/DBCourseVK/pkg"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) repository.RepositoryI {
	return &Repository{
		DB: db,
	}
}

func (repo *Repository) InsertThread(thread *models.Thread) error {
	var rows *sqlx.Rows
	var err error

	if thread.Slug == "" {
		rows, err = repo.DB.NamedQuery(`INSERT INTO threads (title,author,forum,message,created) VALUES (:title,:author,:forum,:message,:created) RETURNING id`,
			thread)
	} else {
		rows, err = repo.DB.NamedQuery(`INSERT INTO threads (title,author,forum,message,slug,created) VALUES (:title,:author,:forum,:message,:slug,:created) RETURNING id`,
			thread)
	}
	defer rows.Close()

	if err != nil {
		return errors.Wrap(err, "database error (table: forums, method: InsertForum)")
	}

	rows.Next()
	err = rows.Scan(&thread.ID)
	if err != nil {
		return errors.Wrap(err, "database error (table: forums, method: InsertForum)")
	}

	return nil
}

func (repo *Repository) SelectThreadBySlug(slug string) (*models.Thread, error) {
	var thread models.Thread

	err := repo.DB.Get(&thread, "SELECT id, title, author, forum, message, votes, COALESCE(slug, '') AS slug, created FROM threads WHERE slug = $1", slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: threads, method: SelectThreadBySlug)")
	}

	return &thread, nil
}

func (repo *Repository) SelectThreadsByForum(forumSlug string, limit int, since string, reverse bool) ([]models.Thread, error) {
	threads := make([]models.Thread, 0, 10)

	query := "SELECT id, title, author, forum, message, votes, COALESCE(slug, '') AS slug, created FROM threads WHERE forum = $1"

	createdCondition := " AND created <= $2"
	orderCondition := " ORDER BY created DESC"

	if !reverse {
		createdCondition = " AND created >= $2"
		orderCondition = " ORDER BY created"
	}

	var err error

	if since != "" {
		query += createdCondition + orderCondition + " LIMIT $3"
		err = repo.DB.Select(&threads, query, forumSlug, since, limit)
	} else {
		query += orderCondition + " LIMIT $2"
		err = repo.DB.Select(&threads, query, forumSlug, limit)
	}

	if err != nil {
		return nil, errors.Wrap(err, "database error (table: users, method: SelectThreadsByForum)")
	}

	return threads, nil
}

func (repo *Repository) SelectThreadByID(ID int64) (*models.Thread, error) {
	var thread models.Thread

	err := repo.DB.Get(&thread, "SELECT id, title, author, forum, message, votes, COALESCE(slug, '') AS slug, created FROM threads WHERE id = $1", ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: threads, method: SelectThreadByID)")
	}

	return &thread, nil
}

func (repo *Repository) InsertVote(vote *models.Vote) error {
	_, err := repo.DB.NamedExec(`INSERT INTO votes (thread,nickname,voice) VALUES (:thread,:nickname,:voice)
ON CONFLICT (thread,nickname) DO UPDATE SET voice=:voice`, vote)
	if err != nil {
		return errors.Wrap(err, "database error (table: votes, method: InsertVote)")
	}

	return nil
}

func (repo *Repository) UpdateThread(thread *models.Thread) error {
	_, err := repo.DB.NamedExec(`UPDATE threads SET title=:title,message=:message WHERE id = :id`, thread)
	if err != nil {
		return errors.Wrap(err, "database error (table: threads, method: UpdateThread)")
	}

	return nil
}
