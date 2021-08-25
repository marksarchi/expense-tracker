package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/data/schema"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

var ErrHelp = errors.New("provided help")

func Migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return errors.Wrap(err, "Connect to database")
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return errors.Wrap(err, "migrate database")
	}
	fmt.Println("Migrations complete")
	return nil
}
