package main

import (
	"encoding/gob"
	"fmt"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/plugin"
)

func init() {
	gob.Register(plugin.Command{})
}

type p struct{}

func (p) Commands() []plugin.Command {
	return nil
}

func (p) Hooks() []plugin.Hook {
	return []plugin.Hook{
		{
			Name:        "hooking",
			PlaceHookOn: "chain init",
		},
	}
}

func (p) Execute(cmd plugin.Command, args []string) error {
	return nil
}

func (p) ExecuteHookPre(hook plugin.Hook, args []string) error {
	fmt.Println("pre")
	return nil
}

func (p) ExecuteHookPost(hook plugin.Hook, args []string) error {
	fmt.Println("post")
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.Hook, args []string) error {
	fmt.Println("clean")
	return nil
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"hooking": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
