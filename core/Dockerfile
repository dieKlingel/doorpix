FROM golang:bookworm AS builder

COPY . /app

RUN apt update && \
	apt-get install -y build-essential libssl-dev uuid-dev libasound2-dev libv4l-dev libvpx-dev libyuv-dev libopenh264-dev libsdl2-dev && \
	apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio && \
	apt install -y libavutil-dev libavformat-dev libavcodec-dev libswscale-dev

COPY config_site.h /tmp/config_site.h
WORKDIR /3rdparty
RUN git clone https://github.com/pjsip/pjproject.git && \
	cp /tmp/config_site.h pjproject/pjlib/include/pj/config_site.h && \
	cd pjproject && \
	./configure && \
	make dep && \
	make

WORKDIR /app
RUN go build -o core cmd/main.go
