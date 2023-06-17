package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	"github.com/vvinokurshin/DBCourseVK/internal/user/repository"
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

func (repo *Repository) InsertUser(user *models.User) error {
	_, err := repo.DB.NamedExec(`INSERT INTO users (nickname,fullname,about,email) VALUES (:nickname,:fullname,:about,:email)`, user)
	if err != nil {
		return errors.Wrap(err, "database error (table: users, method: InsertUser)")
	}

	return nil
}

func (repo *Repository) SelectUsersByNicknameOrEmail(nickname, email string) ([]models.User, error) {
	users := make([]models.User, 0, 10)

	err := repo.DB.Select(&users, "SELECT * FROM users WHERE nickname = $1 OR email = $2", nickname, email)
	if err != nil {
		return nil, errors.Wrap(err, "database error (table: users, method: SelectUsersByNicknameOrEmail)")
	}

	return users, nil
}

func (repo *Repository) SelectUserByNickname(nickname string) (*models.User, error) {
	var user models.User

	err := repo.DB.Get(&user, "SELECT * FROM users WHERE nickname = $1", nickname)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: users, method: SelectUserByNickname)")
	}

	return &user, nil
}

func (repo *Repository) SelectUserByEmail(email string) (*models.User, error) {
	var user models.User

	err := repo.DB.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: users, method: SelectUserByEmail)")
	}

	return &user, nil
}

func (repo *Repository) UpdateUser(user *models.User) error {
	_, err := repo.DB.NamedExec(`UPDATE users SET fullname=:fullname,about=:about,email=:email WHERE nickname = :nickname`, user)
	if err != nil {
		return errors.Wrap(err, "database error (table: users, method: UpdateUser)")
	}

	return nil
}
