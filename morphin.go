package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"time"

	"strings"

	"github.com/codegangsta/cli"
	klog "github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

func morphin(ctx *cli.Context) {
	af, err := openAlphasfile(alphasFile)
	if err != nil {
		log.Fatal(err)
	}

	sklog.Log(logger, sklog.KeyMessage, "Rangers, you must act swiftly, the development environment is in grave danger!", sklog.KeyLevel, sklog.LevelWarning, sklog.KeySubsystem, "zordon")

	for _, r := range af.Service {
		l := klog.NewContext(logger).With(keyColor, r.Color, keyColorReset, colorReset)
		sklog.Info(l, fmt.Sprintf("%s!!!", strings.ToUpper(r.Name)), sklog.KeySubsystem, r.Name)
	}

	end := make(chan struct{}, 1)
	for _, r := range af.Service {
		<-time.After(1 * time.Second)
		go morphRanger(r, logger)
	}
	<-end
}

func morphRanger(s *Service, l klog.Logger) {
	rl := serviceLogger(l, s)
	for {
		cmd := exec.Command(s.Name, JoinArgs(s.Arguments)...)

		if err := run(cmd, s, rl); err != nil {
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				sklog.Error(rl, fmt.Errorf("service will be restarted because of error: %s", err.Error()))
				continue
			}

			sklog.Error(rl, fmt.Errorf("service has stoped with error: %s", err.Error()))
			return
		}
	}
}

func run(c *exec.Cmd, s *Service, l klog.Logger) error {
	var (
		err            error
		stderr, stdout io.ReadCloser
		multi          io.Reader
	)
	stderr, err = c.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err = c.StdoutPipe()
	if err != nil {
		return err
	}
	multi = io.MultiReader(stdout, stderr)

	if err = c.Start(); err != nil {
		return err
	}

	sc(multi, s, l)

	if err = c.Wait(); err != nil {
		return err
	}

	return nil
}

func sc(rc io.Reader, s *Service, l klog.Logger) {
	in := bufio.NewScanner(rc)
	tmp := map[string]interface{}{}
ScanLoop:
	for in.Scan() {
		switch s.Log {
		case "json":
			if !bytes.HasPrefix(in.Bytes(), []byte("{")) {
				sklog.Log(l, sklog.KeyMessage, in.Text())
				continue ScanLoop

			}
			if err := json.Unmarshal(in.Bytes(), &tmp); err != nil {
				sklog.Log(l, sklog.KeyMessage, in.Text(), "error", err.Error())
				continue ScanLoop
			}
			arr := make([]interface{}, 0, len(tmp)*2)
			for k, v := range tmp {
				arr = append(arr, k, v)
			}
			sklog.Log(l, append(arr)...)
		default:
			sklog.Log(l, sklog.KeyMessage, in.Text())
		}
	}
	if err := in.Err(); err != nil {
		sklog.Error(l, err, sklog.KeySubsystem, "alpha")
	}
}
