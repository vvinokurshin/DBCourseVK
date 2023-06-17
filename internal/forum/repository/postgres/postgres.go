package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/forum/repository"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
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

func (repo *Repository) InsertForum(forum *models.Forum) error {
	_, err := repo.DB.NamedExec(`INSERT INTO forums (title,user_nickname,slug) VALUES (:title,:user_nickname,:slug)`, forum)
	if err != nil {
		return errors.Wrap(err, "database error (table: forums, method: InsertForum)")
	}

	return nil
}

func (repo *Repository) SelectForum(slug string) (*models.Forum, error) {
	var forum models.Forum

	err := repo.DB.Get(&forum, "SELECT * FROM forums WHERE slug = $1", slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: forums, method: SelectForum)")
	}

	return &forum, nil
}

func (repo *Repository) InsertForumUser(slug string, nickname string) error {
	_, err := repo.DB.Exec(`INSERT INTO forum_user (user_nickname,forum) VALUES ($1,$2) ON CONFLICT DO NOTHING`, nickname, slug)
	if err != nil {
		return errors.Wrap(err, "database error (table: forum_user, method: InsertForumUser)")
	}

	return nil
}

func (repo *Repository) SelectUsersByForum(slug string, limit int, since string, reverse bool) ([]models.User, error) {
	users := make([]models.User, 0, 10)
	query := "SELECT * FROM users WHERE nickname IN (SELECT user_nickname FROM forum_user WHERE forum = $1"

	nicknameCondition := " AND user_nickname < $2 COLLATE \"C\")"
	orderCondition := " ORDER BY nickname COLLATE \"C\" DESC"

	if !reverse {
		nicknameCondition = " AND user_nickname > $2 COLLATE \"C\")"
		orderCondition = " ORDER BY nickname COLLATE \"C\""
	}

	var err error

	if since != "" {
		query += nicknameCondition + orderCondition + " LIMIT $3"
		err = repo.DB.Select(&users, query, slug, since, limit)
	} else {
		query += ")" + orderCondition + " LIMIT $2"
		err = repo.DB.Select(&users, query, slug, limit)
	}

	if err != nil {
		return nil, errors.Wrap(err, "database error (table: users, method: SelectUsersByForum)")
	}

	return users, nil
}
