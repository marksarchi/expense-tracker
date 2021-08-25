//Package classification of Expensetracker API
//
//Documentation for expense-tracker API
//
// Schemes: http
// BasePath:/
// Version: 1.0.0
//
// Consumes:
//  - application/json
// swagger:meta

package data

import (
	"github.com/sarchimark/expense-tracker/business/data/category"
)

//Information about a category
//swagger:response categoryResponse
type CategoryResponse struct {
	//Newly created category
	//in:body
	Body category.CategoryInfo
}
