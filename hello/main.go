package main

import (
	"encoding/gob"
	"fmt"
	"path/filepath"

	hplugin "github.com/hashicorp/go-plugin"

	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
)

func init() {
	gob.Register(plugin.Manifest{})
	gob.Register(plugin.ExecutedCommand{})
	gob.Register(plugin.ExecutedHook{})
}

type p struct{}

func (p) Manifest() (plugin.Manifest, error) {
	return plugin.Manifest{
		Name: "hello",
		Commands: []plugin.Command{
			{
				Use:               "hello",
				Short:             "Say hello",
				PlaceCommandUnder: "ignite",
			},
		},
		Hooks: []plugin.Hook{
			{
				Name:        "hello",
				PlaceHookOn: "ignite chain init",
			},
		},
	}, nil
}

func (p) Execute(cmd plugin.ExecutedCommand) error {
	fmt.Println("Hello from plugin")
	return nil
}

func (p) ExecuteHookPre(hook plugin.ExecutedHook) error {
	fmt.Println("Hello chain init")
	return nil
}

func (p) ExecuteHookPost(hook plugin.ExecutedHook) error {
	return nil
}

func (p) ExecuteHookCleanUp(hook plugin.ExecutedHook) error {
	return nil
}

func getChain(cmd plugin.ExecutedCommand, chainOption ...chain.Option) (*chain.Chain, error) {
	var (
		home, _ = cmd.Flags().GetString("home")
		path, _ = cmd.Flags().GetString("path")
	)
	if home != "" {
		chainOption = append(chainOption, chain.HomePath(home))
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	return chain.New(absPath, chainOption...)
}

func main() {
	pluginMap := map[string]hplugin.Plugin{
		"hello": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
