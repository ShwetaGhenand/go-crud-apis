// Code generated by sqlc. DO NOT EDIT.

package dbgen

import (
	"database/sql"
)

type User struct {
	ID       int32          `json:"id"`
	Name     string         `json:"name"`
	Password string         `json:"password"`
	Email    string         `json:"email"`
	Phone    sql.NullString `json:"phone"`
	Age      sql.NullInt32  `json:"age"`
	Address  sql.NullString `json:"address"`
}
