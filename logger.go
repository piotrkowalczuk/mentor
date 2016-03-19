package main

import "github.com/piotrkowalczuk/sklog"

const (
	formatMessage     = "- %-60v "
	formatBraces      = "[%v] "
	formatBracesLevel = "[%-5v] "
	formatSubsystem   = "%-15v "
	formatNone        = ""
	keyColor          = "color"
	keyColorReset     = "color_reset"
)

var formatter = sklog.NewSequentialFormatter(
	sklog.NewKeyFormatter("%s", keyColorReset),
	sklog.NewKeyFormatter("%s", keyColor),
	sklog.NewKeyFormatter(formatBraces, sklog.KeyTimestamp),
	sklog.NewKeyFormatter(formatNone, "time"),
	sklog.NewKeyFormatter(formatSubsystem, sklog.KeySubsystem),
	sklog.NewKeyFormatter(formatBracesLevel, sklog.KeyLevel),
	sklog.NewKeyFormatter(formatBraces, sklog.KeyHTTPMethod),
	sklog.NewKeyFormatter(formatBraces, sklog.KeyHTTPPath),
	sklog.NewKeyFormatter(formatBraces, sklog.KeyHTTPStatus),
	sklog.NewKeyFormatter(formatMessage, sklog.KeyMessage),
)
