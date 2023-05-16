#########################################################
# file-normalizer stage
# In order to use BuildKit remote caching, input files must have
# not only the right content hash, but also the right permissions.
# Git only tracks whether the owner can execute a file.
# Here we bring in all files that are going to be used in the
# subsequent stage and normalize the permissions.
#########################################################

FROM mcr.microsoft.com/oss/busybox/busybox:1.33.1 as file-normalizer

COPY environment.yml \
     .devcontainer/devcontainer.bashrc  \
     /data/

RUN chmod -R 555 /data/

#########################################################
# devcontainer stage
# Installs all dependencies and tooling for development.
#########################################################

FROM mcr.microsoft.com/vscode/devcontainers/base:0.201.8-focal AS devcontainer

# Install needed packages and setup non-root user.
ARG USERNAME="vscode"
ARG USER_UID=1000
ARG USER_GID=$USER_UID
ARG CONDA_GID=900
ARG CONDA_ENVIRONMENT_NAME=yardl
ARG VSCODE_DEV_CONTAINERS_SCRIPT_LIBRARY_VERSION=v0.229.0

RUN apt-get update \
    && apt-get install -y \
        libc6-dbg \
        wget \
        bzip2 \
        ca-certificates \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Use the OSS Moby CLI instead of the licensed Docker CLI
ARG USE_MOBY="false"
ARG DOCKER_VERSION="20.10"
RUN script=$(curl -fsSL "https://raw.githubusercontent.com/devcontainers/features/8d3685e09f18dd8b0a6bce50abe3e868dac27a69/src/docker-outside-of-docker/install.sh") \
    && MOBY=${USE_MOBY} VERSION=${DOCKER_VERSION} bash -c "$script"

# Setting the ENTRYPOINT to docker-init.sh will configure non-root access to
# the Docker socket if "overrideCommand": false is set in devcontainer.json.
# The script will also execute CMD if you need to alter startup behaviors.
ENTRYPOINT [ "/usr/local/share/docker-init.sh" ]
CMD [ "sleep", "infinity" ]

ARG MAMBAFORGE_VERSION=22.9.0-2

# Based on https://github.com/conda-forge/miniforge-images/blob/master/ubuntu/Dockerfile
RUN wget --no-hsts --quiet https://github.com/conda-forge/miniforge/releases/download/${MAMBAFORGE_VERSION}/Mambaforge-${MAMBAFORGE_VERSION}-Linux-$(uname -m).sh -O /tmp/miniforge.sh \
    && /bin/bash /tmp/miniforge.sh -b -p /opt/conda \
    && rm /tmp/miniforge.sh \
    && /opt/conda/bin/conda clean --tarballs --index-cache --packages --yes \
    && find /opt/conda -follow -type f -name '*.a' -delete \
    && find /opt/conda -follow -type f -name '*.pyc' -delete \
    && /opt/conda/bin/conda clean --force-pkgs-dirs --all --yes  \
    && groupadd -r conda --gid ${CONDA_GID} \
    && usermod -aG conda ${USERNAME} \
    && chown -R :conda /opt/conda \
    && chmod -R g+w /opt/conda \
    && find /opt -type d | xargs -n 1 chmod g+s

# Create a conda environment from the environment file in the repo root.
COPY --from=file-normalizer --chown=$USER_UID:conda /data/environment.yml /tmp/build/
RUN umask 0002 \
    && /opt/conda/bin/mamba env create -f /tmp/build/environment.yml \
    && /opt/conda/bin/mamba clean -fy \
    && sudo chown -R :conda /opt/conda/envs

# Install Go
ARG GO_VERSION=1.20.2
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV "PATH"="/usr/local/go/bin:/go/bin:${PATH}"
RUN umask 0002 \
    && script=$(curl -fsSL "https://raw.githubusercontent.com/devcontainers/features/8d3685e09f18dd8b0a6bce50abe3e868dac27a69/src/go/install.sh") \
    && VERSION=${GO_VERSION} TARGET_GOPATH=${GOPATH} TARGET_GOROOT=${GOROOT} bash -c "$script" \
    && chown -R "${USERNAME}:conda" "${GOROOT}" "${GOPATH}"

# Install watchexec
ARG WATCHEXEC_VERSION=1.20.4
RUN wget --quiet https://github.com/watchexec/watchexec/releases/download/cli-v${WATCHEXEC_VERSION}/watchexec-${WATCHEXEC_VERSION}-x86_64-unknown-linux-gnu.deb -O watchexec.deb \
    && dpkg -i watchexec.deb

# Add a file that is to be sourced from .bashrc and from the devops pipeline stages
COPY --from=file-normalizer /data/devcontainer.bashrc /opt/devcontainer/

# Add a section to /etc/bash.bashrc that ensures that a section is present at the end of ~/.bashrc.
# We can't just write to .bashrc from here because it will be overwritten if the devcontainer user has
# opted to use their own dotfiles repo. The dotfiles repo is cloned after the postCreateCommand
# in the devcontainer.json file is executed.
RUN echo "\n\
if ! grep -q \"^source /opt/devcontainer/devcontainer.bashrc\" \${HOME}/.bashrc; then\n\
	echo \"source /opt/devcontainer/devcontainer.bashrc\" >> \${HOME}/.bashrc\n\
fi\n" >> /etc/bash.bashrc

ENV CMAKE_GENERATOR=Ninja

# Create a kits file for the VSCode CMake Tools extension, so you are not prompted for which kit to select whenever you open VSCode
RUN mkdir -p /home/vscode/.local/share/CMakeTools \
    && echo '[{"name":"GCC-10","compilers":{"C":"/opt/conda/envs/yardl/bin/x86_64-conda_cos6-linux-gnu-gcc","CXX":"/opt/conda/envs/yardl/bin/x86_64-conda_cos6-linux-gnu-g++"}}]' > /home/vscode/.local/share/CMakeTools/cmake-tools-kits.json \
    && chown vscode:conda /home/vscode/.local/share/CMakeTools/cmake-tools-kits.json
