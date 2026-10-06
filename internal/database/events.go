package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	Id            int    `json:"id"`
	OwnerId       int    `json:"ownerId"`
	Title         string `json:"title" binding:"required,min=3,max=100"`
	Description   string `json:"description" binding:"required,min=3,max=400"`
	ImageUrl      string `json:"imageUrl"`
	Location      string `json:"location" binding:"required,min=3,max=200"`
	StartDateTime string `json:"startDateTime" binding:"required"`
	EndDateTime   string `json:"endDateTime" binding:"required"`
	CategoryId    string `json:"categoryId" binding:"required"`
	Url           string `json:"url"`
}

func (m *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO events(owner_id, title, description,
		image_url, location, start_date_time, end_date_time,
		category_id, url)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := m.DB.ExecContext(
		ctx,
		query,
		event.OwnerId,
		event.Title,
		event.Description,
		event.ImageUrl,
		event.Location,
		event.StartDateTime,
		event.EndDateTime,
		event.CategoryId,
		event.Url,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	event.Id = int(id)
	return nil

}

func (m *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, title, description, image_url,
	location, start_date_time, end_date_time,
	category_id, url FROM events`

	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.Id,
			&event.OwnerId,
			&event.Title,
			&event.Description,
			&event.ImageUrl,
			&event.Location,
			&event.StartDateTime,
			&event.EndDateTime,
			&event.CategoryId,
			&event.Url,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Get(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, title, description, image_url, location,
	start_date_time, end_date_time,
	category_id, url FROM events WHERE id = ?`

	var event Event

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.Id,
		&event.OwnerId,
		&event.Title,
		&event.Description,
		&event.ImageUrl,
		&event.Location,
		&event.StartDateTime,
		&event.EndDateTime,
		&event.CategoryId,
		&event.Url,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &event, nil
}

func (m *EventModel) GetByOwnerID(ownerID int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, title, description, image_url,
	location, start_date_time, end_date_time,
	category_id, url FROM events WHERE owner_id = ?`

	rows, err := m.DB.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event

		err := rows.Scan(
			&event.Id,
			&event.OwnerId,
			&event.Title,
			&event.Description,
			&event.ImageUrl,
			&event.Location,
			&event.StartDateTime,
			&event.EndDateTime,
			&event.CategoryId,
			&event.Url,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (m *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE events 
		SET title = ?, description = ?, image_url = ?, location = ?,
		start_date_time = ?, end_date_time = ?,
		category_id = ?, url = ? 
		WHERE id = ?
	`
	result, err := m.DB.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.ImageUrl,
		event.Location,
		event.StartDateTime,
		event.EndDateTime,
		event.CategoryId,
		event.Url,
		event.Id,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (m *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "DELETE FROM events WHERE id = ?"

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil

}
