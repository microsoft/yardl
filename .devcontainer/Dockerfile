ARG DEVCONTAINER_BASE=mcr.microsoft.com/devcontainers/base:1.0.9-ubuntu-22.04

#########################################################
# devcontainer stage
# Installs all dependencies and tooling for development.
#########################################################

FROM ${DEVCONTAINER_BASE} AS devcontainer

#########################################################
# Install Matlab (conditional on INSTALL_MATLAB build arg)
#   Based on mathworks/{matlab-deps:r2023b,matlab:r2023b}
#
# Configure the Matlab License Server to use Matlab within the devcontainer (including VSCode extensions):
#   In your HOST environment, export the environment variable MATLAB_LICENSE_FILE, e.g.
#       export MATLAB_LICENSE_FILE=/mnt/c/Users/username/Documents/MATLAB/license.lic
#
# To skip MATLAB installation, set the INSTALL_MATLAB build arg to "false" in devcontainer.json.
#########################################################

ENV DEBIAN_FRONTEND="noninteractive" TZ="Etc/UTC"

ARG INSTALL_MATLAB=true
ARG MATLAB_RELEASE=r2023b
ARG MATLAB_PRODUCT_LIST="MATLAB"
ARG MATLAB_INSTALL_LOCATION="/opt/matlab/${MATLAB_RELEASE}"
ARG LICENSE_SERVER=

ENV MLM_LICENSE_FILE=${LICENSE_SERVER}

COPY .devcontainer/matlab-r2023b.ubuntu-22.04.dependencies.txt /tmp/matlab-dependencies.txt

RUN if [ "${INSTALL_MATLAB}" = "true" ]; then \
        set -e; \
        export DEBIAN_FRONTEND=noninteractive; \
        apt-get update; \
        apt-get install --no-install-recommends --yes \
            $(cat /tmp/matlab-dependencies.txt) \
            wget \
            unzip \
            ca-certificates; \
        apt-get clean; \
        apt-get -y autoremove; \
        rm -rf /var/lib/apt/lists/*; \
        [ -d /usr/share/X11/xkb ] || mkdir -p /usr/share/X11/xkb; \
        adduser --shell /bin/bash --disabled-password --gecos "" matlab; \
        echo "matlab ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/matlab; \
        chmod 0440 /etc/sudoers.d/matlab; \
        cd /home/matlab; \
        wget -q https://www.mathworks.com/mpm/glnxa64/mpm -O /home/matlab/mpm; \
        chmod +x /home/matlab/mpm; \
        HOME=/home/matlab /home/matlab/mpm install \
            --release=${MATLAB_RELEASE} \
            --destination=${MATLAB_INSTALL_LOCATION} \
            --products ${MATLAB_PRODUCT_LIST} \
            || (echo "MPM Installation Failure. See below for more information:" && cat /tmp/mathworks_root.log && false); \
        rm -f /home/matlab/mpm /tmp/mathworks_root.log; \
        ln -s ${MATLAB_INSTALL_LOCATION}/bin/matlab /usr/local/bin/matlab; \
        ln -s ${MATLAB_INSTALL_LOCATION} /opt/matlab/latest; \
        mkdir -p ${MATLAB_INSTALL_LOCATION}/licenses; \
        chown -R matlab:matlab /opt/matlab; \
        wget -O /usr/local/bin/run-matlab-command https://ssd.mathworks.com/supportfiles/ci/run-matlab-command/v2/glnxa64/run-matlab-command; \
        chmod +x /usr/local/bin/run-matlab-command; \
    else \
        echo "Skipping MATLAB installation (INSTALL_MATLAB=${INSTALL_MATLAB})"; \
    fi


#########################################################
# Install everything else
#########################################################


# Install needed packages and setup non-root user.
ARG USERNAME="vscode"
ARG USER_UID=1000
ARG USER_GID=$USER_UID

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
ARG DOCKER_VERSION="24.0.6"
ARG DOCKER_DASH_COMPOSE_VERSION="none"
RUN script=$(curl -fsSL "https://raw.githubusercontent.com/devcontainers/features/4fca96b5e8a4bfc93679098cb19d73c65ce571eb/src/docker-outside-of-docker/install.sh") \
    && MOBY=${USE_MOBY} VERSION=${DOCKER_VERSION} DOCKERDASHCOMPOSEVERSION=${DOCKER_DASH_COMPOSE_VERSION} bash -c "$script"

# Setting the ENTRYPOINT to docker-init.sh will configure non-root access to
# the Docker socket if "overrideCommand": false is set in devcontainer.json.
# The script will also execute CMD if you need to alter startup behaviors.
ENTRYPOINT [ "/usr/local/share/docker-init.sh" ]
CMD [ "sleep", "infinity" ]

# Install pixi. The conda/pypi environment itself is created from pixi.toml /
# pixi.lock at container creation time (see postCreateCommand), so that it can
# live in the mounted .pixi volume rather than being baked into the image.
ARG PIXI_VERSION=0.70.0
RUN mkdir -p /home/${USERNAME}/.pixi/bin \
    && curl -L -o /home/${USERNAME}/.pixi/bin/pixi -fsSL --compressed "https://github.com/prefix-dev/pixi/releases/download/v${PIXI_VERSION}/pixi-$(uname -m)-unknown-linux-musl" \
    && chmod +x /home/${USERNAME}/.pixi/bin/pixi \
    && chown -R ${USERNAME}:${USERNAME} /home/${USERNAME}/.pixi
ENV PATH="/home/${USERNAME}/.pixi/bin:${PATH}"

# Install Go
COPY tooling/go.mod /tmp/build/go.mod
ENV GOROOT="/usr/local/go"
ENV GOPATH="/go"
ENV "PATH"="/usr/local/go/bin:/go/bin:${PATH}"
RUN umask 0002 \
    && GO_VERSION="$(awk '/^go /{print $2; exit}' /tmp/build/go.mod)" \
    && [ -n "${GO_VERSION}" ] \
    && script=$(curl -fsSL "https://raw.githubusercontent.com/devcontainers/features/72df8a5f191f840a66dc2e2ced10a136e4d75173/src/go/install.sh") \
    && VERSION=${GO_VERSION} TARGET_GOPATH=${GOPATH} TARGET_GOROOT=${GOROOT} bash -c "$script" \
    && chown -R "${USERNAME}:${USERNAME}" "${GOROOT}" "${GOPATH}"

# Install watchexec
ARG WATCHEXEC_VERSION=2.5.1
RUN wget --quiet https://github.com/watchexec/watchexec/releases/download/v${WATCHEXEC_VERSION}/watchexec-${WATCHEXEC_VERSION}-$(uname -m)-unknown-linux-musl.deb -O watchexec.deb \
    && dpkg -i watchexec.deb

ENV CMAKE_GENERATOR=Ninja
