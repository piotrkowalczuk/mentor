package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/codegangsta/cli"
	"github.com/piotrkowalczuk/sklog"
)

func recruit(ctx *cli.Context) error {
	af, err := openAlphasfile(alphasFile)
	if err != nil {
		log.Fatal(err)
	}

	sklog.Log(logger, sklog.KeyMessage, "Alpha, Rita's escaped! Recruit a team of services with attitude!", sklog.KeyLevel, sklog.LevelWarning, sklog.KeySubsystem, "zordon")
	sklog.Log(logger, sklog.KeyMessage, "Understood, Zordon!", sklog.KeySubsystem, "alpha", sklog.KeyLevel, sklog.LevelInfo)

	for _, s := range af.Service {
		goget := exec.Command("go", "get", "-t", "-d", s.Import)
		if err = run(goget, s, logger); err != nil {
			sklog.Fatal(logger, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
		}

		if s.Branch != "" && s.Branch != "master" {
			if err := updateRepository(s); err != nil {
				sklog.Fatal(logger, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
			}
		}

		sklog.Info(logger, fmt.Sprintf("%s ready", s.Name), sklog.KeySubsystem, "alpha", "branch", s.Branch)
	}

	sklog.Log(logger, sklog.KeyMessage, "Use extreme caution Rangers, you are dealing with an evil here that is beyond all imagination!", sklog.KeySubsystem, "zordon", keyColorReset, colorReset, sklog.KeyLevel, sklog.LevelInfo)
	return nil
}
