FROM mcr.microsoft.com/devcontainers/base:bookworm

LABEL org.opencontainers.image.authors="Kai Mayer <kai.mayer@dieklingel.com>"

### Install PJSIP Dependencies
RUN apt-get update && \
	apt-get install -y \
	build-essential \
	cmake \
	git \
	libssl-dev \
	uuid-dev \
	libasound2-dev \
	libpulse-dev \
	libssl-dev \
	libssl-dev \
	libxml2-dev \
	libsqlite3-dev \
	libv4l-dev \
	libvpx-dev \
	libyuv-dev \
	libopenh264-dev \
	libsdl2-dev \
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
	gstreamer1.0-pulseaudio \
	libavutil-dev \
	libavformat-dev \
	libavcodec-dev \
	libswscale-dev \
	libopus-dev \
	libvpx-dev \
	libspeexdsp-dev

# Install dev dependencies
RUN apt-get install -y \
	swig \
	protobuf-compiler

ARG PJ_PROJECT_SRC_DIR="pjproject"
COPY ${PJ_PROJECT_SRC_DIR} /3rdparty/pjproject
WORKDIR /3rdparty/pjproject
RUN ./configure CFLAGS="-DPJMEDIA_HAS_VIDEO=1 -DPJ_ICE_MAX_CAND=32 -DPJSIP_MAX_PKT_LEN=12288 -DPJ_ICE_MAX_CHECKS=512"  && \
	make dep && \
	make && \
	make install
