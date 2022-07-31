package auto

import (
	"general_ledger_golang/models"
	"general_ledger_golang/pkg/logger"
)

// Migrate will run migrations automatically, when server starts.
//
// In local, run migrations, once models are ready, get ddl from the
// models and put that in manual migration folder as sql.
//
// Run that in prod. `Never run auto migration in prod.`
func Migrate() {
	db, _ := models.GetDB()
	if err := db.AutoMigrate(models.Book{}, models.Operation{}, models.Posting{}, models.BookBalance{}); err != nil {
		logger.Logger.Fatalf("Automigration failed, error: %+v", err) // fataF is printf followed by panic
	}
	if hasConstraint := db.Migrator().HasConstraint(&models.BookBalance{}, "non_negative_balance"); hasConstraint == false {
		// create this constraint only if it doesn't exist, if any modification needed,
		// drop and create the constraint as postgres doesn't have update constraint provision.
		err := db.Migrator().CreateConstraint(&models.BookBalance{}, "non_negative_balance")
		if err != nil {
			logger.Logger.Fatalf("non_negative_balance check constraint add failed, error: %+v", err)
		}
	}
}
