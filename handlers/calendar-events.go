package handlers

import (
	"events-api/sotrage"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func GetEvent(db storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, _ := strconv.Atoi(c.Param("id"))
		event, err := db.GetEvent(id)
		if err != nil {
			return c.String(http.StatusNotFound, "Event not found?")
		}

		return c.JSON(http.StatusOK, event)
	}
}

func GetEventsList(db storage.Storage) echo.HandlerFunc {
	return func(c echo.Context) error {
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		offset, _ := strconv.Atoi(c.QueryParam("offset"))

		// todo: tags

		list, err := db.GetEventsList(limit, offset, []string{})
		if err != nil {
			return c.String(http.StatusInternalServerError, "")
		}

		return c.JSON(http.StatusOK, list)
	}
}
