package main

import (
	"context"
	"path/filepath"

	"github.com/gobuffalo/genny"

	"github.com/ignite/cli/ignite/pkg/gocmd"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/multiformatname"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/ibc"
)

const (
	bandImport  = "github.com/bandprotocol/bandchain-packet"
	bandVersion = "v0.0.2"
)

// OracleOption configures options for AddOracle.
type OracleOption func(*oracleOptions)

type oracleOptions struct {
	signer string
}

// newOracleOptions returns a oracleOptions with default options
func newOracleOptions() oracleOptions {
	return oracleOptions{
		signer: "creator",
	}
}

// OracleWithSigner provides a custom signer name for the message
func OracleWithSigner(signer string) OracleOption {
	return func(m *oracleOptions) {
		m.signer = signer
	}
}

// AddOracle adds a new BandChain oracle envtest.
func AddOracle(
	ctx context.Context,
	path string,
	tracer *placeholder.Tracer,
	moduleName,
	queryName string,
	options ...OracleOption,
) (sm xgenny.SourceModification, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return sm, err
	}

	if err := installBandPacket(ctx, path); err != nil {
		return sm, err
	}

	modpath, path, err := gomodulepath.Find(path)
	if err != nil {
		return sm, err
	}

	o := newOracleOptions()
	for _, apply := range options {
		apply(&o)
	}

	mfName, err := multiformatname.NewName(moduleName, multiformatname.NoNumber)
	if err != nil {
		return sm, err
	}
	moduleName = mfName.LowerCase

	name, err := multiformatname.NewName(queryName)
	if err != nil {
		return sm, err
	}

	/*
		if err := checkComponentValidity(path, moduleName, name, false); err != nil {
			return sm, err
		}

		mfSigner, err := multiformatname.NewName(o.signer, checkForbiddenOracleFieldName)
		if err != nil {
			return sm, err
		}

		// Module must implement IBC
		ok, err := isIBCModule(path, moduleName)
		if err != nil {
			return sm, err
		}
		if !ok {
			return sm, fmt.Errorf("the module %s doesn't implement IBC module interface", moduleName)
		}
	*/

	// Generate the packet
	var (
		g    *genny.Generator
		opts = &ibc.OracleOptions{
			AppName:    modpath.Package,
			AppPath:    path,
			ModulePath: modpath.RawPath,
			ModuleName: moduleName,
			QueryName:  name,
			// MsgSigner:  mfSigner,
		}
	)
	g, err = ibc.NewOracle(tracer, opts)
	if err != nil {
		return sm, err
	}
	sm, err = xgenny.RunWithValidation(tracer, g)
	if err != nil {
		return sm, err
	}
	if err := gocmd.ModTidy(ctx, path); err != nil {
		return sm, err
	}
	return sm, gocmd.Fmt(ctx, path)
}

func installBandPacket(ctx context.Context, path string) error {
	return gocmd.Get(ctx, path, []string{
		gocmd.PackageLiteral(bandImport, bandVersion),
	})
}
