// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	_ "embed" // Calls init function.

	"github.com/ardanlabs/darwin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sarchimark/expense-tracker/foundation/database"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string

	//go:embed sql/seed.sql
	seedDoc string

	//go:embed sql/delete.sql
	deleteDoc string
)
var migrations = []darwin.Migration{
	{
		Version:     1.1,
		Description: "Create table et_users",
		Script: `
		create table et_users(
			user_id uuid primary key not null,
			first_name varchar(20) not null,
			last_name varchar(20) not null,
			email varchar(30) not null,
			password text not null
		);`,
	},
	{
		Version:     1.2,
		Description: "Create table et_categories",
		Script: `
		create table et_categories(
			category_id integer primary key not null,
			user_id uuid not null,
			title varchar(20) not null,
			description varchar(50) not null
		);`,
	},
	{
		Version:     1.3,
		Description: "Alter table et_categories with user column",
		Script: `

		alter table et_categories add constraint cat_users_fk
		foreign key (user_id) references et_users(user_id);
`,
	},
	{
		Version:     2.1,
		Description: "Create table et_transactions",
		Script: `
		create table et_transactions(
			transaction_id integer primary key not   null ,
			category_id integer not null,
			user_id uuid not null,
			amount numeric(10,2),
			note varchar(50) not null,
			transaction_date bigint not null
		);`,
	},
	{
		Version:     2.2,
		Description: "Alter table et_transactions with user column",
		Script: `
		alter table et_transactions add constraint trans_cat_fk
		foreign  key (category_id) references  et_categories(category_id);
		alter table et_transactions add constraint trans_users_fk
		foreign key (user_id) references et_users(user_id);
		
`,
	},
	{
		Version:     2.3,
		Description: "Add sequence",
		Script: `
		
		create sequence et_categories_seq increment 1 start 1;
		create sequence et_transactions_seq increment 1 start 1000;
		
		
`,
	},
}

const seeds = `INSERT INTO et_users(user_id ,first_name, last_name,email , password) VALUES
('638389b7-7a2b-43f2-9eff-e3596a21f789ykgVqc','mark ', 'sarchi', 'sarchimark@example.com', '$2a$12$JOvqT7bBQlEZ7Cnkiur8teSRU8xwpzH4L9475X3zMpin/u7lBQueK'),
('3d266f28-5d49-4702-9528-9b266afc618aOQGXJa','eva','max' ,'evamax@yahoo.com', '$2a$12$lcaFYoHrKTAFOtaaa5DgKO0b9GEULsAG23Z.q6cy91nFCujy91Z1y')
ON CONFLICT DO NOTHING;

INSERT INTO ET_CATEGORIES ( CATEGORY_ID, USER_ID, TITLE, DESCRIPTION) VALUES
(NEXTVAL('ET_CATEGORIES_SEQ'), 1, 'Travel costs', 'All travel costs'),
(NEXTVAL('ET_CATEGORIES_SEQ'), 2, 'Car Costs', 'all car costs')
ON CONFLICT DO NOTHING;

INSERT INTO ET_TRANSACTIONS (TRANSACTION_ID, CATEGORY_ID, USER_ID, AMOUNT, NOTE, TRANSACTION_DATE) VALUES
(NEXTVAL('ET_TRANSACTIONS_SEQ'), 1, 1, 2580, 'Travel upcountry', 1616014260159428300),
(NEXTVAL('ET_TRANSACTIONS_SEQ'), 1, 1 , 33450, 'Repaired windscreen', 1616097061795640800),
(NEXTVAL('ET_TRANSACTIONS_SEQ'), 2, 2 , 33450, 'Bought new tyre', 1616097061795640800)
ON CONFLICT DO NOTHING;
`

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return errors.Wrap(err, "status check database")
	}
	if err := database.StatusCheck(ctx, db); err != nil {
		return errors.Wrap(err, "status check database")
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return errors.Wrap(err, "construct darwin driver")
	}

	//d := darwin.New(driver, migrations)
	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))
	return d.Migrate()
}

// Seed runs the set of seed-data queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func Seed(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return errors.Wrap(err, "status check database")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seedDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

// DeleteAll runs the set of Drop-table queries against db. The queries are ran in a
// transaction and rolled back if any fail.
func DeleteAll(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(deleteDoc); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}
