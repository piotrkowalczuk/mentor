package main

import (
	"fmt"
	"log"
	"os/exec"

	"strings"

	"github.com/codegangsta/cli"
	klog "github.com/go-kit/kit/log"
	"github.com/piotrkowalczuk/sklog"
)

func powerup(ctx *cli.Context) {
	af, err := openAlphasfile(alphasFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(af.Rangers) == 0 {
		return
	}

	sklog.Warning(rangerLogger(logger, af.Rangers[0]), "We need Thunderzord power now!", sklog.KeySubsystem, af.Rangers[0].Name)

	for _, r := range af.Rangers {
		rl := rangerLogger(logger, r)
		check := exec.Command("git", "-C", src(gopath, r.Import), "diff", "--exit-code")
		if err := run(check, r, klog.NewNopLogger()); err != nil {
			if check.ProcessState.Exited() {
				sklog.Warning(logger, fmt.Sprintf("Alpha and I will have to analyze your powers, %s, to see if they can be restored to you permanently.", r.Name), sklog.KeySubsystem, "zordon", keyColorReset, colorReset)
				continue
			}
			sklog.Fatal(rl, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
		}

		update := exec.Command("go", "get", "-u", "-t", r.Import)
		if err := run(update, r, rl); err != nil {
			sklog.Fatal(rl, fmt.Errorf("Ayiyiyiyi!: %s", err.Error()), sklog.KeySubsystem, "alpha")
		}
		sklog.Info(rl, fmt.Sprintf("%s Thunderzord Power!", strings.ToTitle(r.Name)), sklog.KeySubsystem, r.Name)
	}

	sklog.Log(logger, sklog.KeyMessage, "ThunderZords shall be yours, powerful and agile. When joined together, all shall form the Thunder MegaZord.", sklog.KeySubsystem, "zordon", sklog.KeyLevel, sklog.LevelInfo, keyColorReset, colorReset)
}
