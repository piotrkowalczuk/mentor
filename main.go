package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	klog "github.com/go-kit/kit/log"
	"github.com/hashicorp/hcl"
	"github.com/mgutz/ansi"
	"github.com/piotrkowalczuk/sklog"
)

var (
	gopath string
	dir    string
	colors = []string{
		ansi.LightGreen,
		ansi.LightYellow,
		ansi.LightBlue,
		ansi.LightMagenta,
		ansi.LightCyan,
	}
	colorLog   = ansi.ColorCode("black:green")
	colorError = ansi.ColorCode("black:red")
	colorReset = ansi.ColorCode("reset")
)

func init() {
	var err error
	gopath = os.Getenv("GOPATH")
	dir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	logger := sklog.NewHumaneLogger(os.Stdout, formatter)
	b, err := ioutil.ReadFile("Mentorfile")
	if err != nil {
		log.Fatalf("mentor: Mentorfile read error: %s", err.Error())
	}

	var mf Mentorfile
	if err = hcl.Unmarshal(b, &mf); err != nil {
		log.Fatalf("mentor: Mentorfile parsing error: %s", err.Error())
	}

	logf("Mentorfile successfully parsed")

	for i, s := range mf.Services {
		if s.Color == "" {
			s.Color = colors[i%len(colors)]
		}
		install := exec.Command("go", "get", s.Path)
		l := klog.NewContext(logger).With(keyColor, s.Color, keyColorReset, colorReset)
		if err = ex(install, s, l); err != nil {
			errorf("mentor", "Mentorfile parsing error: %s%s", err.Error())
		}

		logf("%s has been installed", s.Name)
	}

	end := make(chan struct{}, 1)
	for _, ser := range mf.Services {
		<-time.After(1 * time.Second)
		go runService(ser, logger)
	}
	<-end
}

func runService(s *Service, l klog.Logger) {
	l = klog.NewContext(l).With(sklog.KeySubsystem, s.Name, keyColor, s.Color, keyColorReset, colorReset)
	for {
		cmd := exec.Command(s.Name, JoinArgs(s.Arguments)...)

		err := ex(cmd, s, l)
		if err != nil {
			if cmd.ProcessState.Exited() {
				errorf(s.Name, "service will be restarted because of error: %s", err.Error())
				continue
			}

			errorf(s.Name, "service has stoped with error: %s", err.Error())
			return
		}
	}
}

func ex(c *exec.Cmd, s *Service, l klog.Logger) error {
	var err error

	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}

	if err = c.Start(); err != nil {
		return err
	}

	go sc(stdout, s, l)
	sc(stderr, s, l)

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
		errorf("scan error: %s", err.Error())
	}
}

func logf(s string, args ...interface{}) {
	fmt.Printf(" %s MENTOR %s %s \n", colorLog, colorReset, fmt.Sprintf(s, args...))
}

func errorf(n, s string, args ...interface{}) {
	fmt.Printf(" %s %s %s %s \n", colorError, strings.ToUpper(n), colorReset, fmt.Sprintf(s, args...))
}
