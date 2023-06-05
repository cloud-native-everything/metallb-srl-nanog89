ARG version=1.22.2
ARG base=kindest/node
FROM ${base}:v${version}

LABEL org.opencontainers.image.authors="olivier.tilmans@nokia.com"
LABEL description="Disk image for kind workers able to run ENC/vSR"
LABEL version="$version"

# See https://kind.sigs.k8s.io/docs/design/node-image for a walkthrough of
# the kind base image. In particular, that its entrypoint enables to start
# a systemd instance, hence honoring the below systemctl enable commands.
RUN for i in $(seq 5); \
        do DEBIAN_FRONTEND=noninteractive apt update -qq && break || sleep 5; \
    done && \
    for i in $(seq 5); \
        do DEBIAN_FRONTEND=noninteractive apt install -qq -y \
            iputils-ping libnl-utils net-tools \
            tcpdump lldpd ssh && break || sleep 5; \
    done && \
    systemctl enable lldpd && \
    systemctl enable containerd && \
    bash -c 'curl -sL https://get-gnmic.kmrd.dev | bash'
