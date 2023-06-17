package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	postUC "github.com/vvinokurshin/DBCourseVK/internal/post/usecase"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"net/http"
	"strconv"
	"strings"
)

type DeliveryI interface {
	CreatePosts(c echo.Context) error
	GetPostsByThread(c echo.Context) error
	GetPost(c echo.Context) error
	UpdatePost(c echo.Context) error
}

type Delivery struct {
	postUC postUC.UseCaseI
}

func NewDelivery(postUC postUC.UseCaseI) DeliveryI {
	return &Delivery{
		postUC: postUC,
	}
}

func (delivery *Delivery) CreatePosts(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	posts := make([]models.Post, 0, 10)

	err := c.Bind(&posts)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	err = delivery.postUC.CreatePosts(slugOrID, posts)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		case errors.Is(err, pkg.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, pkg.ErrConflict.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, posts)
}

func (delivery *Delivery) GetPostsByThread(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	sinceStr := c.QueryParam("since")
	since, err := strconv.Atoi(sinceStr)
	if sinceStr != "" && err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	reverse, err := strconv.ParseBool(c.QueryParam("desc"))
	sort := c.QueryParam("sort")
	if sort != "flat" && sort != "tree" && sort != "parent_tree" {
		sort = "flat"
	}

	posts, err := delivery.postUC.GetPostsByThread(slugOrID, limit, since, reverse, sort)
	if err != nil {
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, posts)
}

func (delivery *Delivery) GetPost(c echo.Context) error {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	queryRelated := c.QueryParam("related")
	var related []string

	if queryRelated != "" {
		related = strings.Split(queryRelated, ",")
		for _, elem := range related {
			if elem != "user" && elem != "forum" && elem != "thread" {
				c.Logger().Error(pkg.ErrBadRequest)
				return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
			}
		}
	}

	postDetails, err := delivery.postUC.GetPost(ID, related)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, postDetails)
}

func (delivery *Delivery) UpdatePost(c echo.Context) error {
	ID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	var post models.Post

	err = c.Bind(&post)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	post.ID = ID
	err = delivery.postUC.UpdatePost(&post)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, post)
}
