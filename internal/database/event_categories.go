package database

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var ErrInvalidCategory = errors.New("invalid or missing category")

// erstellt ein Event und löst dessen Kategorie innerhalb derselben Transaktion auf
func (m *EventModel) CreateWithCategory(ctx context.Context, event *Event) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	categoryID, err := resolveCategory(ctx, tx, event)
	if err != nil {
		return err
	}

	query := `
	INSERT INTO events (
		owner_id, title, description, image_url, location,
		start_date_time, end_date_time, category_id, url
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`

	result, err := tx.ExecContext(
		ctx,
		query,
		event.OwnerId,
		event.Title,
		event.Description,
		event.ImageUrl,
		event.Location,
		event.StartDateTime,
		event.EndDateTime,
		categoryID,
		event.Url,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	event.Id = int(id)
	event.CategoryId = categoryID
	event.CategoryName = ""

	return nil

}

// Aktualisiert ein Event und löst dessen Kategorie innerhalb derselben Transaktion auf
// Eine ungenutzte benutzerdefinierte Kategorie wird entfernt
func (m *EventModel) UpdateWithCategory(ctx context.Context, event *Event) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldCategoryID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT category_id FROM events WHERE id = ?`,
		event.Id,
	).Scan(&oldCategoryID)

	if err != nil {
		return err
	}

	newCategoryID, err := resolveCategory(ctx, tx, event)
	if err != nil {
		return err
	}

	query := `
	UPDATE events
	SET title = ?, description = ?, image_url = ?, location = ?,
	    start_date_time = ?, end_date_time = ?, category_id = ?, url = ?
	WHERE id = ?
`

	result, err := tx.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.ImageUrl,
		event.Location,
		event.StartDateTime,
		event.EndDateTime,
		newCategoryID,
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

	if oldCategoryID != newCategoryID {
		if err := deleteUnusedCustomCategory(ctx, tx, oldCategoryID); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	event.CategoryId = newCategoryID
	event.CategoryName = ""

	return nil

}

// löscht ein Event und entfernt dessen benutzerdefinierte Kategorie
// falls kein anderes Event diese Kategorie verwendet
func (m *EventModel) DeleteWithCategory(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var categoryID string
	err = tx.QueryRowContext(
		ctx,
		`SELECT category_id FROM events WHERE id = ?`,
		id,
	).Scan(&categoryID)

	if err != nil {
		return err
	}

	result, err := tx.ExecContext(
		ctx,
		`DELETE FROM events WHERE id = ?`,
		id,
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

	if err := deleteUnusedCustomCategory(ctx, tx, categoryID); err != nil {
		return err
	}

	return tx.Commit()

}

// Gibt eine vorhandene Kategorie-ID zurück oder erstellt eine benutzerdefinierte Kategorie
// CategoryId || CategoryName, nicht beides
func resolveCategory(
	ctx context.Context,
	tx *sql.Tx,
	event *Event,
) (string, error) {
	categoryName := strings.TrimSpace(event.CategoryName)
	categoryID := strings.TrimSpace(event.CategoryId)

	if categoryName != "" && categoryID != "" {
		return "", ErrInvalidCategory
	}

	if categoryName != "" {
		var existingID string

		err := tx.QueryRowContext(
			ctx,
			`SELECT id FROM categories WHERE name = ?`,
			categoryName,
		).Scan(&existingID)

		if err == nil {
			return existingID, nil
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}

		randomBytes := make([]byte, 16)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}
		newID := hex.EncodeToString(randomBytes)

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO categories (id, name, is_default) VALUES (?, ?, 0)`,
			newID,
			categoryName,
		)
		if err != nil {
			return "", err
		}

		return newID, nil
	}

	if categoryID == "" {
		return "", ErrInvalidCategory
	}

	var existingID string
	err := tx.QueryRowContext(
		ctx,
		`SELECT id FROM categories WHERE id = ?`,
		categoryID,
	).Scan(&existingID)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCategory
	}
	if err != nil {
		return "", err
	}

	return existingID, nil

}

// Standardkategorie wird niemals gelöscht
func deleteUnusedCustomCategory(
	ctx context.Context,
	tx *sql.Tx,
	categoryID string,
) error {

	_, err := tx.ExecContext(
		ctx,
		`DELETE FROM categories WHERE id = ? AND is_default = 0 AND NOT EXISTS ( SELECT 1 FROM events WHERE category_id = ? )`,
		categoryID,
		categoryID,
	)

	return err

}
