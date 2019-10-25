package cmd

import "github.com/spf13/cobra"

func init() {
	webserverCmd.Flags().IntVarP(&port, "port", "p", 3000, "the port on which to run the web server")
	webserverCmd.Flags().IntVarP(&workers, "workers", "w", 4, "number of workers to run the webserver on")
	webserverCmd.Flags().StringVarP(&workerClass, "worker_class", "k", "sync", "the worker class to use")
	webserverCmd.Flags().IntVarP(&workerTimeout, "worker_timeout", "t", 120, "the timeout for waiting on webserver workers")
	webserverCmd.Flags().StringVarP(&hostName, "hostname", "h", "127.0.0.1", "set the hostname on with to run the webserver")
	webserverCmd.Flags().BoolVarP(&daemon, "daemon", "D", false, "deamonize instead of running in the foreground")
	webserverCmd.Flags().StringVarP(&stdout, "stdout", "", "", "redirect stdout to this file")
	webserverCmd.Flags().StringVarP(&stderr, "stderr", "", "", "redirect stderr to this file")
	webserverCmd.Flags().StringVarP(&accessLogFile, "access_logfile", "A", "-", "the logfile to store the webserver access log. use '-' to print to stderr")
	webserverCmd.Flags().StringVarP(&errorLogFile, "error_logfile", "E", "-", "the logfile to store the webserver error log. use '-' to print to stderr.")
	webserverCmd.Flags().StringVarP(&logFile, "log-file", "l", "", "location of the log file")
	webserverCmd.Flags().StringVarP(&sslCert, "ssl_cert", "", "", "path to the SSL certificat wfor the webserver")
	webserverCmd.Flags().StringVarP(&sslKey, "ssl_key", "", "", "path to the key to use with the SSL certificate")
	webserverCmd.Flags().BoolVarP(&debug, "debug", "d", false, "use the server in debug mode")
}

var webserverCmd = &cobra.Command{
	Use:     "webserver",
	Short:   "start a goflow webserver instance",
	PreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE:    func(cmd *cobra.Command, args []string) error { return nil },
}
