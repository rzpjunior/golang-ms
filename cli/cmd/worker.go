package cmd

import (
	"os"
	"path"
	"runtime"

	"git.edenfarm.id/edenlabs/cli/app"
	"git.edenfarm.id/edenlabs/cli/log"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Build & Running your project as a worker",
	Long:  `Build & Running your project as a worker`,
	Run: func(c *cobra.Command, args []string) {
		RunWorker()
	},
}

func RunWorker() {

	log.Log = log.New("Edenlabs - WOrker")

	gps := getGoPath()
	if len(gps) == 0 {
		log.Log.Errorln("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	app.PrintHeader("Starting Worker server...")

	var paths []string
	readDirectory(cwd, &paths)

	var files []string
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	app.Watch(appName, paths, files, "worker")
	app.Build("worker")

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
