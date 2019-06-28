package storage

import (
	"database/sql"
	. "github.com/tobira-shoe/event-models"
	"log"
)

type MySQLStorage struct {
	user string
	pass string
	name string

	db *sql.DB
}

func NewMySQLStorage(user, pass, name string) MySQLStorage {
	s := MySQLStorage{}

	s.name = name
	s.user = user
	s.pass = pass

	return s
}

func (s *MySQLStorage) Open() error {
	dataSourceName := s.user + ":" + s.pass + "@/" + s.name + "?charset=utf8&parseTime=true"

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *MySQLStorage) Close() error {
	err := s.db.Close()

	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) SaveEventsList(events []DouEvent) error {
	// todo: truncate table before insert
	stmt, err := s.db.Prepare(`
		INSERT INTO look_events_db.events
		    (
				id,
				title,
				short_description,
				full_description,
				location,
				raw_date,
		     	start_date,
		     	end_date,
				online,
				cost,
				image
		    )

		    VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)

	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	for _, event := range events {
		_, err := stmt.Exec(event.ID,
			event.Title,
			event.ShortDescription,
			event.FullDescription,
			event.Location,
			event.RawDate,
			event.Start,
			event.End,
			event.Online,
			event.Cost,
			event.Image)

		if err != nil {
			return err
		}
	}

	// todo: insert many2many tags-events

	return nil
}

func (s *MySQLStorage) SaveTagsList(tags []string) error {
	// todo: truncate table before insert
	stmt, err := s.db.Prepare(`
		INSERT INTO look_events_db.tags (value) VALUES(?)
	`)

	if err != nil {
		log.Println(err)
		return err
	}

	defer stmt.Close()

	for _, tag := range tags {
		_, err := stmt.Exec(tag)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MySQLStorage) GetEventsList(limit, offset int, tags []string) ([]DouEvent, error) {
	events := make([]DouEvent, 0)

	rows, err := s.db.Query(`
		SELECT
			id,
			title, 
			short_description, 
			full_description, 
			location, 
			online, 
			cost, 
			image, 
			raw_date, 
			start_date, 
			end_date
		FROM look_events_db.events
		LIMIT ? OFFSET ?
	`, limit, offset)

	if err != nil {
		return events, err
	}

	defer rows.Close()

	for rows.Next() {
		var event DouEvent

		err = rows.Scan(&event.ID,
			&event.Title,
			&event.ShortDescription,
			&event.FullDescription,
			&event.Location,
			&event.Online,
			&event.Cost,
			&event.Image,
			&event.RawDate,
			&event.Start,
			&event.End)
		if err != nil {
			return events, err
		}

		tags, err := s.getEventTags(event.ID)
		if err != nil {
			return events, err
		}
		event.Tags = tags

		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return events, err
	}

	return events, nil
}

func (s *MySQLStorage) GetEvent(id int) (DouEvent, error) {
	var event DouEvent

	row := s.db.QueryRow(`
		SELECT
			id,
			title, 
			short_description, 
			full_description, 
			location, 
			online, 
			cost, 
			image, 
			raw_date, 
			start_date, 
			end_date
		FROM look_events_db.events
		WHERE id = ?
	`, id)

	err := row.Scan(&event.ID,
		&event.Title,
		&event.ShortDescription,
		&event.FullDescription,
		&event.Location,
		&event.Online,
		&event.Cost,
		&event.Image,
		&event.RawDate,
		&event.Start,
		&event.End)

	if err != nil {
		return event, err
	}

	tags, err := s.getEventTags(id)
	if err != nil {
		return event, err
	}
	event.Tags = tags

	return event, nil
}

func (s *MySQLStorage) getEventTags(eventID int) ([]string, error) {
	tags := make([]string, 0)

	rows, err := s.db.Query(`
		SELECT
		   tag.value
		FROM look_events_db.tags tag
			JOIN look_events_db.tags_events te on tag.id = te.tag_id
		WHERE te.event_id = ?
	`, eventID)
	if err != nil {
		return tags, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return tags, err
	}

	return tags, nil
}

func (s *MySQLStorage) GetTagsList() ([]string, error) {
	tags := make([]string, 0)

	rows, err := s.db.Query(`SELECT value FROM look_events_db.tags`)
	if err != nil {
		return tags, err
	}

	for rows.Next() {
		var tag string

		err = rows.Scan(&tag)
		if err != nil {
			return tags, err
		}

		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return tags, err
	}

	return tags, nil
}
