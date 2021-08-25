package category

import (
	//"fmt"

	//"database/sql"
	"database/sql"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	//uuid "github.com/satori/go.uuid"
	"github.com/google/uuid"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Category manages the set of API's for category access
type Category struct {
	db  *sqlx.DB
	log *log.Logger
}

//New creates a new instance of type Category
func New(log *log.Logger, db *sqlx.DB) Category {
	return Category{
		log: log,
		db:  db,
	}
}

//Create adds a new Category
func (cg Category) Create(newCategory NewCategory, userID string) (CategoryInfo, error) {
	q := "INSERT INTO ET_CATEGORIES ( CATEGORY_ID, USER_ID, TITLE, DESCRIPTION) VALUES(NEXTVAL('ET_CATEGORIES_SEQ'), $1, $2, $3)"

	cat := CategoryInfo{
		UserID:      userID,
		Title:       newCategory.Title,
		Description: newCategory.Description,
	}

	_, err := cg.db.Exec(q, cat.UserID, cat.Title, cat.Description)
	if err != nil {
		cg.log.Println(err)
		return CategoryInfo{}, errors.Wrap(err, "inserting category")

	}

	return cat, nil

}

//GetAllCategories Returns all categories with specified userID
func (cg Category) GetAllCategories(userID string) ([]CategoryInfo, error) {
	if _, err := uuid.Parse(userID); err != nil {
		return []CategoryInfo{}, ErrInvalidID
	}

	const q = `
	SELECT 
	    c.* ,
	    COALESCE(SUM(t.amount),0) AS total_expense 
	FROM 
	    et_categories as c
	RIGHT OUTER JOIN 
	    et_transactions AS t ON c.category_id = t.category_id
	WHERE 
	    c.user_id = $1
	GROUP BY c.category_id`

	categories := []CategoryInfo{}
	err := cg.db.Select(&categories, q, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Selecting categories")
	}
	return categories, nil
}

//GetCategoryByID Returns all categories with specified userID and categoryID
func (cg Category) GetCategoryByID(userID string, categoryID int) (CategoryInfo, error) {

	q := `
	SELECT 
	    c.* ,
	    COALESCE(SUM(t.amount),0) AS total_expense 
	FROM 
	    et_categories as c
	RIGHT OUTER JOIN 
	    et_transactions AS t ON c.category_id = t.category_id
	WHERE  
	    c.user_id = $1 AND c.category_id = $2 
	GROUP BY c.category_id`

	var ci CategoryInfo

	if _, err := uuid.Parse(userID); err != nil {
		return CategoryInfo{}, ErrInvalidID
	}

	if err := cg.db.Get(&ci, q, userID, categoryID); err != nil {
		if err == sql.ErrNoRows {
			return CategoryInfo{}, ErrNotFound
		}
		return CategoryInfo{}, errors.Wrap(err, "selecting category")

	}

	return ci, nil

}

//DeleteCategory deletes a category
func (cg Category) DeleteCategory(userID string, categoryID int) error {
	q := `DELETE FROM ET_CATEGORIES WHERE USER_ID = $1 AND CATEGORY_ID = $2`

	_, err := cg.db.Exec(q, userID, categoryID)
	if err != nil {
		return errors.Wrapf(err, "Deleting category", categoryID)
	}
	return nil
}

//UpdateCategory update a single category
func (cg Category) UpdateCategory(userID string, categoryID int, uc UpdateCategory) error {
	q := `UPDATE ET_CATEGORIES SET TITLE = $1, DESCRIPTION = $2 " +
	"WHERE USER_ID = $3 AND CATEGORY_ID = $4`

	cat, err := cg.GetCategoryByID(userID, categoryID)
	if err != nil {
		return err
	}

	if uc.Description != nil {
		cat.Description = *uc.Description
	}
	if uc.Title != nil {
		cat.Title = *uc.Title
	}

	if _, err := cg.db.Exec(q, cat.Title, cat.Description, userID, categoryID); err != nil {
		return errors.Wrap(err, "updating product")
	}

	return nil

}
