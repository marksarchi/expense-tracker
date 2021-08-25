package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/business/data/schema"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

func Seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		errors.Wrap(err, "connect to database")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := schema.Seed(ctx, db); err != nil {
		return errors.Wrap(err, "seed database")
	}

	fmt.Println("seed data is complete")
	return nil

}
