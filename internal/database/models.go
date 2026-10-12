package database

import "database/sql"

type Models struct {
	Users      UserModel
	Events     EventModel
	Attendees  AttendeeModel
	Categories CategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:      UserModel{DB: db},
		Events:     EventModel{DB: db},
		Attendees:  AttendeeModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
