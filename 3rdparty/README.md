# 3rdparty

## linphone-sdk

```sh
apt-get update && apt-get install cmake python3 pip yasm nasm doxygen python3-pystache python3-six ninja-build meson libv4l-dev libglew-dev libpulse-dev libxext-dev
cd 3rdparty/linphone-sdk
rm -rf build
cmake --preset "default" -B "build" -G "Ninja Multi-Config" -DENABLE_DB_STORAGE=OFF -DENABLE_VCARD=OFF -DENABLE_NON_FREE_FEATURES=ON -DENABLE_OPENH264=ON
cmake --build "build" --config "Release" --parallel 8
sudo cmake --install "build" --prefix "/usr"
````
