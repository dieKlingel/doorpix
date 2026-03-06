# Raspberry Pi

As the name suggests, _doorpix_ is build to be run on a Raspberry Pi. The Raspberry Pi is a pretty good fit, because it has a great form factor and we can use the gpio pins, to directly attach buttons to it.

## Install Docker

As first step you need to install docker (or any other container runtime of your choice), as doorpix will be run in a container. See https://docs.docker.com/engine/install/debian/ for further help.

## Runn DoorPiX

### Audio

Make sure you have installed pipewire and it is up and running. You can do so by verifying that a the pipewire socket exists:

```sh
ls -l $XDG_RUNTIME_DIR/pipewire-0
# output: srw-rw-rw- 1 pi pi 0 Mar  6 18:07 /run/user/1000/pipewire-0
```

If you have not installed pipewire you can install and start it with

```sh
sudo apt install -y pipewire wireplumber libspa-0.2-bluetooth
systemctl --user restart pipewire.service
```

### Video

TBD: document usage with raspicam/libcamera

### DoorPiX

You can use the docker compose file, to run doorpix.

```sh
docker compose up
```
