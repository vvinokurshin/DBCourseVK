package delivery

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/vvinokurshin/DBCourseVK/internal/models"
	userUC "github.com/vvinokurshin/DBCourseVK/internal/user/usecase"
	"github.com/vvinokurshin/DBCourseVK/pkg"
	"net/http"
)

type DeliveryI interface {
	CreateUser(c echo.Context) error
	GetUser(c echo.Context) error
	UpdateUser(c echo.Context) error
}

type Delivery struct {
	UserUC userUC.UseCaseI
}

func NewDelivery(uUC userUC.UseCaseI) DeliveryI {
	return &Delivery{
		UserUC: uUC,
	}
}

func (delivery *Delivery) CreateUser(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest)
	}

	user.Nickname = c.Param("nickname")
	conflictUsers, err := delivery.UserUC.CreateUser(&user)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrConflict):
			return c.JSON(http.StatusConflict, conflictUsers)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, pkg.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusCreated, user)
}

func (delivery *Delivery) GetUser(c echo.Context) error {
	nickname := c.Param("nickname")
	user, err := delivery.UserUC.GetUserByNickname(nickname)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, pkg.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, user)
}

func (delivery *Delivery) UpdateUser(c echo.Context) error {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, pkg.ErrBadRequest.Error())
	}

	user.Nickname = c.Param("nickname")
	err = delivery.UserUC.UpdateUser(&user)

	if err != nil {
		c.Logger().Error(err)
		switch {
		case errors.Is(err, pkg.ErrConflict):
			return echo.NewHTTPError(http.StatusConflict, pkg.ErrConflict.Error())
		case errors.Is(err, pkg.ErrNotFound):
			return echo.NewHTTPError(http.StatusNotFound, pkg.ErrNotFound.Error())
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, pkg.ErrInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, user)
}
