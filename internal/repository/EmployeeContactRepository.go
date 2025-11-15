package repository

import (
	"context"
	"diplom_back/internal/entity"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetAllEmployeeContacts(ctx context.Context, db *pgxpool.Pool) ([]entity.EmployeeContact, error) {
	var contacts []entity.EmployeeContact

	rows, err := db.Query(ctx, `select * from employee_contact`)
	if err != nil {
		slog.Warn("ERROR GETTING EMPLOYEE CONTACTS")
	}
	defer rows.Close()
	for rows.Next() {
		var contact entity.EmployeeContact
		err := rows.Scan(
			&contact.Id,
			&contact.RoomNum,
			&contact.Title,
			&contact.Comment)
		if err != nil {
			slog.Warn("GetAllEmployeeContacts", err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	if err := rows.Err(); err != nil {
		slog.Warn("GetAllEmployeeContacts", err)
		return nil, err
	}
	return contacts, nil
}

func PostNewEmployeeContact(ctx context.Context, db *pgxpool.Pool, contact *entity.NewEmployeeContact) (int, error) {
	var id int
	err := db.QueryRow(ctx, `INSERT INTO employee_contact (room_num, title, comment)
			VALUES ($1, $2, $3) RETURNING id`, contact.RoomNum, contact.Title, contact.Comment).Scan(&id)

	return id, err
}

func DeleteEmployeeContact(ctx context.Context, db *pgxpool.Pool, contactId uint) error {
	_, err := db.Exec(ctx, `DELETE FROM employee_contact WHERE id = $1`, contactId)
	return err
}
