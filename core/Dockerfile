FROM golang:bookworm AS builder

COPY . /app

RUN apt update && \
	apt-get install -y build-essential libssl-dev uuid-dev libasound2-dev libv4l-dev libvpx-dev libyuv-dev libopenh264-dev libsdl2-dev libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio libavutil-dev libavformat-dev libavcodec-dev libswscale-dev

WORKDIR /3rdparty
RUN git clone https://github.com/pjsip/pjproject.git && \
	cp /app/pj_config_site.h.in pjproject/pjlib/include/pj/config_site.h && \
	cd pjproject && \
	./configure && \
	make dep && \
	make && \
	make install

WORKDIR /app
RUN go build -o core cmd/main.go



FROM debian:bookworm

RUN apt-get update && \
	apt-get install -y libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio	

COPY --from=builder /app/core /usr/local/bin/core

ENTRYPOINT [ "/usr/local/bin/core" ]
