package main

import (
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	hplugin "github.com/hashicorp/go-plugin"

	ignitecmd "github.com/ignite/cli/ignite/cmd"
	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/services/chain"
	"github.com/ignite/cli/ignite/services/plugin"
)

func init() {
	gob.Register(plugin.Command{})
}

type p struct {
	logger hclog.Logger
}

func (p) Commands() ([]plugin.Command, error) {
	cmd := plugin.Command{
		Use:               "band oracle-name",
		Short:             "Scaffold an IBC BandChain query oracle to request real-time data",
		Long:              "Scaffold an IBC BandChain query oracle to request real-time data from BandChain scripts in a specific IBC-enabled Cosmos SDK module",
		PlaceCommandUnder: "ignite scaffold",
	}
	cmd.Flags().StringP("path", "p", ".", "path of the app")
	cmd.Flags().String("signer", "", "Label for the message signer (default: creator)")
	cmd.Flags().String("module", "", "IBC Module to add the packet info")
	return []plugin.Command{cmd}, nil
}

func (p *p) Execute(cmd plugin.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("require the oracle query name as unique arg")
	}
	var (
		oracle     = args[0]
		appPath, _ = cmd.Flags().GetString("path")
		signer, _  = cmd.Flags().GetString("signer")
		module, _  = cmd.Flags().GetString("module")
	)
	if module == "" {
		return errors.New("please specify a module to create the BandChain oracle into: --module <module_name>")
	}

	session := cliui.New(cliui.StartSpinnerWithText("Scaffolding..."))
	defer session.End()

	var options []OracleOption
	if signer != "" {
		options = append(options, OracleWithSigner(signer))
	}

	sm, err := AddOracle(context.Background(), appPath, placeholder.New(), module, oracle, options...)
	if err != nil {
		return err
	}
	_ = sm

	modificationsStr, err := ignitecmd.SourceModificationToString(sm)
	if err != nil {
		return err
	}

	fmt.Println(modificationsStr)

	fmt.Printf(`
				🎉 Created a Band oracle query "%[1]v".

				Note: BandChain module uses version "bandchain-1".
				Make sure to update the keys.go file accordingly.

				// x/%[2]v/types/keys.go
				const Version = "bandchain-1"

				`, oracle, module)

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
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	p := &p{logger: logger}
	pluginMap := map[string]hplugin.Plugin{
		"band": &plugin.InterfacePlugin{Impl: p},
	}

	hplugin.Serve(&hplugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig(),
		Plugins:         pluginMap,
	})
}
