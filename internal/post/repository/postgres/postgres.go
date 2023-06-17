package postgres

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	"github.com/vvinokurshin/DBCourseVK/internal/post/repository"
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

func (repo *Repository) InsertPosts(posts []models.Post) error {
	rows, err := repo.DB.NamedQuery(`INSERT INTO posts (parent,author,message,is_edited,forum,thread,created) 
VALUES (:parent,:author,:message,:is_edited,:forum,:thread,:created) RETURNING id`, posts)
	defer rows.Close()

	if err != nil {
		return errors.Wrap(err, "database error (table: posts, method: InsertPosts)")
	}

	for idx := range posts {
		rows.Next()
		err = rows.Scan(&posts[idx].ID)
		if err != nil {
			return errors.Wrap(err, "database error (table: posts, method: InsertPosts)")

		}
	}

	return nil
}

func (repo *Repository) SelectPostByID(ID int64) (*models.Post, error) {
	var post models.Post

	err := repo.DB.Get(&post, "SELECT id,parent,author,message,is_edited,forum,thread,created FROM posts WHERE id = $1", ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, pkg.ErrNotFound
		}

		return nil, errors.Wrap(err, "database error (table: posts, method: SelectPostByID)")
	}

	return &post, nil
}

func (repo *Repository) selectPostsFlatSort(threadID int64, limit, since int, reverse bool) ([]models.Post, error) {
	posts := make([]models.Post, 0, 10)
	var err error

	query := `SELECT id,parent,author,message,is_edited,forum,thread,created FROM posts WHERE thread = $1`

	if reverse {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND id < $2 ORDER BY id DESC LIMIT $3`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY id DESC LIMIT $2`, threadID, limit)
		}
	} else {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND id > $2 ORDER BY id LIMIT $3`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY id LIMIT $2`, threadID, limit)
		}
	}
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "database error (table: posts, method: selectPostsFlatSort)")
	}

	return posts, nil
}

func (repo *Repository) selectPostsTreeSort(threadID int64, limit, since int, reverse bool) ([]models.Post, error) {
	posts := make([]models.Post, 0, 10)
	var err error

	query := `SELECT id,parent,author,message,is_edited,forum,thread,created FROM posts WHERE thread = $1`

	if reverse {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND post_tree < (SELECT post_tree FROM posts WHERE ID = $2) 
			ORDER BY post_tree DESC LIMIT $3`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY post_tree DESC LIMIT $2`, threadID, limit)
		}
	} else {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND post_tree > (SELECT post_tree FROM posts WHERE ID = $2) 
			ORDER BY post_tree LIMIT $3`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY post_tree LIMIT $2`, threadID, limit)
		}
	}
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "database error (table: posts, method: selectPostsTreeSort)")
	}

	return posts, nil
}

func (repo *Repository) selectPostsParentTreeSort(threadID int64, limit, since int, reverse bool) ([]models.Post, error) {
	posts := make([]models.Post, 0, 10)
	var err error

	query := `SELECT id,parent,author,message,is_edited,forum,thread,created FROM posts WHERE post_tree[1] IN 
	(SELECT id FROM posts WHERE parent = 0 AND thread = $1`

	if reverse {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND id < (SELECT post_tree[1] FROM posts WHERE id = $2) ORDER BY id DESC LIMIT $3) 
			ORDER BY post_tree[1] desc, post_tree`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY id DESC LIMIT $2) ORDER BY post_tree[1] DESC, post_tree`, threadID, limit)
		}
	} else {
		if since != 0 {
			err = repo.DB.Select(&posts, query+` AND id > (SELECT post_tree[1] FROM posts WHERE id = $2) ORDER BY id LIMIT $3) 
			ORDER BY post_tree`, threadID, since, limit)
		} else {
			err = repo.DB.Select(&posts, query+` ORDER BY id LIMIT $2) ORDER BY post_tree`, threadID, limit)
		}
	}
	if err != nil {
		return []models.Post{}, errors.Wrap(err, "database error (table: posts, method: selectPostsParentTreeSort)")
	}

	return posts, nil
}

func (repo *Repository) SelectPostsByThread(threadID int64, limit, since int, reverse bool, sort string) ([]models.Post, error) {
	switch sort {
	case "flat":
		return repo.selectPostsFlatSort(threadID, limit, since, reverse)
	case "tree":
		return repo.selectPostsTreeSort(threadID, limit, since, reverse)
	case "parent_tree":
		return repo.selectPostsParentTreeSort(threadID, limit, since, reverse)
	}

	return []models.Post{}, nil
}

func (repo *Repository) UpdatePost(post *models.Post) error {
	post.IsEdited = true

	_, err := repo.DB.NamedExec(`UPDATE posts SET message=:message,is_edited=:is_edited WHERE id = :id`, post)
	if err != nil {
		return errors.Wrap(err, "database error (table: posts, method: UpdatePost)")
	}

	return nil
}
