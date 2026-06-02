#! /usr/bin/env bash

set -euo pipefail

# The .pixi directory is a docker volume that is mounted into the workspace.
# Ensure the current user owns it (it is created as root by the Docker daemon).
workspace_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${workspace_dir}"
sudo chown "$(id -u):$(id -g)" .pixi

# Create the pixi "dev" environment from pixi.toml / pixi.lock. This lives in the
# mounted .pixi volume rather than being baked into the image.
pixi install --locked --environment dev

# Create a stable gcov path for VS Code settings regardless of architecture.
gcov_symlink="${workspace_dir}/.pixi/envs/dev/bin/gcov"
gcov_candidates=("${workspace_dir}"/.pixi/envs/dev/bin/*-conda-linux-gnu-gcov)
if [[ -e "${gcov_candidates[0]}" ]]; then
    ln -sf "$(basename "${gcov_candidates[0]}")" "${gcov_symlink}"
fi

# Create stable compiler symlinks used by VS Code C/C++ tooling.
gxx_symlink="${workspace_dir}/.pixi/envs/dev/bin/g++"
gxx_candidates=("${workspace_dir}"/.pixi/envs/dev/bin/*-conda-linux-gnu-g++)
if [[ -e "${gxx_candidates[0]}" ]]; then
    ln -sf "$(basename "${gxx_candidates[0]}")" "${gxx_symlink}"
fi

gcc_symlink="${workspace_dir}/.pixi/envs/dev/bin/gcc"
gcc_candidates=("${workspace_dir}"/.pixi/envs/dev/bin/*-conda-linux-gnu-gcc)
if [[ -e "${gcc_candidates[0]}" ]]; then
    ln -sf "$(basename "${gcc_candidates[0]}")" "${gcc_symlink}"
fi

# Create a kits file for the VSCode CMake Tools extension, so you are not
# prompted for which kit to select whenever you open VSCode. The compiler paths
# come from the conda compiler packages, which set $GCC and $GXX when the pixi
# environment is activated.
kits_file="${HOME}/.local/share/CMakeTools/cmake-tools-kits.json"
mkdir -p "$(dirname "${kits_file}")"

# $GCC and $GXX are expanded by the inner shell (inside the activated pixi
# environment), so the body is single-quoted and the path is passed as an arg.
# shellcheck disable=SC2016
pixi run --environment dev bash -c \
    'echo "[{\"name\":\"Pixi\",\"compilers\":{\"C\":\"${GCC}\",\"CXX\":\"${GXX}\"}}]" > "$1"' \
    bash "${kits_file}"
