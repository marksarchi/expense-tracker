package handlers

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/business/data/transaction"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

type transactionGroup struct {
	transaction transaction.Transaction
}

func (tg transactionGroup) addTransaction(w http.ResponseWriter, r *http.Request) error {
	var nt transaction.NewTransaction
	err := web.Decode(r, &nt)
	if err != nil {
		errors.Wrap(err, "Unable to decode payload")
	}

	categoryId := web.GetURLParamInt64(w, r, "categoryId")

	userID := tg.userID(r)

	transaction, err := tg.transaction.AddTransaction(nt, userID, int(categoryId))

	if err != nil {
		return errors.Wrapf(err, "Transaction: %+v", &transaction)
	}
	return web.Respond(w, transaction, http.StatusCreated)

}
func (tg transactionGroup) getTransactionById(w http.ResponseWriter, r *http.Request) error {

	userId := tg.userID(r)
	categoryID := int(web.GetURLParamInt64(w, r, "categoryId"))
	transactionID := int(web.GetURLParamInt64(w, r, "transactionId"))

	trans, err := tg.transaction.GetTransactionByID(userId, categoryID, transactionID)
	if err != nil {
		switch err {
		case transaction.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case transaction.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "%s", transactionID)
		}

	}
	return web.Respond(w, trans, http.StatusOK)

}

func (tg transactionGroup) getAllTransactions(w http.ResponseWriter, r *http.Request) error {
	userID := tg.userID(r)
	categoryID := web.GetURLParamInt64(w, r, "categoryId")

	categories, err := tg.transaction.GetAllTransactions(userID, int(categoryID))
	if err != nil {
		return err
	}
	return web.Respond(w, categories, http.StatusOK)

}
func (tg transactionGroup) updateTransaction(w http.ResponseWriter, r *http.Request) error {
	userID := tg.userID(r)
	categoryID := web.GetURLParamInt64(w, r, "categoryId")
	transactionID := web.GetURLParamInt64(w, r, "transactionId")
	var transUp transaction.UpdateTransaction
	if err := web.Decode(r, &transUp); err != nil {
		return errors.Wrapf(err, "Unable to decode payload")
	}
	if err := tg.transaction.UpdateTransaction(userID, int(categoryID), int(transactionID), transUp); err != nil {
		switch err {
		case transaction.ErrInvalidID:
			web.NewRequestError(err, http.StatusBadRequest)
		case transaction.ErrNotFound:
			web.NewRequestError(err, http.StatusNotFound)
		case transaction.ErrForbidden:
			web.NewRequestError(err, http.StatusForbidden)
		default:
			errors.Wrapf(err, "updating transaction", transUp)
		}

	}
	return web.Respond(w, nil, http.StatusNoContent)

}

func (tg transactionGroup) removeTransaction(w http.ResponseWriter, r *http.Request) error {
	userID := tg.userID(r)
	categoryID := web.GetURLParamInt64(w, r, "categoryId")
	transactionID := web.GetURLParamInt64(w, r, "transactionId")

	if err := tg.transaction.RemoveTransactionByID(userID, int(categoryID), int(transactionID)); err != nil {
		switch err {
		case transaction.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID :%s", transactionID)
		}
	}
	return web.Respond(w, nil, http.StatusNoContent)
}

func (tg transactionGroup) userID(r *http.Request) int {
	claimsMap := r.Context().Value(auth.Key).(jwt.MapClaims)
	idstr := claimsMap["user_id"].(string)
	userID, err := strconv.Atoi(idstr)
	if err != nil {
		errors.Wrap(err, "converting claims string to id")
	}
	return userID
}
