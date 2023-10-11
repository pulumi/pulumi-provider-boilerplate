# Example Usage:
# - docker build -t provider -f docker/Dockerfile .
# - docker run -it --rm -v $PWD:/provider provider

FROM docker.io/library/ubuntu:22.04

ARG PIP_PKGS="\
setuptools \
"

ARG APT_PKGS="\
git \
curl \
gnupg \
python3 \
python3-pip \
dotnet-sdk-7.0 \
ca-certificates \
build-essential \
dotnet-runtime-7.0 \
# golang-go \ # default packages go 1.18 is not new enough
"

# Append rootfs directory tree into container
# place any additional files to be included into the rootfs directory
ADD docker/rootfs /
# Disable Timezone prompts
ENV TZ=UTC
# Disable time zone prompts etc.
ENV DEBIAN_FRONTEND=noninteractive
# apt + pip packages
RUN set -ex \
    && apt-get update \
    && apt-get install ${APT_PKGS} \
    && update-alternatives --install \
        /usr/bin/python \
        python \
        /usr/bin/python3 1 \
    && python3 -m pip install ${PIP_PKGS} \
    && apt-get clean \
    && apt-get autoremove -y \
    && apt-get purge -y --auto-remove \
    && rm -rf \
        /var/lib/{apt,dpkg,cache,log} \
        /usr/share/{doc,man,locale} \
        /var/cache/apt \
        /root/.cache \
        /var/tmp/* \
        /tmp/* \
    && true

# golang
ENV PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/usr/games:/usr/local/games:/usr/local/go/bin"
RUN set -ex \
    && export arch=$(uname -m | awk '{print ($1 == "x86_64") ? "amd64" : ($1 == "aarch64" ? "arm64" : "unknown")}') \
    && export goversion="$(curl -s https://go.dev/dl/?mode=json \
        | awk -F'[":go]' '/  "version"/{print $8}' \
        | head -n1)" \
    && curl -L https://go.dev/dl/go${goversion}.linux-${arch}.tar.gz \
        | tar -C /usr/local/ -xzvf - \
    && which go \
    && go version \
    && true

# nodejs npm yarn
RUN set -ex \
    && export NODE_MAJOR=20 \
    && curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key \
        | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg \
    && echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" \
        | tee /etc/apt/sources.list.d/nodesource.list \
    && apt-get update \
    && apt-get install nodejs \
    && apt-get clean \
    && apt-get autoremove -y \
    && apt-get purge -y --auto-remove \
    && rm -rf \
        /var/lib/{apt,dpkg,cache,log} \
        /usr/share/{doc,man,locale} \
        /var/cache/apt \
        /root/.cache \
        /var/tmp/* \
        /tmp/* \
    && node --version \
    && npm --version \
    && npm install --global yarn \
    && yarn --version \
    && true

# pulumi
RUN set -ex \
    && export arch=$(uname -m | awk '{print ($1 == "x86_64") ? "x64" : ($1 == "aarch64" ? "arm64" : "unknown")}') \
    && export urlPulumiRelease="https://api.github.com/repos/pulumi/pulumi/releases/latest" \
    && export urlPulumiVersion=$(curl -s ${urlPulumiRelease} | awk -F '["v,]' '/tag_name/{print $5}') \
    && export urlPulumiBase="https://github.com/pulumi/pulumi/releases/download" \
    && export urlPulumiBin="pulumi-v${urlPulumiVersion}-linux-x64.tar.gz" \
    && export urlPulumi="${urlPulumiBase}/v${urlPulumiVersion}/${urlPulumiBin}" \
    && curl -L ${urlPulumi} \
        | tar xzvf - --directory /tmp \
    && mv /tmp/pulumi/* /usr/local/bin/ \
    && rm -rf /tmp/pulumi \
    && which pulumi \
    && pulumi version \
    && true

# pulumictl
RUN set -ex \
    && export arch=$(uname -m | awk '{print ($1 == "x86_64") ? "amd64" : ($1 == "aarch64" ? "arm64" : "unknown")}') \
    && export urlPulumiRelease="https://api.github.com/repos/pulumi/pulumictl/releases/latest" \
    && export urlPulumiVersion=$(curl -s ${urlPulumiRelease} | awk -F '["v,]' '/tag_name/{print $5}') \
    && export urlPulumiBase="https://github.com/pulumi/pulumictl/releases/download" \
    && export urlPulumiBin="pulumictl-v${urlPulumiVersion}-linux-arm64.tar.gz" \
    && export urlPulumi="${urlPulumiBase}/v${urlPulumiVersion}/${urlPulumiBin}" \
    && curl -L ${urlPulumi} \
        | tar xzvf - --directory /tmp \
    && mv /tmp/pulumictl /usr/local/bin/ \
    && rm -rf /tmp/* \
    && which pulumictl \
    && pulumictl version \
    && true

WORKDIR /provider
CMD ["make", "build"]

ARG VERSION
ARG BUILD_DATE
ARG PULUMICTL
ARG PULUMI
LABEL \
    org.opencontainers.image.version=$VERSION \
    org.opencontainers.image.created=$BUILD_DATE \
    org.opencontainers.image.licenses="APACHE2" \
    org.opencontainers.image.vendor="Pulumi" \
    org.opencontainers.image.title="Pulumi native provider build image." \
    org.opencontainers.image.description="A containerized environment for building native pulumi providers." \
    org.opencontainers.image.documentation="https://github.com/pulumi/pulumi-provider-boilerplate" \
    org.opencontainers.image.url="https://github.com/pulumi/pulumi-provider-boilerplate" \
    org.opencontainers.image.authors="https://github.com/pulumi"