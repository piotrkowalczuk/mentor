package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/piotrkowalczuk/sklog"
)

func recruit(ctx *cli.Context) {
	af, err := openAlphasfile(alphasFile)
	if err != nil {
		log.Fatal(err)
	}

	sklog.Log(logger, sklog.KeyMessage, "Alpha, Rita's escaped! Recruit a team of services with attitude!", sklog.KeyLevel, sklog.LevelWarning, sklog.KeySubsystem, "zordon")
	sklog.Log(logger, sklog.KeyMessage, "Understood, Zordon!", sklog.KeySubsystem, "alpha", sklog.KeyLevel, sklog.LevelInfo)

	for _, r := range af.Rangers {
		install := exec.Command("go", "get", "-t", r.Import)
		if err = run(install, r, logger); err != nil {
			sklog.Fatal(logger, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
		}
		sklog.Info(logger, fmt.Sprintf("%s ready!", r.Name), sklog.KeySubsystem, "alpha")
	}

	sklog.Log(logger, sklog.KeyMessage, "Use extreme caution Rangers, you are dealing with an evil here that is beyond all imagination!", sklog.KeySubsystem, "zordon", sklog.KeyLevel, sklog.LevelInfo)
}
