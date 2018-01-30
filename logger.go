package bot

import (
	"os"

	"github.com/op/go-logging"
)

var stdout = logging.AddModuleLevel(
	logging.NewBackendFormatter(
		logging.NewLogBackend(os.Stdout, "", 0),
		logging.MustStringFormatter(`%{time:2006-01-02T15:04:05} bot ▶ %{color}%{level:.4s}%{color:reset} %{message}`),
	),
)

func newLogger() *logging.Logger {
	log := logging.MustGetLogger("bot")
	stdout.SetLevel(logging.INFO, "")
	log.SetBackend(stdout)

	return log
}
