package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Owner struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	Id            int    `json:"id"`
	OwnerId       int    `json:"ownerId"`
	Owner         *Owner `json:"owner,omitempty"`
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

	query := `
		SELECT e.id, e.owner_id, e.title, e.description, e.image_url,
		       e.location, e.start_date_time, e.end_date_time,
		       e.category_id, e.url, COALESCE(u.name, '') AS owner_name
		FROM events e
		LEFT JOIN users u ON e.owner_id = u.id
	`
	rows, err := m.DB.QueryContext(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event
		var ownerName string

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
			&ownerName,
		)
		if err != nil {
			return nil, err
		}

		if event.OwnerId > 0 {
			event.Owner = &Owner{
				Id:   event.OwnerId,
				Name: ownerName,
			}
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

	query := `
		SELECT e.id, e.owner_id, e.title, e.description, e.image_url,
		       e.location, e.start_date_time, e.end_date_time,
		       e.category_id, e.url, COALESCE(u.name, '') AS owner_name
		FROM events e
		LEFT JOIN users u ON e.owner_id = u.id
		WHERE e.id = ?
	`

	var event Event
	var ownerName string

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
		&ownerName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if event.OwnerId > 0 {
		event.Owner = &Owner{
			Id:   event.OwnerId,
			Name: ownerName,
		}
	}

	return &event, nil
}

func (m *EventModel) GetByOwnerID(ownerID int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT e.id, e.owner_id, e.title, e.description, e.image_url,
		       e.location, e.start_date_time, e.end_date_time,
		       e.category_id, e.url, COALESCE(u.name, '') AS owner_name
		FROM events e
		LEFT JOIN users u ON e.owner_id = u.id
		WHERE e.owner_id = ?
	`

	rows, err := m.DB.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []*Event{}

	for rows.Next() {
		var event Event
		var ownerName string

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
			&ownerName,
		)
		if err != nil {
			return nil, err
		}

		if event.OwnerId > 0 {
			event.Owner = &Owner{
				Id:   event.OwnerId,
				Name: ownerName,
			}
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
