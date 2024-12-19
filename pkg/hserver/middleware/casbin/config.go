package casbin

import (
	"embed"
	_ "embed"
)

//go:embed model.conf
var modelConf embed.FS
