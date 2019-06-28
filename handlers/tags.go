package handlers

import (
	"events-api/sotrage"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetTagsList(db storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		tags, err := db.GetTagsList()

		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(http.StatusOK, tags)
	}
}
