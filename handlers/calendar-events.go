package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tobira-shoe/dou-events-parser"
	"net/http"
	"strconv"
)

func GetEvent(c echo.Context) error {
	err, list := parser.ParseCalendarEvents()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	id, _ := strconv.Atoi(c.Param("id"))
	for _, event := range list {
		if event.ID == id {
			return c.JSON(http.StatusOK, event)
		}
	}

	return c.String(http.StatusNotFound, "Event not found")
}

func GetEventsList(c echo.Context) error {
	err, list := parser.ParseCalendarEvents()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	return c.JSON(http.StatusOK, list[offset:offset+limit])
}
