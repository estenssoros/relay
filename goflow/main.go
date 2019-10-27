package main

import (
	"time"

	"github.com/estenssoros/goflow/cmd"
	"github.com/sirupsen/logrus"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	logrus.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
}

func main() {
	start := time.Now()
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("process took %v", time.Since(start))
}
