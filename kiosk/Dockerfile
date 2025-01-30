FROM debian:bookworm as builder

RUN apt update && \
	apt install -y git unzip curl clang cmake pkg-config libgles2-mesa-dev libxkbcommon-x11-0 libwayland-dev libglib2.0-dev libgstreamer1.0-dev

WORKDIR /opt
ENV PATH="${PATH}:/opt/flutter-elinux/bin"
RUN git clone https://github.com/sony/flutter-elinux.git && \
	flutter-elinux precache

WORKDIR /app
COPY . /app
RUN flutter-elinux pub get && \
	flutter-elinux build elinux



FROM debian:bookworm

RUN apt update && \
	apt install -y libfontconfig libxkbcommon-x11-0 libwayland-client0 libwayland-cursor0 libwayland-egl1-mesa libgl1-mesa-glx libgl1-mesa-dri libglib2.0-dev pipewire-alsa libgstreamer1.0-dev libgstreamer-plugins-base1.0-dev libgstreamer-plugins-bad1.0-dev gstreamer1.0-plugins-base gstreamer1.0-plugins-good gstreamer1.0-plugins-bad gstreamer1.0-plugins-ugly gstreamer1.0-libav gstreamer1.0-tools gstreamer1.0-x gstreamer1.0-alsa gstreamer1.0-gl gstreamer1.0-gtk3 gstreamer1.0-qt5 gstreamer1.0-pulseaudio

COPY --from=builder /app/build/elinux/*/release/bundle /app
ENTRYPOINT [ "/app/kiosk", "--bundle=.", "--fullscreen" ]
