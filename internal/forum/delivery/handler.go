package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	forumUC "github.com/vvinokurshin/DBCourseVK/internal/forum/usecase"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"net/http"
	"strconv"
)

type DeliveryI interface {
	CreateForum(c echo.Context) error
	GetForum(c echo.Context) error
	GetUsersByForum(c echo.Context) error
}

type Delivery struct {
	ForumUC forumUC.UseCaseI
}

func NewDelivery(forumUC forumUC.UseCaseI) DeliveryI {
	return &Delivery{
		ForumUC: forumUC,
	}
}

func (delivery *Delivery) CreateForum(c echo.Context) error {
	var forum models.Forum

	err := c.Bind(&forum)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest)
	}

	err = delivery.ForumUC.CreateForum(&forum)
	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrConflict):
			return c.JSON(http.StatusConflict, forum)
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, pkg.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, forum)
}

func (delivery *Delivery) GetForum(c echo.Context) error {
	slug := c.Param("slug")
	forum, err := delivery.ForumUC.GetForum(slug)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, pkg.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, forum)
}

func (delivery *Delivery) GetUsersByForum(c echo.Context) error {
	forumSlug := c.Param("slug")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")
	reverse, err := strconv.ParseBool(c.QueryParam("desc"))

	users, err := delivery.ForumUC.GetUsersByForum(forumSlug, limit, since, reverse)
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

	return c.JSON(http.StatusOK, users)
}
