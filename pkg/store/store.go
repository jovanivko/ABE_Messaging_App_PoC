package store

import (
	"abeProofOfConcept/pkg/store/models"
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var dsn = "postgres://jovan:1234@localhost:5432/company?sslmode=disable"
var DB *bun.DB

func ConnectDB() *bun.DB {

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	// Create a Bun database instance with the PostgreSQL dialect
	DB = bun.NewDB(sqldb, pgdialect.New())

	DB.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
	))
	return DB
}

func TableExists(ctx context.Context, tableName string) bool {
	var exists bool
	_ = DB.NewSelect().
		ColumnExpr("EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)", tableName).
		Scan(ctx, &exists)

	return exists
}

func createUserTable(ctx context.Context) error {
	_, err := DB.NewCreateTable().Model((*models.User)(nil)).IfNotExists().Exec(ctx)
	return err
}

func createMessagesTable(ctx context.Context) error {
	_, err := DB.NewCreateTable().Model((*models.Message)(nil)).IfNotExists().Exec(ctx)
	return err
}

func createMsgFragmentTable(ctx context.Context) error {
	_, err := DB.NewCreateTable().Model((*models.MessageFragment)(nil)).IfNotExists().
		ForeignKey(`("msg_id") REFERENCES "messages" ("id") ON DELETE CASCADE`).Exec(ctx)
	return err
}

func createKeysTable(ctx context.Context) error {
	_, err := DB.NewCreateTable().Model((*models.MasterKey)(nil)).IfNotExists().Exec(ctx)
	return err
}

func CreateDatabase(ctx context.Context) error {
	if DB == nil {
		return errors.New("database connection not established")
	}
	err := createUserTable(ctx)
	if err != nil {
		return err
	}
	err = createMessagesTable(ctx)
	if err != nil {
		return err
	}
	err = createMsgFragmentTable(ctx)
	if err != nil {
		return err
	}
	err = createKeysTable(ctx)
	if err != nil {
		return err
	}
	return nil
}
