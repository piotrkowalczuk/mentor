package main

import (
	"os"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/go-kit/kit/log"
	"github.com/mgutz/ansi"
	"github.com/piotrkowalczuk/sklog"
)

var (
	gopath     string
	alphasFile string
	colors     = []string{
		ansi.LightGreen,
		ansi.LightYellow,
		ansi.LightBlue,
		ansi.LightMagenta,
		ansi.LightCyan,
	}
	logger     log.Logger
	colorReset = ansi.ColorCode("reset")
)

func init() {
	logger = sklog.NewHumaneLogger(os.Stdout, formatter)
}

func main() {
	app := cli.NewApp()
	app.Name = "zordon"
	app.Usage = "Defends development environment from Rita, and her endless waves of containers!"
	app.Authors = []cli.Author{
		{
			Name:  "Piotr Kowalczuk",
			Email: "p.kowalczuk.priv@gmail.com",
		},
	}
	app.Version = "0.1.0"
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "gopath",
			EnvVar:      "GOPATH",
			Destination: &gopath,
		},
		cli.StringFlag{
			Name:        "alphasfile",
			Value:       "Alphasfile",
			Destination: &alphasFile,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:        "morphintime",
			Aliases:     []string{"mt"},
			Description: "Rangers, you must act swiftly, the development environment is in grave danger!",
			Action:      morphin,
			Before:      summon,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "install",
				},
			},
		},
		{
			Name:        "recruit",
			Aliases:     []string{"r"},
			Description: "Alpha, Rita's escaped! Recruit a team of services with attitude!",
			Action:      recruit,
			Before:      summon,
		},
		{
			Name:        "powerup",
			Aliases:     []string{"pu"},
			Description: "We need Thunderzord power now!",
			Action:      powerup,
			Before:      summon,
		},
	}

	app.Run(os.Args)
}

func summon(ctx *cli.Context) error {
	// TODO: implement, something goes wrong
	// fmt.Println(ctx.Command)
	// sklog.Log(logger, sklog.KeyMessage, ctx.Command.Description, sklog.KeyLevel, sklog.LevelFatal, sklog.KeySubsystem, "zordon")
	return nil
}

func src(gopath, pkg string) string {
	return gopath + "/src/" + pkg
}

func serviceLogger(l log.Logger, s *Service) log.Logger {
	return log.NewContext(l).With(sklog.KeySubsystem, s.Name, keyColor, s.Color, keyColorReset, colorReset)
}

func isGitModifiedLocaly(s *Service) (bool, error) {
	check := exec.Command("git", "-C", src(gopath, s.Import), "diff", "--exit-code")
	if err := run(check, s, log.NewNopLogger()); err != nil {
		if check.ProcessState.Exited() {
			return true, nil
		}
		return false, err
	}

	return false, nil
}

func updateRepository(s *Service) (err error) {
	fetch := exec.Command("git", "-C", src(gopath, s.Import), "fetch", "-q", "origin", s.Branch)
	if err = run(fetch, s, logger); err != nil {
		return
	}
	checkout := exec.Command("git", "-C", src(gopath, s.Import), "checkout", "-q", s.Branch)
	if err = run(checkout, s, logger); err != nil {
		return
	}
	pull := exec.Command("git", "-C", src(gopath, s.Import), "pull", "-q", "origin", s.Branch)
	if err = run(pull, s, logger); err != nil {
		return
	}

	return
}
