package entity

import "time"

type Account struct {
	Id        int       `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
}
