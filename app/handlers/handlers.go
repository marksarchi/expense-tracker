package handlers

import (
	//"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/auth"
	"github.com/sarchimark/expense-tracker/business/data/category"
	"github.com/sarchimark/expense-tracker/business/data/transaction"
	"github.com/sarchimark/expense-tracker/business/data/user"
	"github.com/sarchimark/expense-tracker/business/mid"
	"github.com/sarchimark/expense-tracker/foundation/web"
)

func toHandler(fn http.HandlerFunc) http.Handler {
	return http.HandlerFunc(fn)
}
func SetupRoutes(db *sqlx.DB, log *log.Logger, shutdown chan os.Signal) http.Handler {
	auth, err := auth.New("RS256")
	if err != nil {
		log.Println(errors.Wrap(err, "constructing auth"))
	}

	app := web.NewApp(shutdown, mid.Logger(log), mid.Errors(log))
	cg := categoryGroup{
		category: category.New(log, db),
		auth:     auth,
	}
	app.Handle("POST", "/api/categories", cg.createCategory, mid.Authenticate(auth))
	app.Handle("GET", "/api/categories", cg.getCategories, mid.Authenticate(auth))
	app.Handle("GET", "/api/categories/{categoryId}", cg.getCategoryByID, mid.Authenticate(auth))
	app.Handle("PUT", "/api/categories/{categoryId}", cg.updateCategory, mid.Authenticate(auth))
	app.Handle("DELETE", "/api/categories/{categoryId}", cg.removeCategory, mid.Authenticate(auth))

	tg := transactionGroup{
		transaction: transaction.New(log, db),
	}
	app.Handle("POST", "/api/categories/{categoryId}/transactions", tg.addTransaction, mid.Authenticate(auth))
	app.Handle("GET", "/api/categories/{categoryId}/transactions/{transactionId}", tg.getTransactionById, mid.Authenticate(auth))
	app.Handle("GET", "/api/categories/{categoryId}/transactions", tg.getAllTransactions, mid.Authenticate(auth))
	app.Handle("PUT", "/api/categories/{categoryId}/transactions/{transactionId}", tg.updateTransaction, mid.Authenticate(auth))
	app.Handle("DELETE", "/api/categories/{categoryId}/transactions/{transactionId}", tg.removeTransaction, mid.Authenticate(auth))

	ug := UserGroup{
		user: user.New(log, db),
		auth: auth,
	}

	app.Handle("POST", "/users/signup", ug.createUser)
	app.Handle("POST", "/users/login", ug.Login)

	return app

}
