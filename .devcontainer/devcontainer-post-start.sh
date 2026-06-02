#! /usr/bin/env bash

set -euo pipefail

# Ensure the user's .bashrc sources the devcontainer.bashrc.
# This is done in postStartCommand (rather than the Dockerfile) so that it runs
# after any custom dotfiles have already been applied.

workspace_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

line="source ${workspace_dir}/.devcontainer/devcontainer.bashrc"

if ! grep -Fqx "${line}" "${HOME}/.bashrc" 2>/dev/null; then
    printf '%s\n' "${line}" >> "${HOME}/.bashrc"
fi
