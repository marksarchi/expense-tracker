package handlers

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/business/data/user"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

type UserGroup struct {
	user user.User
	auth *auth.Auth
}

func (ug *UserGroup) createUser(w http.ResponseWriter, r *http.Request) error {

	var newUser user.NewUser

	if err := web.Decode(r, &newUser); err != nil {
		errors.Wrap(err, "Unable to decode payload")
	}

	user, err := ug.user.CreateUser(newUser)
	if err != nil {
		errors.Wrapf(err, "creating new user : %+v", newUser)

	}
	log.Println(user)
	return web.Respond(w, user, http.StatusCreated)

}
func (ug *UserGroup) Login(w http.ResponseWriter, r *http.Request) error {
	email, password, ok := r.BasicAuth()
	if !ok {
		err := errors.New("Must provide email and password")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}
	claims, err := ug.user.Authenticate(email, password)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "Authenticating")
		}
	}

	// var tkn struct {
	// 	Token string `json:"token"`
	// }

	// tkn.Token, err = ug.auth.GenerateToken(claims)
	// if err != nil {
	// 	return errors.Wrap(err, "generating token")
	// }
	tk, err := auth.CreateToken(claims)
	if err != nil {
		return errors.Wrap(err, "Creating token")
	}

	return web.Respond(w, tk, http.StatusOK)
}
