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
		errors.Wrapf(err, "creating new user : %+v", &user)

	}
	log.Println(&user)
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
func (ug UserGroup) queryByID(w http.ResponseWriter, r *http.Request) error {
	usr, err := ug.user.QueryByID()

	if err != nil {
		switch err {
		case user.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case user.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case user.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s", "3d266f28-5d49-4702-9528-9b266afc618a")
		}

	}
	return web.Respond(w, usr, http.StatusOK)
}
