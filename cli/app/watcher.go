package app

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"

	"git.edenfarm.id/edenlabs/cli/log"
	"github.com/howeyc/fsnotify"
)

var (
	// appName, packages name, get from directory name.
	appName string

	// cmd is external command.
	cmd *exec.Cmd

	// state is mutual exclusion lock,
	state sync.Mutex

	// watching, list of extention file that need to watch.
	watching []string

	// lastBuild, time last build performed.
	lastBuild time.Time

	// isStarted, if true then application is on running.
	isStarted = make(chan bool)

	// runTime, slices of runnable application unix times.
	runTime = make(map[string]int64)
)

// Watch performs initializing fsnotify to watch files on current directory.
// and trigger to rebuild and restart the packages after some changes has been made
// on files which that we watch.
func Watch(appname string, paths []string, files []string, command string) {
	appName = appname
	watching = append(files, ".go")
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Log.Errorln("Failed to create new Watcher")
		log.Log.Error(err)
		os.Exit(2)
	}

	go func() {
		for {
			select {
			case e := <-watcher.Event:
				isbuild := true
				if !isWatched(e.Name) {
					continue
				}

				if lastBuild.Add(1 * time.Second).After(time.Now()) {
					continue
				}

				lastBuild = time.Now()
				mt := lastModified(e.Name)
				if t := runTime[e.Name]; mt == t {
					log.Log.Warningf("Skipped # %s #\n", e.String())
					isbuild = false
				}

				runTime[e.Name] = mt
				if isbuild {
					log.Log.Info(e)
					go Build(command)
				}
			case err := <-watcher.Error:
				log.Log.Error(err)
			}
		}
	}()

	log.Log.Infoln("Initializing watcher...")
	for i, path := range paths {
		bar := progress(i+1, len(paths), 42)
		time.Sleep(time.Millisecond * 2)
		os.Stdout.Write([]byte(bar + "\r"))

		err = watcher.Watch(path)
		if err != nil {
			log.Log.Errorln("Failed to watch directory.")
			log.Log.Error(err)
			os.Exit(2)
		}

		os.Stdout.Sync()
	}

	os.Stdout.Write([]byte("\n"))
}

// Build performs executing go build of packages.
func Build(command string) {
	var err error

	state.Lock()
	defer state.Unlock()

	path, _ := os.Getwd()
	os.Chdir(path)

	// run go lint, testing and vet
	// runTest()

	// using script from go-fast-build
	// github.com/kovetskiy/go-fast
	cmd := exec.Command("/bin/sh", "-c", "export GOBIN=$(pwd); exec go install -gcflags \"-trimpath $GOPATH/src\" \"$@\";")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Log.Errorln("Failed to build.")
		log.Log.Error(err)
		return
	}

	restart(appName, command)
}

func runTest() {
	log.Log.Infoln("Running vet, lint and test ...")
	cmd := exec.Command("/bin/sh", "-c", "go test -cover;")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Log.Errorln("Failed to run test.")
		log.Log.Error(err)
	}
}

// restart performs restarting application binary.
func restart(app string, command string) {
	log.Log.Infof("Restarting %s ...", app)
	kill()

	go start(app, command)
}

// kill performs killing current running application.
func kill() {
	defer func() {
		if e := recover(); e != nil {
			log.Log.Errorf("Kill.recover -> %s", e)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
	}
}

// start performs to start the application binary.
func start(app string, command string) {
	log.Log.Infof("Rebuild %s ...", app)

	if strings.Index(app, "./") == -1 {
		app = "./" + app
	}

	cmd = exec.Command(app, command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go cmd.Run()
	log.Log.Infof("Started %s ...", app)
	isStarted <- true
}

// lastModified returns unix timestamp of file last modified.
func lastModified(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		log.Log.Errorf("Cannot find file [ %s ]\n", err)
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Log.Errorf("Failed to get information from file [ %s ]\n", err)
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

// isWatched returns if ext files was on watching.
func isWatched(fileName string) bool {
	for _, s := range watching {
		if strings.HasSuffix(fileName, s) {
			return true
		}
	}

	return false
}

// progress returns a string as progress bar from scaning directory.
func progress(current, total, cols int) string {
	prefix := strconv.Itoa(current) + "/" + strconv.Itoa(total)
	bar_start := " ["
	bar_end := "] "

	bar_size := cols - len(prefix+bar_start+bar_end)
	amount := int(float32(current) / (float32(total) / float32(bar_size)))
	remain := bar_size - amount

	bar := strings.Repeat("#", amount) + strings.Repeat(" ", remain)
	return "🔍 Scanning (\033[1m" + prefix + "\033[0m)" + bar_start + bar + bar_end
}
