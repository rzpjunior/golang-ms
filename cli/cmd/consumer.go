package cmd

import (
	"os"
	"path"
	"runtime"

	"git.edenfarm.id/edenlabs/cli/app"
	"git.edenfarm.id/edenlabs/cli/log"
	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Build & Running your project as a Consumer",
	Long:  `Build & Running your project as a Consumer`,
	Run: func(c *cobra.Command, args []string) {
		RunConsumer()
	},
}

func RunConsumer() {
	log.Log = log.New("Edenlabs - Consumer")

	gps := getGoPath()
	if len(gps) == 0 {
		log.Log.Errorln("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	app.PrintHeader("Starting Consumer server...")

	var paths []string
	readDirectory(cwd, &paths)

	var files []string
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	app.Watch(appName, paths, files, "consumer")
	app.Build("consumer")

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
