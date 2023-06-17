package usecase

import (
	"github.com/pkg/errors"
	forumRepo "github.com/vvinokurshin/DBCourseVK/internal/forum/repository"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	threadRepo "github.com/vvinokurshin/DBCourseVK/internal/thread/repository"
	userRepo "github.com/vvinokurshin/DBCourseVK/internal/user/repository"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"strconv"
)

type UseCaseI interface {
	CreateThread(thread *models.Thread) error
	GetThreadsByForum(forumSlug string, limit int, since string, reverse bool) ([]models.Thread, error)
	CreateVote(threadSlugOrID string, vote *models.Vote) (*models.Thread, error)
	GetThread(threadSlugOrID string) (*models.Thread, error)
	UpdateThread(threadSlugOrID string, thread *models.Thread) error
}

type UseCase struct {
	threadRepo threadRepo.RepositoryI
	forumRepo  forumRepo.RepositoryI
	userRepo   userRepo.RepositoryI
}

func NewUseCase(threadRepo threadRepo.RepositoryI, forumRepo forumRepo.RepositoryI, userRepo userRepo.RepositoryI) UseCaseI {
	return &UseCase{
		threadRepo: threadRepo,
		forumRepo:  forumRepo,
		userRepo:   userRepo,
	}
}

func (uc *UseCase) CreateThread(thread *models.Thread) error {
	user, err := uc.userRepo.SelectUserByNickname(thread.Author)
	if err != nil {
		return err
	}

	forum, err := uc.forumRepo.SelectForum(thread.Forum)
	if err != nil {
		return err
	}

	if thread.Slug != "" {
		oldThread, err := uc.threadRepo.SelectThreadBySlug(thread.Slug)
		if err == nil {
			thread.ID = oldThread.ID
			thread.Title = oldThread.Title
			thread.Forum = oldThread.Forum
			thread.Author = oldThread.Author
			thread.Message = oldThread.Message
			thread.Slug = oldThread.Slug
			thread.Created = oldThread.Created
			thread.Votes = oldThread.Votes
			return pkg.ErrConflict
		} else if !errors.Is(err, pkg.ErrNotFound) {
			return err
		}
	}

	thread.Forum = forum.Slug
	thread.Author = user.Nickname

	err = uc.threadRepo.InsertThread(thread)
	if err != nil {
		return err
	}

	return uc.forumRepo.InsertForumUser(thread.Forum, thread.Author)
}

func (uc *UseCase) GetThreadsByForum(forumSlug string, limit int, since string, reverse bool) ([]models.Thread, error) {
	_, err := uc.forumRepo.SelectForum(forumSlug)
	if err != nil {
		return nil, err
	}

	return uc.threadRepo.SelectThreadsByForum(forumSlug, limit, since, reverse)
}

func (uc *UseCase) CreateVote(threadSlugOrID string, vote *models.Vote) (*models.Thread, error) {
	var thread *models.Thread

	ID, err := strconv.ParseInt(threadSlugOrID, 10, 64)
	if err != nil {
		thread, err = uc.threadRepo.SelectThreadBySlug(threadSlugOrID)
	} else {
		thread, err = uc.threadRepo.SelectThreadByID(ID)
	}
	if err != nil {
		return nil, err
	}

	vote.Thread = thread.ID

	_, err = uc.userRepo.SelectUserByNickname(vote.Nickname)
	if err != nil {
		return nil, err
	}

	err = uc.threadRepo.InsertVote(vote)
	if err != nil {
		return nil, err
	}

	return uc.threadRepo.SelectThreadByID(thread.ID)
}

func (uc *UseCase) GetThread(threadSlugOrID string) (*models.Thread, error) {
	ID, err := strconv.ParseInt(threadSlugOrID, 10, 64)
	if err != nil {
		return uc.threadRepo.SelectThreadBySlug(threadSlugOrID)
	}

	return uc.threadRepo.SelectThreadByID(ID)
}

func (uc *UseCase) UpdateThread(threadSlugOrID string, thread *models.Thread) error {
	var oldThread *models.Thread

	ID, err := strconv.ParseInt(threadSlugOrID, 10, 64)
	if err != nil {
		oldThread, err = uc.threadRepo.SelectThreadBySlug(threadSlugOrID)
	} else {
		oldThread, err = uc.threadRepo.SelectThreadByID(ID)
	}
	if err != nil {
		return err
	}

	if thread.Title != "" {
		oldThread.Title = thread.Title
	}
	if thread.Message != "" {
		oldThread.Message = thread.Message
	}

	*thread = *oldThread

	return uc.threadRepo.UpdateThread(thread)
}
