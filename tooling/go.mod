module github.com/microsoft/yardl/tooling

go 1.21

require (
	github.com/alecthomas/participle/v2 v2.1.1
	github.com/dlclark/regexp2 v1.11.4
	github.com/fsnotify/fsnotify v1.8.0
	github.com/inancgumus/screen v0.0.0-20190314163918-06e984b86ed3
	github.com/knadh/koanf/providers/structs v0.1.0
	github.com/knadh/koanf/v2 v2.1.2
	github.com/rs/zerolog v1.33.0
	github.com/spf13/cobra v1.8.1
	github.com/stretchr/testify v1.9.0
	gopkg.in/yaml.v3 v3.0.1
)

require gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.15.0 // indirect
)

// Replace go-yaml with a fork that contains pending PR
// https://github.com/go-yaml/yaml/pull/691
replace gopkg.in/yaml.v3 => github.com/johnstairs/go-yaml-yaml v0.0.0-20221109150101-483fca0d3ee9
