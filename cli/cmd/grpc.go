package cmd

import (
	"os"
	"path"
	"runtime"

	"git.edenfarm.id/edenlabs/cli/app"
	"git.edenfarm.id/edenlabs/cli/log"
	"github.com/spf13/cobra"
)

var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Build & Running your project as a GRPC server",
	Long:  `Build & Running your project as a GRPC server`,
	Run: func(c *cobra.Command, args []string) {
		RunGRPC()
	},
}

func RunGRPC() {
	log.Log = log.New("Edenlabs - GRPC")

	gps := getGoPath()
	if len(gps) == 0 {
		log.Log.Errorln("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	app.PrintHeader("Starting GRPC server...")

	var paths []string
	readDirectory(cwd, &paths)

	var files []string
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	app.Watch(appName, paths, files, "grpc")
	app.Build("grpc")

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}
