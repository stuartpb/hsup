package hsup

import (
	"errors"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"bitbucket.org/kardianos/osext"
)

var ErrNoReleases = errors.New("No releases found")

type Notifier interface {
	Notify() <-chan *Processes
}

type Processes struct {
	Rel   *Release
	Forms []Formation

	Dd        DynoDriver
	OneShot   bool
	Executors []*Executor
}

type Formation interface {
	Args() []string
	Quantity() int
	Type() string
}

func linuxAmd64Path() string {
	exe, err := osext.Executable()
	if err != nil {
		log.Fatalf("could not locate own executable:", err)
	}

	if runtime.GOOS == "linux" && runtime.GOARCH == "amd64" {
		return exe
	}

	return exe + "-linux-amd64"
}

func copyFile(src, dst string, mode os.FileMode) (err error) {
	r, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() {
		err = combine(err, r.Close())
	}()

	w, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		err = combine(err, w.Close())
	}()

	if _, err := io.Copy(w, r); err != nil {
		return err
	}
	if err := os.Chmod(dst, mode); err != nil {
		return err
	}
	return nil
}

type combinedError []error

func combine(errors ...error) combinedError {
	errors = make(combinedError, 0, len(errors))
	for _, err := range errors {
		if err == nil {
			continue
		}
		errors = append(errors, err)
	}
	return errors
}

func (e combinedError) Error() string {
	msgs := make([]string, len(e))
	for i, entry := range e {
		msgs[i] = entry.Error()
	}
	return strings.Join(msgs, " | ")
}
