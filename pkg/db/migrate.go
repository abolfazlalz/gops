package db

import "command-server/pkg/models"

func Migrate() error {
	db, err := Database()
	if err != nil {
		return err
	}
	return db.AutoMigrate(
		models.CommandLog{},
		models.Bash{},
	)
}
