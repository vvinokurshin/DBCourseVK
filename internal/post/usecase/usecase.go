package usecase

import (
	"github.com/pkg/errors"
	forumRepo "github.com/vvinokurshin/DBCourseVK/internal/forum/repository"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	postRepo "github.com/vvinokurshin/DBCourseVK/internal/post/repository"
	threadRepo "github.com/vvinokurshin/DBCourseVK/internal/thread/repository"
	userRepo "github.com/vvinokurshin/DBCourseVK/internal/user/repository"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"strconv"
	"time"
)

type UseCaseI interface {
	CreatePosts(threadSlugOrID string, posts []models.Post) error
	GetPostsByThread(threadSlugOrID string, limit, since int, reverse bool, sort string) ([]models.Post, error)
	GetPost(ID int64, related []string) (*models.PostDetails, error)
	UpdatePost(post *models.Post) error
}

type UseCase struct {
	postRepo   postRepo.RepositoryI
	threadRepo threadRepo.RepositoryI
	userRepo   userRepo.RepositoryI
	forumRepo  forumRepo.RepositoryI
}

func NewUseCase(postRepo postRepo.RepositoryI, threadRepo threadRepo.RepositoryI, userRepo userRepo.RepositoryI, forumRepo forumRepo.RepositoryI) UseCaseI {
	return &UseCase{
		postRepo:   postRepo,
		threadRepo: threadRepo,
		userRepo:   userRepo,
		forumRepo:  forumRepo,
	}
}

func (uc *UseCase) CreatePosts(threadSlugOrID string, posts []models.Post) error {
	var thread *models.Thread

	ID, err := strconv.ParseInt(threadSlugOrID, 10, 64)
	if err != nil {
		thread, err = uc.threadRepo.SelectThreadBySlug(threadSlugOrID)
	} else {
		thread, err = uc.threadRepo.SelectThreadByID(ID)
	}
	if err != nil {
		return err
	}

	if len(posts) == 0 {
		return nil
	}

	for idx := range posts {
		posts[idx].Thread = thread.ID
		posts[idx].Forum = thread.Forum

		_, err = uc.userRepo.SelectUserByNickname(posts[idx].Author)
		if err != nil {
			return err
		}

		if posts[idx].Parent != 0 {
			parentPost, err := uc.postRepo.SelectPostByID(posts[idx].Parent)
			if err != nil {
				if errors.Is(err, pkg.ErrNotFound) {
					return pkg.ErrConflict
				}

				return err
			}

			if parentPost.Thread != posts[idx].Thread {
				return pkg.ErrConflict
			}
		}
	}

	timeNow := time.Now().Format("2006-01-02T15:04:05.999999Z")
	for idx := range posts {
		posts[idx].Created = timeNow
	}

	err = uc.postRepo.InsertPosts(posts)
	if err != nil {
		return err
	}

	for _, post := range posts {
		err = uc.forumRepo.InsertForumUser(post.Forum, post.Author)
		if err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) GetPostsByThread(threadSlugOrID string, limit, since int, reverse bool, sort string) ([]models.Post, error) {
	var thread *models.Thread

	ID, err := strconv.ParseInt(threadSlugOrID, 10, 64)
	if err != nil {
		thread, err = uc.threadRepo.SelectThreadBySlug(threadSlugOrID)
	} else {
		thread, err = uc.threadRepo.SelectThreadByID(ID)
	}
	if err != nil {
		return []models.Post{}, err
	}

	return uc.postRepo.SelectPostsByThread(thread.ID, limit, since, reverse, sort)
}

func (uc *UseCase) GetPost(ID int64, related []string) (*models.PostDetails, error) {
	postDetails := models.PostDetails{}

	post, err := uc.postRepo.SelectPostByID(ID)
	if err != nil {
		return nil, err
	}

	postDetails.Post = post

	for _, elem := range related {
		switch elem {
		case "user":
			user, err := uc.userRepo.SelectUserByNickname(post.Author)
			if err != nil {
				return nil, err
			}

			postDetails.User = user
		case "thread":
			thread, err := uc.threadRepo.SelectThreadByID(post.Thread)
			if err != nil {
				return nil, err
			}

			postDetails.Thread = thread
		case "forum":
			forum, err := uc.forumRepo.SelectForum(post.Forum)
			if err != nil {
				return nil, err
			}

			postDetails.Forum = forum
		}
	}

	return &postDetails, nil
}

func (uc *UseCase) UpdatePost(post *models.Post) error {
	oldPost, err := uc.postRepo.SelectPostByID(post.ID)
	if err != nil {
		return err
	}

	if post.Message != "" && oldPost.Message != post.Message {
		oldPost.Message = post.Message
		*post = *oldPost

		return uc.postRepo.UpdatePost(post)
	}

	*post = *oldPost

	return nil
}
