package user

// Info represents an individual user.
type Info struct {
	ID           string `db:"user_id" json:"id"`
	FirstName    string `db:"first_name" json:"firstname"`
	SecondName   string `db:"last_name" json:"secondname"`
	Email        string `db:"email" json:"email"`
	PasswordHash []byte `db:"password" json:"-"`
	// DateCreated  time.Time `db:"date_created" json:"date_created"`
	// DateUpdated  time.Time `db:"date_updated" json:"date_updated"`
}

// NewUser contains information needed to create a new User.
type NewUser struct {
	FirstName       string `json:"firstname" validate:"required"`
	SecondName      string `json:"secondname" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"eqfield=Password"`
}
