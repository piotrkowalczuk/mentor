package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"github.com/piotrkowalczuk/sklog"
	"github.com/go-kit/kit/log"
)

const (
	zordonDir = ".zordon"
)

func openPIDFile(name string) (*os.File, error) {
	fpath := filepath.Join(zordonDir, name+".pid")
	if _, err := os.Stat(filepath.Dir(fpath)); os.IsNotExist(err) {
		if err := os.Mkdir(filepath.Dir(fpath), 0777); err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(fpath, os.O_EXCL|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("cannot create %s.pid: %v", name, err)
	}
	return f, nil
}
func killAll(l log.Logger) error {
	fi, err := ioutil.ReadDir(zordonDir)
	if err != nil {
		return err
	}

	for _, f := range fi {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".pid") {
			continue
		}

		fp := filepath.Join(zordonDir, f.Name())
		b, err := ioutil.ReadFile(fp)
		if err != nil {
			return err
		}

		pid, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}

		p, err := os.FindProcess(int(pid))
		if err != nil {
			return err
		}

		if err := p.Kill(); err != nil {
			return err
		}

		sklog.Info(l, fmt.Sprintf("process %s (%d) has been killed", f.Name(), pid))

		if err := os.Remove(fp); err != nil {
			return err
		}
	}

	return nil
}

// getProcess gets a Process from a pid and checks that the
// process is actually running. If the process
// is not running, then getProcess returns a nil
// Process and the error ErrNotRunning.
func getProcess(pid int) (*os.Process, error) {
	p, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	// try to check if the process is actually running by sending
	// it signal 0.
	err = p.Signal(syscall.Signal(0))
	if err == nil {
		return p, nil
	}
	if err == syscall.ESRCH {
		return nil, errors.New("zordon: service is not running")
	}
	return nil, errors.New("server running but inaccessible")
}
