package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	connID       string
	connType     string
	connHost     string
	connLogin    string
	connPassword string
	connSchema   string
	connPort     int
	connURI      string
	connExtra    string
)

func init() {
	connectionCmd.AddCommand(connectionListCmd)
	connectionCmd.AddCommand(connectionAddCmd)
	connectionCmd.AddCommand(connectionDeleteCmd)
}

var connectionCmd = &cobra.Command{
	Use:   "connection",
	Short: "list/add/delete connections",
}

var connectionListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list connections",
	PreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("listing connections")
		return nil
	},
}

func init() {
	connectionAddCmd.Flags().StringVarP(&connID, "conn_id", "", "", "connection name")
	connectionAddCmd.Flags().StringVarP(&connType, "conn_type", "", "", "connection type")
	connectionAddCmd.Flags().StringVarP(&connURI, "conn_uri", "", "", "connection uri")
}

var connectionAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add connections",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if connID == "" {
			return errors.New("must supply conn_id")
		}
		if connType == "" {
			return errors.New("must supply conn_type")
		}
		if connURI == "" {
			return errors.New("must supply conn_uri")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("adding connection")
		return nil
	},
}

var connectionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete connections",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if connID == "" {
			return errors.New("must supply conn_id")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("deleting connection")
		return nil
	},
}
