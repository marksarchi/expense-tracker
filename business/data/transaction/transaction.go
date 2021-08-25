package transaction

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

var (
	// ErrNotFound is used when a specific Product is requested but does not exist.
	ErrNotFound = errors.New("not found")

	// ErrInvalidID occurs when an ID is not in a valid form.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("attempted action is not allowed")
)

// Transaction manages the set of API's for transaction access
type Transaction struct {
	db  *sqlx.DB
	log *log.Logger
}

//New returns a new instance of type Transaction
func New(log *log.Logger, db *sqlx.DB) Transaction {
	return Transaction{
		log: log,
		db:  db,
	}
}

//AddTransaction creates a new transaction
func (t Transaction) AddTransaction(nt NewTransaction, userID string, CategoryID int) (Info, error) {
	q := `INSERT INTO ET_TRANSACTIONS (TRANSACTION_ID, CATEGORY_ID, USER_ID, AMOUNT, NOTE, TRANSACTION_DATE) 
	VALUES(NEXTVAL('ET_TRANSACTIONS_SEQ'), $1, $2, $3, $4, $5)`

	tr := Info{
		CategoryID:      CategoryID,
		UserID:          userID,
		Amount:          nt.Amount,
		Note:            nt.Note,
		TransactionDate: int(time.Now().UTC().UnixNano()),
	}

	_, err := t.db.Exec(q, tr.CategoryID, tr.UserID, tr.Amount, tr.Note, tr.TransactionDate)
	if err != nil {
		return Info{}, errors.Wrap(err, "Inserting transaction")

	}
	return tr, nil

}

//GetTransactionByID finds a single transaction identified by given userId, categoryId and transactionId
func (t Transaction) GetTransactionByID(userID string, categoryID, transactionID int) (Info, error) {

	q := `SELECT TRANSACTION_ID, CATEGORY_ID, USER_ID, AMOUNT, NOTE, TRANSACTION_DATE FROM ET_TRANSACTIONS 
	WHERE USER_ID = $1 AND CATEGORY_ID  = $2 AND TRANSACTION_ID = $3`

	var tr Info
	if err := t.db.Get(&tr, q, userID, categoryID, transactionID); err != nil {
		if err == sql.ErrNoRows {
			return Info{}, ErrNotFound
		}
		return Info{}, errors.Wrap(err, "Selecting transaction")
	}

	return tr, nil

}

//GetAllTransactions finds transactions identified by given userId and categoryId
func (t Transaction) GetAllTransactions(userID string, categoryID int) ([]Info, error) {
	q := `SELECT TRANSACTION_ID, CATEGORY_ID, USER_ID, AMOUNT, NOTE, TRANSACTION_DATE FROM ET_TRANSACTIONS
	 WHERE USER_ID = $1 AND CATEGORY_ID =$2`

	transactions := []Info{}
	if err := t.db.Select(&transactions, q, userID, categoryID); err != nil {
		return []Info{}, errors.Wrap(err, "selecting transactions")
	}
	return transactions, nil

}

//UpdateTransaction updates a single transaction
func (t Transaction) UpdateTransaction(userID string, categoryID, transactionID int, ut UpdateTransaction) error {
	q := `UPDATE ET_TRANSACTIONS SET  AMOUNT = $4, NOTE = $5" +
"WHERE USER_ID = $1 AND CATEGORY_ID = $2 AND TRANSACTION_ID = $3"`
	trans, err := t.GetTransactionByID(userID, categoryID, transactionID)
	if err != nil {
		return ErrNotFound
	}

	if ut.Amount != nil {
		trans.Amount = *ut.Amount
	}

	if ut.Note != nil {
		trans.Note = *ut.Note
	}

	t.log.Printf("%s: %s", "transaction.Update", database.Log(q, trans.TransactionID, trans.Amount, trans.Note))

	if _, err := t.db.Exec(q, userID, categoryID, transactionID, trans.Amount, trans.Note); err != nil {
		return errors.Wrap(err, "updating transaction")
	}
	return nil

}

//RemoveTransactionByID  deletes a single transaction
func (t Transaction) RemoveTransactionByID(userID string, categoryID, transactionID int) error {
	q := `DELETE FROM ET_TRANSACTIONS WHERE USER_ID = $1 AND CATEGORY_ID = $2 AND TRANSACTION_ID = $3`
	t.log.Printf("%s: %s", "transaction.Delete",
		database.Log(q, transactionID),
	)

	if _, err := t.db.Exec(q, userID, categoryID, transactionID); err != nil {
		return errors.Wrapf(err, "deleting transaction %s", transactionID)

	}
	return nil

}
