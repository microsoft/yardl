#! /bin/bash
# shellcheck source=/dev/null

# Activate the pixi "dev" environment for interactive shells.
parent_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)"
eval "$(pixi shell-hook --manifest-path "${parent_dir}/pixi.toml" --environment dev)"

# Shell completions
source <(pixi completion --shell bash)
source <(just --completions bash)

PATH=${PATH}:${HOME}/go/bin

if [[ "${BASH_ENV:-}" == "$(readlink -f "${BASH_SOURCE[0]:-}")" ]]; then
    # We don't want subshells to unnecessarily source this again.
    unset BASH_ENV
fi
