package main

import (
	"HCPlatform/code/cmd"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {

	err := cmd.Execute()
	if err != nil {
		return
	}

}

func init() {
	// Output to stdout instead of the default stderr
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	//If you wish to add the calling method as a field, instruct the logger via:
	log.SetReportCaller(true)
}
