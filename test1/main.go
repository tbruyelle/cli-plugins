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
	gob.Register(plugin.Command{})
}

type p struct{}

func (p) Commands() ([]plugin.Command, error) {
	// TODO: write your command list here
	cmds := []plugin.Command{
		{
			Use:               "test1",
			Short:             "Explain what the command is doing...",
			Long:              "Long description goes here...",
			PlaceCommandUnder: "ignite",
			// Examples of subcommands:
			Commands: []plugin.Command{
				{Use: "add"},
				{Use: "list"},
				{Use: "delete"},
			},
		},
	}
	return cmds, nil
}

func (p) Execute(cmd plugin.Command, args []string) error {
	// TODO: write command execution here
	fmt.Printf("Hello I'm the test1 plugin!\nargs=%v, with=%v\n", args, cmd.With)

	// This is how the plugin can access the chain:
	c, err := getChain(cmd)
	if err != nil {
		return err
	}
	_ = c

	// According to the number of declared commands, you may need a switch:
	switch cmd.Use {
	case "add":
		fmt.Println("Adding stuff...")
	case "list":
		fmt.Println("Listing stuff...")
	case "delete":
		fmt.Println("Deleting stuff...")
	}
	return nil
}

func getChain(cmd plugin.Command, chainOption ...chain.Option) (*chain.Chain, error) {
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
		"test1": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
