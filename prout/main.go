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
	gob.Register(plugin.Hook{})
}

type p struct{}

func (p) Commands() ([]plugin.Command, error) {
	// TODO: write your command list here
	// Here the default is a single prout command added to the root ignite
	// command.
	cmd := plugin.Command{
		Use:               "prout",
		Short:             "Explain what the command is doing...",
		Long:              "Long description goes here...",
		PlaceCommandUnder: "ignite",
		// Examples of adding subcommands:
		/*
			Commands: []plugin.Command{
				{Use: "add"},
				{Use: "list"},
				{Use: "delete"},
			},
		*/
	}
	// Example of adding flags
	cmd.Flags().String("my-flag", "", "a flag example")

	return []plugin.Command{cmd}, nil
}

func (p) Hooks() ([]plugin.Hook, error) {
	return []plugin.Hook{}, nil
}

func (p) Execute(cmd plugin.Command, args []string) error {
	// TODO: write command execution here
	fmt.Printf("Hello I'm the prout plugin\n")
	fmt.Printf("My executed command: %s\n", cmd.Use)
	fmt.Printf("My args: %v\n", args)
	myFlag, _ := cmd.Flags().GetString("my-flag")
	fmt.Printf("My flags: my-flag=%s\n", myFlag)
	fmt.Printf("My config parameters: %v\n", cmd.With)

	// This is how the plugin can access the chain:
	c, err := getChain(cmd)
	if err != nil {
		return err
	}
	_ = c

	// According to the number of declared commands, you may need a switch:
	/*
		switch cmd.Use {
		case "add":
			fmt.Println("Adding stuff...")
		case "list":
			fmt.Println("Listing stuff...")
		case "delete":
			fmt.Println("Deleting stuff...")
		}
	*/
	return nil
}

func (p) ExecuteHookPre(hook plugin.Hook, args []string) error {
	switch hook.Name {
	default:
		return fmt.Errorf("hook not defined")
	}
}

func (p) ExecuteHookPost(hook plugin.Hook, args []string) error {
	switch hook.Name {
	default:
		return fmt.Errorf("hook not defined")
	}
}

func (p) ExecuteHookCleanUp(hook plugin.Hook, args []string) error {
	switch hook.Name {
	default:
		return fmt.Errorf("hook not defined")
	}
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
		"prout": &plugin.InterfacePlugin{Impl: &p{}},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
