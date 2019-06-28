package main

import (
	storage "events-api/sotrage"
	"fmt"
	"github.com/briandowns/spinner"
	douparser "github.com/tobira-shoe/dou-events-parser"
	"time"
)

func parseDOU(db storage.Storage) {
	s := spinner.New(spinner.CharSets[7], 100*time.Millisecond)
	s.Start()

	fmt.Println("Start scheduled parsing...")

	err, events := douparser.ParseCalendarEvents()
	if err != nil {
		fmt.Println(err)
	} else {
		err = db.SaveEventsList(events)
		if err != nil {
			fmt.Println(err)
		}
	}

	err, tags := douparser.ParseEventTags()
	if err != nil {
		fmt.Println(err)
	} else {
		err = db.SaveTagsList(tags)
		if err != nil {
			fmt.Println(err)
		}
	}

	// todo: log parsing time
	s.Stop()
}
