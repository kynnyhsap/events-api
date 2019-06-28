package main

import (
	"events-api/handlers"
	. "events-api/sotrage"
	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"log"
	"os"
)

func main() {
	db := NewMySQLStorage(os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))

	err := db.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	r := echo.New()

	r.Pre(middleware.RemoveTrailingSlash())

	r.GET("/tags", handlers.GetTagsList(&db))
	r.GET("/calendar/events", handlers.GetEventsList(&db))
	r.GET("/calendar/events/:id", handlers.GetEvent(&db))

	parseDOU(&db)

	c := cron.New()
	err = c.AddFunc("@every 1h30m", func() {
		parseDOU(&db)
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	r.Logger.Fatal(r.Start(":8080"))

}
