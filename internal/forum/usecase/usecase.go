package usecase

import (
	"errors"
	forumRepo "github.com/vvinokurshin/DBCourseVK/internal/forum/repository"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	userRepo "github.com/vvinokurshin/DBCourseVK/internal/user/repository"
	"github.com/vvinokurshin/DBCourseVK/pkg"
)

type UseCaseI interface {
	CreateForum(forum *models.Forum) error
	GetForum(slug string) (*models.Forum, error)
	GetUsersByForum(slug string, limit int, since string, reverse bool) ([]models.User, error)
}

type UseCase struct {
	forumRepo forumRepo.RepositoryI
	userRepo  userRepo.RepositoryI
}

func NewUseCase(forumRepo forumRepo.RepositoryI, userRepo userRepo.RepositoryI) UseCaseI {
	return &UseCase{
		forumRepo: forumRepo,
		userRepo:  userRepo,
	}
}

func (uc *UseCase) CreateForum(forum *models.Forum) error {
	user, err := uc.userRepo.SelectUserByNickname(forum.User)
	if err != nil {
		return err
	}

	existForum, err := uc.forumRepo.SelectForum(forum.Slug)
	if err == nil {
		forum.User = existForum.User
		forum.Slug = existForum.Slug
		forum.Title = existForum.Title
		forum.Threads = existForum.Threads
		forum.Posts = existForum.Posts
		return pkg.ErrConflict
	} else if !errors.Is(err, pkg.ErrNotFound) {
		return err
	}

	forum.User = user.Nickname

	return uc.forumRepo.InsertForum(forum)
}

func (uc *UseCase) GetForum(slug string) (*models.Forum, error) {
	forum, err := uc.forumRepo.SelectForum(slug)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (uc *UseCase) GetUsersByForum(slug string, limit int, since string, reverse bool) ([]models.User, error) {
	_, err := uc.forumRepo.SelectForum(slug)
	if err != nil {
		return nil, err
	}

	return uc.forumRepo.SelectUsersByForum(slug, limit, since, reverse)
}
