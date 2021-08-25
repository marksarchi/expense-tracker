package user

import (
	//"database/sql"
	"database/sql"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	//uuid "github.com/satori/go.uuid"
	"github.com/google/uuid"

	//"github.com/dgrijalva/jwt-go"
	"github.com/sarchimark/expense-tracker/business/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrNotFound is used when a specific User is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrAuthenticationFailure occurs when a user attempts to authenticate but
	// anything goes wrong.
	ErrAuthenticationFailure = errors.New("authentication failed")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

type User struct {
	db  *sqlx.DB
	log *log.Logger
}

//New create a new instance of mark
func New(log *log.Logger, db *sqlx.DB) User {
	return User{
		db:  db,
		log: log,
	}

}

//CreateUser creates a new user
func (u User) CreateUser(newUser NewUser) (Info, error) {
	q := "INSERT INTO ET_USERS (USER_ID ,FIRST_NAME, LAST_NAME,EMAIL , PASSWORD) VALUES($1 ,$2 ,$3, $4,$5)"

	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 12)
	if err != nil {
		return Info{}, errors.Wrap(err, "generating hash password")
	}

	userInfo := Info{
		ID:           uuid.New().String(),
		FirstName:    newUser.FirstName,
		SecondName:   newUser.SecondName,
		Email:        newUser.Email,
		PasswordHash: hash,
	}
	log.Printf("User : %s", userInfo)

	_, err = u.db.Exec(q, userInfo.ID, userInfo.FirstName, userInfo.SecondName, userInfo.Email, userInfo.PasswordHash)
	if err != nil {
		return Info{}, errors.Wrapf(err, "Inserting NewUser")
	}

	return userInfo, nil

}

//Authenticates users
func (u User) Authenticate(email, password string) (auth.Claims, error) {

	q := "SELECT * FROM et_users WHERE email =  $1"

	u.log.Printf("user.Authenticate : %s", email)
	var usr Info

	if err := u.db.Get(&usr, q, email); err != nil {
		if err == sql.ErrNoRows {
			return auth.Claims{}, errors.Wrap(err, "selecting user")
		}
		return auth.Claims{}, ErrAuthenticationFailure
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "Expense tracker project",
			Subject:   usr.ID,
			Audience:  "students",
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	// atClaims := jwt.MapClaims{}
	// atClaims["authorized"] = true
	// atClaims["user_id"] = usr.ID
	// atClaims["exp"] = time.Now().Add(time.Hour).Unix()

	return claims, nil

}

func (u User) QueryByID() (Info, error) {
	userID := "3d266f28-5d49-4702-9528-9b266afc618a"
	const q = `
	SELECT
		*
	FROM
		et_users
	WHERE 
		user_id = $1`

	var userInfo Info
	if err := u.db.Get(&userInfo, q, userID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrapf(err, "selecting user %q", userID)
	}
	return userInfo, nil
}
