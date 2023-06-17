package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	threadUC "github.com/vvinokurshin/DBCourseVK/internal/thread/usecase"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"net/http"
	"strconv"
)

type DeliveryI interface {
	CreateThread(c echo.Context) error
	GetThreadsByForum(c echo.Context) error
	CreateVote(c echo.Context) error
	GetThread(c echo.Context) error
	UpdateThread(c echo.Context) error
}

type Delivery struct {
	threadUC threadUC.UseCaseI
}

func NewDelivery(threadUC threadUC.UseCaseI) DeliveryI {
	return &Delivery{
		threadUC: threadUC,
	}
}

func (delivery *Delivery) CreateThread(c echo.Context) error {
	var thread models.Thread

	err := c.Bind(&thread)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest)
	}

	thread.Forum = c.Param("slug")

	err = delivery.threadUC.CreateThread(&thread)
	if err != nil {
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		case errors.Is(err, pkg.ErrConflict):
			c.Logger().Error(err)
			return c.JSON(http.StatusConflict, thread)
		default:
			c.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusCreated, thread)
}

func (delivery *Delivery) GetThreadsByForum(c echo.Context) error {
	forumSlug := c.Param("slug")
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 100
	}

	since := c.QueryParam("since")
	reverse, err := strconv.ParseBool(c.QueryParam("desc"))

	threads, err := delivery.threadUC.GetThreadsByForum(forumSlug, limit, since, reverse)
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

	return c.JSON(http.StatusOK, threads)
}

func (delivery *Delivery) CreateVote(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	var vote models.Vote

	err := c.Bind(&vote)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest)
	}

	thread, err := delivery.threadUC.CreateVote(slugOrID, &vote)
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

	return c.JSON(http.StatusOK, thread)
}

func (delivery *Delivery) GetThread(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	thread, err := delivery.threadUC.GetThread(slugOrID)
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

	return c.JSON(http.StatusOK, thread)
}

func (delivery *Delivery) UpdateThread(c echo.Context) error {
	slugOrID := c.Param("slug_or_id")

	var thread models.Thread

	err := c.Bind(&thread)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest)
	}

	err = delivery.threadUC.UpdateThread(slugOrID, &thread)
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

	return c.JSON(http.StatusOK, thread)
}
