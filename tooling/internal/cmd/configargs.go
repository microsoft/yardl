package cmd

import (
	"fmt"

	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"github.com/microsoft/yardl/tooling/pkg/packaging"
	"github.com/rs/zerolog/log"
)

var k = koanf.New(".")

// updatePackageInfoFromArgs overrides the fields in packageInfo using command-line arguments
func updatePackageInfoFromArgs(packageInfo *packaging.PackageInfo, configArgs map[string]string) error {
	if err := k.Load(structs.Provider(packageInfo, "yaml"), nil); err != nil {
		log.Panic().Msgf("error loading package info: %v", err)
	}

	for key, value := range configArgs {
		if k.Exists(key) {
			log.Info().Msgf("Overriding config key %s with value %s", key, value)
			k.Set(key, value)
		} else {
			return fmt.Errorf("invalid config key %s", key)
		}
	}

	if err := k.UnmarshalWithConf("", packageInfo, koanf.UnmarshalConf{Tag: "yaml"}); err != nil {
		return fmt.Errorf("error overriding package info: %w", err)
	}

	return nil
}
