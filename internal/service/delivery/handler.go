package delivery

import (
	"github.com/labstack/echo/v4"
	serviceUC "github.com/vvinokurshin/DBCourseVK/internal/service/usecase"
	"net/http"
)

type DeliveryI interface {
	ClearAll(c echo.Context) error
	GetStatus(c echo.Context) error
}

type Delivery struct {
	ServUC serviceUC.UseCaseI
}

func NewDelivery(servUC serviceUC.UseCaseI) DeliveryI {
	return &Delivery{
		ServUC: servUC,
	}
}

func (delivery *Delivery) ClearAll(c echo.Context) error {
	err := delivery.ServUC.ClearAll()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

func (delivery *Delivery) GetStatus(c echo.Context) error {
	status, err := delivery.ServUC.GetStatus()
	if err != nil {
		c.Logger().Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, status)
}
