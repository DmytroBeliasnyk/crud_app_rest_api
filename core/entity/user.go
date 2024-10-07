package entity

import "github.com/DmytroBeliasnyk/crud_app_rest_api/core/dto"

type User struct {
	Id    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`

	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

func FromSignUpDTO(u dto.SignUpDTO, passwordHash string) *User {
	return &User{
		Name:         u.Name,
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: passwordHash,
	}
}
