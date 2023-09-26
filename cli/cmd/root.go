package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"git.edenfarm.id/edenlabs/cli/app"
	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "edenlabs",
		Short: "Tool for easier serve some instance as a you need.",
		Long:  `Tool for easier serve some instance as a you need.`,
		Run: func(c *cobra.Command, args []string) {
			c.HelpFunc()(c, args)
		},
	}

	command.CompletionOptions.HiddenDefaultCmd = false
	command.DisableFlagParsing = true
	command.DisableFlagsInUseLine = true
	command.CompletionOptions = cobra.CompletionOptions{
		DisableDefaultCmd:   true,
		DisableNoDescFlag:   true,
		DisableDescriptions: true,
		HiddenDefaultCmd:    true,
	}
	command.ResetFlags()
	command.AddCommand(runCmd, grpcCmd, consumerCmd, versionCmd)

	b := &bytes.Buffer{}

	headerCLI := `

███████╗██████╗░███████╗███╗░░██╗██╗░░░░░░█████╗░██████╗░░██████╗
██╔════╝██╔══██╗██╔════╝████╗░██║██║░░░░░██╔══██╗██╔══██╗██╔════╝
█████╗░░██║░░██║█████╗░░██╔██╗██║██║░░░░░███████║██████╦╝╚█████╗░
██╔══╝░░██║░░██║██╔══╝░░██║╚████║██║░░░░░██╔══██║██╔══██╗░╚═══██╗
███████╗██████╔╝███████╗██║░╚███║███████╗██║░░██║██████╦╝██████╔╝
╚══════╝╚═════╝░╚══════╝╚═╝░░╚══╝╚══════╝╚═╝░░╚═╝╚═════╝░╚═════╝░
`

	helpText := `%s%s
%s%s

%s%s%s

  %sedenlabs%s [command]

%s%s%s

  %srun%s         Build & Running your project as a HTTP server
  %sgrpc%s        Build & Running your project as a GRPC server
  %sconsumer%s    Build & Running your project as a Consumer
  %shelp%s        Help about any command
  %sversion%s     Print the version number of Edenlabs

%sUse "edenlabs [command] --help" for more information about a command.

`

	fmt.Fprintf(b, helpText, ansi.LightCyan, headerCLI, ansi.LightYellow, command.Short, ansi.LightYellow, "Usage:", ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightYellow, "Command", ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightCyan, ansi.Reset, ansi.LightYellow)

	command.SetHelpTemplate(string(b.Bytes()))

	return command
}

var mainFiles app.ListOpts

func readDirectory(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), "docs") {
			continue
		}

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirectory(directory+"/"+fileInfo.Name(), paths)
			continue
		}

		if useDirectory == true {
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}

	return
}

// getGoPath returns list of go path on system.
func getGoPath() (p []string) {
	gopath := os.Getenv("GOPATH")
	p = strings.Split(gopath, ":")

	return
}
