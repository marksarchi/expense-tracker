package category

//NewCategory represents information needed to create a new category
type NewCategory struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

//CategoryInfo represents a category

type CategoryInfo struct {
	ID           int     `db:"category_id" json:"category_id"`
	Title        string  `db:"title" json:"title"`
	Description  string  `db:"description" json:"description"`
	UserID       string  `db:"user_id" json:"user_id"`
	TotalExpense float64 `db:"total_expense" json:"total_expense"`
}

//UpdateCategory has information needed to update a category
type UpdateCategory struct {
	Title       *string `db:"title" json:"title"`
	Description *string `db:"description" json:"description"`
}
