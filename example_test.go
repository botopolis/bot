package gobot_test

import (
	"fmt"

	"github.com/berfarah/gobot"
)

type ExamplePlugin struct {
	Username string
}

func (p *ExamplePlugin) Load(r *gobot.Robot) { p.Username = "beardroid" }

func ExamplePluginUsage() {
	robot := gobot.New(NewChat())
	robot.Install(
		&ExamplePlugin{},
	)

	var plugin ExamplePlugin
	if ok := robot.Plugin(&plugin); ok {
		fmt.Println(plugin.Username)
		// "beardroid"
	}
}
