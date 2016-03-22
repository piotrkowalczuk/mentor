package main

import (
	"fmt"
	"log"
	"os/exec"

	"strings"

	"github.com/codegangsta/cli"
	"github.com/piotrkowalczuk/sklog"
)

func powerup(ctx *cli.Context) {
	af, err := openAlphasfile(alphasFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(af.Service) == 0 {
		return
	}

	sklog.Warning(serviceLogger(logger, af.Service[0]), "We need Thunderzord power now!", sklog.KeySubsystem, af.Service[0].Name)

	for _, s := range af.Service {
		rl := serviceLogger(logger, s)
		modified, err := isGitModifiedLocaly(s)
		if err != nil {
			sklog.Fatal(rl, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
		}
		if modified {
			sklog.Warning(logger, fmt.Sprintf("Alpha and I will have to analyze your powers, %s, to see if they can be restored to you permanently.", s.Name), sklog.KeySubsystem, "zordon", keyColorReset, colorReset)
			continue
		}

		if s.Branch == "" || s.Branch == "master" {
			update := exec.Command("go", "get", "-u", "-t", s.Import)
			if err := run(update, s, rl); err != nil {
				sklog.Fatal(rl, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
			}
		} else {
			if err := updateRepository(s); err != nil {
				sklog.Fatal(rl, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
			}
		}
		sklog.Info(rl, fmt.Sprintf("%s Thunderzord Power!", strings.ToTitle(s.Name)), sklog.KeySubsystem, s.Name)
	}

	sklog.Log(logger, sklog.KeyMessage, "ThunderZords shall be yours, powerful and agile. When joined together, all shall form the Thunder MegaZord.", sklog.KeySubsystem, "zordon", sklog.KeyLevel, sklog.LevelInfo, keyColorReset, colorReset)
}
