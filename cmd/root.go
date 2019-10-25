package cmd

import (
	"github.com/spf13/cobra"
)

var (
	port          int
	workers       int
	workerClass   string
	workerTimeout int
	hostName      string
	pid           int
	daemon        bool
	stdout        string
	stderr        string
	accessLogFile string
	errorLogFile  string
	logFile       string
	sslCert       string
	sslKey        string
	debug         bool
)

func init() {
	rootCmd.AddCommand(initDBCmd)
	rootCmd.AddCommand(connectionCmd)
	rootCmd.AddCommand(webserverCmd)
}

var rootCmd = &cobra.Command{
	Use:   "goflow",
	Short: "",
}

func Execute() error {
	return rootCmd.Execute()
}
