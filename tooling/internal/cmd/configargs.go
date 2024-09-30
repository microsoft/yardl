package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/rs/zerolog/log"
)

var k = koanf.New(".")

// ConfigArgs is a provider that parses configuration key=value pairs from command line arguments.
// It is used to override package information in packaging.PackageInfo.
// koanf provides a posflag provider that works with spf13/pflag (used by Cobra), but it requires
// use to define every packaging.PackageInfo field as a command-line flag.
type ConfigArgs struct {
	delim string
	args  []string
	ko    *koanf.Koanf
}

func ConfigArgsProvider(args []string, delim string, ko *koanf.Koanf) *ConfigArgs {
	return &ConfigArgs{
		delim: delim,
		args:  args,
		ko:    ko,
	}
}

func (aa *ConfigArgs) Read() (map[string]interface{}, error) {
	mp := make(map[string]interface{})

	for _, arg := range aa.args {
		parts := strings.Split(arg, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid config 'key=value' pair: %s", arg)
		}

		key := parts[0]
		value := parts[1]

		if key == "" || value == "" {
			return nil, fmt.Errorf("invalid config 'key=value' pair: %s", arg)
		}

		if aa.ko.Exists(key) {
			log.Debug().Msgf("Overriding config key %s with value %s", key, value)
			aa.ko.Set(key, value)
		} else {
			log.Debug().Msgf("Ignoring invalid config key %s", key)
		}

		/* We don't need to set the value in the map because we've already
		added it directly to the koanf instance IFF the key is a valid
		field in packaging.PackageInfo. */
		// mp[key] = value
	}

	return maps.Unflatten(mp, aa.delim), nil
}

func (*ConfigArgs) ReadBytes() ([]byte, error) {
	return nil, errors.New("method unsupported")
}

func (*ConfigArgs) Watch(cb func(event interface{}, err error)) error {
	return errors.New("method unsupported")
}

// updatePackageInfoFromArgs overrides the fields in packageInfo using command-line arguments
func updatePackageInfoFromArgs(packageInfo *packaging.PackageInfo, configArgs []string) error {
	if err := k.Load(structs.Provider(packageInfo, "yaml"), nil); err != nil {
		log.Panic().Msgf("error loading package info: %v", err)
	}

	if err := k.Load(ConfigArgsProvider(configArgs, ".", k), nil); err != nil {
		return err
	}

	if err := k.UnmarshalWithConf("", packageInfo, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return fmt.Errorf("error overriding package info: %w", err)
	}

	return nil
}
