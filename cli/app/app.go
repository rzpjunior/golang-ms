package app

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"git.edenfarm.id/edenlabs/cli/log"
)

type DocVal string

func (d *DocVal) String() string {
	return fmt.Sprint(*d)
}

func (d *DocVal) Set(value string) error {
	*d = DocVal(value)
	return nil
}

type ListOpts []string

func (opts *ListOpts) String() string {
	return fmt.Sprint(*opts)
}

func (opts *ListOpts) Set(value string) error {
	*opts = append(*opts, value)
	return nil
}

type StubTemplate struct {
	AppPath           string
	ProjectPath       string
	PackagePath       string
	PackageName       string
	ModuleName        string
	ModelName         string
	ModelNameSingular string
	ModelNamePlural   string
	TableName         string
}

func GetPackagePath(currentPath string) string {
	gp := os.Getenv("GOPATH")
	log.Log.Debugf("gopath:%s", gp)
	if gp == "" {
		log.Log.Errorln("you should set GOPATH in the env")
		os.Exit(2)
	}

	appPath := ""
	haspath := false
	for _, wg := range filepath.SplitList(gp) {
		wg, _ = filepath.EvalSymlinks(path.Join(wg, "src"))

		if filepath.HasPrefix(strings.ToLower(currentPath), strings.ToLower(wg)) {
			haspath = true
			appPath = wg
			break
		}
	}

	if !haspath {
		log.Log.Errorf("Can't generate application code outside of GOPATH '%s'", gp)
		os.Exit(2)
	}

	return strings.Join(strings.Split(currentPath[len(appPath)+1:], string(filepath.Separator)), "/")
}

func GetGoPath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		log.Log.Errorln("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	log.Log.Debugf("gopath:%s", gopath)
	return gopath
}

func FormatSourceCode(filename string) {
	cmd := exec.Command("gofmt", "-w", filename)
	if err := cmd.Run(); err != nil {
		log.Log.Errorf("gofmt err: %s\n", err)
	}
}

func MakeDir(appPath string) {
	var ps = []struct {
		Name string
		Path string
	}{
		{"project", appPath},
		{"datastore", path.Join(appPath, "datastore")},
		{"model", path.Join(appPath, "datastore", "model")},
		{"repository", path.Join(appPath, "datastore", "repository")},
		{"engine", path.Join(appPath, "engine")},
		{"test", path.Join(appPath, "test")},
		{"src", path.Join(appPath, "src")},
	}

	PrintHeader("Generating project directory stucture.")

	for _, p := range ps {
		log.Log.Infof("%-20s : \t\t%s", p.Name+" directory", p.Path)
		os.Mkdir(p.Path, 0755)
	}

	PrintFooter()
}

func WriteFile(file *os.File, content string, tpl *StubTemplate) {
	if tpl != nil {
		content = StubReplaces(content, tpl)
	}

	if _, err := file.WriteString(content); err != nil {
		log.Log.Errorf("Could not write file %s\n%s", file.Name(), err)
		os.Exit(2)
	}

	file.Close()
}

func StubReplaces(content string, tpl *StubTemplate) string {
	content = strings.Replace(content, "{{ProjectPath}}", tpl.ProjectPath, -1)
	content = strings.Replace(content, "{{PackagePath}}", tpl.PackagePath, -1)
	content = strings.Replace(content, "{{PackageName}}", tpl.PackageName, -1)
	content = strings.Replace(content, "{{ModuleName}}", tpl.ModuleName, -1)
	content = strings.Replace(content, "{{ModelName}}", tpl.ModelName, -1)
	content = strings.Replace(content, "{{ModelNameSingular}}", tpl.ModelNameSingular, -1)
	content = strings.Replace(content, "{{ModelNamePlural}}", tpl.ModelNamePlural, -1)
	content = strings.Replace(content, "{{TableName}}", tpl.TableName, -1)

	return content
}

func DirExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Log.Fatal(err)
	}
	ok := []string{"y", "Y", "yes", "Yes", "YES"}
	notOk := []string{"n", "N", "no", "No", "NO"}
	if ContainsString(ok, response) {
		return true
	} else if ContainsString(notOk, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

func ContainsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}
	return false
}

func FileReader(file string) (f *os.File, err error) {
	if DirExist(file) {
		log.Log.Warnf("%v is exist, do you want to overwrite it? Yes or No?", file)
		if AskForConfirmation() {
			if f, err = os.OpenFile(file, os.O_RDWR|os.O_TRUNC, 0666); err != nil {
				log.Log.Error(err)
				return
			}
		} else {
			log.Log.Infoln("Skip creating file.")
			return
		}
	} else {
		if f, err = os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0666); err != nil {
			log.Log.Error(err)
			return
		}
	}

	return
}

func PrintHeader(msg string) {
	log.Log.Infoln(msg)
	log.Log.Infoln("--------------------------------------")
}

func PrintFooter() {
	log.Log.Infoln("--------------------------------------")
	log.Log.Infoln("")
}
