FROM golang:trixie AS builder

ARG PJSIP_VERSION=2.15.1
ARG DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y \
    git \
    pkg-config \
    build-essential \
    libssl-dev \
    libopenh264-dev \
    libvpx-dev \
    libv4l-dev \
    libopus-dev \
    libbcg729-dev \
    libsdl2-dev \
    libx11-dev \
    libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer-plugins-bad1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-alsa \
    gstreamer1.0-gl \
    gstreamer1.0-gtk3 \
    gstreamer1.0-qt5 \
    gstreamer1.0-pulseaudio

# Build pjsip
WORKDIR /build
RUN git clone --branch ${PJSIP_VERSION} --depth 1 https://github.com/pjsip/pjproject.git
COPY overlays/config_site.h pjproject/pjlib/include/pj/config_site.h
WORKDIR /build/pjproject
RUN ./configure CFLAGS="-fPIC" && \
    make && \
    make dep && \
    make install

# Build doorpix 
WORKDIR /build/doorpix
COPY . .
RUN go build -o /usr/local/bin/doorpix cmd/doorpix/main.go



FROM debian:trixie-slim AS app

RUN apt update && apt install --no-install-recommends -y \
    pipewire-alsa \ 
    libbcg729-dev \
    libgstreamer1.0-dev \
    libgstreamer-plugins-base1.0-dev \
    libgstreamer-plugins-bad1.0-dev \
    gstreamer1.0-plugins-base \
    gstreamer1.0-plugins-good \
    gstreamer1.0-plugins-bad \
    gstreamer1.0-plugins-ugly \
    gstreamer1.0-libav \
    gstreamer1.0-tools \
    gstreamer1.0-x \
    gstreamer1.0-alsa \
    gstreamer1.0-gl \
    gstreamer1.0-gtk3 \
    gstreamer1.0-qt5 \
    gstreamer1.0-pulseaudio

COPY --from=builder /usr/local/bin/doorpix /usr/local/bin/doorpix
RUN groupadd -g 1000 doorpix && \
    useradd -u 1000 -g doorpix -s /bin/sh -m doorpix
WORKDIR /var/doorpix
USER doorpix
EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/doorpix" ]
