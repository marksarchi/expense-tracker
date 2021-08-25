package handlers

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/business/data/category"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

type categoryGroup struct {
	category category.Category
	auth     *auth.Auth
}

// swagger:route POST /api/categories category createCategory
// creates a new category
// responses:
//  201:categoryResponse

//createCategory creates a new category
func (cg categoryGroup) createCategory(w http.ResponseWriter, r *http.Request) error {

	var nc category.NewCategory
	if err := web.Decode(r, &nc); err != nil {
		return errors.Wrap(err, "unable to decode payload")
	}
	userID, _ := userID(r)
	cat, err := cg.category.Create(nc, userID)
	if err != nil {
		return errors.Wrapf(err, "creating new category : %+v", nc)
	}

	return web.Respond(w, cat, http.StatusCreated)

}

func (cg categoryGroup) getCategories(w http.ResponseWriter, r *http.Request) error {

	userID, err := userID(r)
	if err != nil {
		return errors.Wrap(err, "UserID not found")
	}

	categories, err := cg.category.GetAllCategories(userID)
	if err != nil {
		return err
	}
	return web.Respond(w, categories, http.StatusOK)
}
func (cg categoryGroup) getCategoryByID(w http.ResponseWriter, r *http.Request) error {

	userID, err := userID(r)
	if err != nil {
		return errors.Wrap(err, "UserID not found")
	}

	log.Printf("UserID : %s", userID)

	categoryID := web.GetURLParamInt64(w, r, "categoryId")
	cat, err := cg.category.GetCategoryByID(userID, int(categoryID))
	if err != nil {
		switch err {
		case category.ErrForbidden:
			return web.NewRequestError(err, http.StatusBadRequest)
		case category.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		default:
			return errors.Wrapf(err, "ID: %s", categoryID)
		}
	}
	return web.Respond(w, cat, http.StatusOK)

}
func (cg categoryGroup) updateCategory(w http.ResponseWriter, r *http.Request) error {
	userID, err := userID(r)
	if err != nil {
		return errors.Wrap(err, "Error fetching UserID")
	}
	categoryID := int(web.GetURLParamInt64(w, r, "categoryId"))
	var up category.UpdateCategory
	if err := web.Decode(r, &up); err != nil {
		return errors.Wrapf(err, "unable to decode payload")
	}

	if err := cg.category.UpdateCategory(userID, categoryID, up); err != nil {
		switch err {
		case category.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case category.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case category.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "ID: %s  User: %+v", userID, &up)
		}
	}
	return nil
}
func (cg categoryGroup) removeCategory(w http.ResponseWriter, r *http.Request) error {
	userID, err := userID(r)
	if err != nil {
		return errors.Wrap(err, "Error fetching UserID")
	}
	categoryID := int(web.GetURLParamInt64(w, r, "categoryId"))
	if err := cg.category.DeleteCategory(userID, categoryID); err != nil {
		switch err {
		case category.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "ID : %s", categoryID)
		}
	}

	return web.Respond(w, nil, http.StatusNetworkAuthenticationRequired)
}

// func userID(r *http.Request) string {
// 	claimsMap := r.Context().Value(auth.Key).(jwt.MapClaims)
// 	userID := claimsMap["user_id"].(string)
// 	//userID, err := strconv.Atoi(idstr)
// 	// if err != nil {
// 	// 	errors.Wrap(err, "converting claims string to id")
// 	// }
// 	return userID
// }

func userID(r *http.Request) (string, error) {
	claims, ok := r.Context().Value(auth.Key).(auth.Claims)
	if !ok {
		return "", web.NewShutdownError("web value missing from context")
	}
	return claims.Subject, nil

}
