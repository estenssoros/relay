package cmd

import (
	"github.com/estenssoros/relay/db"
	"github.com/estenssoros/relay/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var initDBCmd = &cobra.Command{
	Use:     "initdb",
	Short:   "initialize testing database",
	PreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := db.Connect()
		if err != nil {
			return errors.Wrap(err, "gorm open")
		}
		defer db.Close()
		if err := db.AutoMigrate(models.Migrations...).Error; err != nil {
			return errors.Wrap(err, "automigrate")
		}
		return nil
	},
}
