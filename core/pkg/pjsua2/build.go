package pjsua2

//go:generate swig -c++ -go -cgo -intgosize 64 -outcurrentdir -I../../../3rdparty/pjproject/pjsip/include -I../../../3rdparty/pjproject/pjlib/include ../../../3rdparty/pjproject/pjsip-apps/src/swig/pjsua2.i

// #cgo CPPFLAGS: -DPJ_AUTOCONF=1 -O2 -DPJ_IS_BIG_ENDIAN=0 -DPJ_IS_LITTLE_ENDIAN=1
// #cgo CPPFLAGS: -I${SRCDIR}/../../../3rdparty/pjproject/pjsip/include -I${SRCDIR}/../../../3rdparty/pjproject/pjlib/include -I${SRCDIR}/../../../3rdparty/pjproject/pjlib-util/include -I${SRCDIR}/../../../3rdparty/pjproject/pjmedia/include -I${SRCDIR}/../../../3rdparty/pjproject/pjnath/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/pjsip/lib -lpjsua2-aarch64-unknown-linux-gnu -lpjsua-aarch64-unknown-linux-gnu -lpjsip-ua-aarch64-unknown-linux-gnu -lpjsip-simple-aarch64-unknown-linux-gnu -lpjsip-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/pjmedia/lib  -lpjmedia-codec-aarch64-unknown-linux-gnu -lpjmedia-videodev-aarch64-unknown-linux-gnu -lpjmedia-audiodev-aarch64-unknown-linux-gnu -lpjsdp-aarch64-unknown-linux-gnu -lpjmedia-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/pjnath/lib -lpjnath-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/pjlib-util/lib -lpjlib-util-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/pjlib/lib -lpj-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -L${SRCDIR}/../../../3rdparty/pjproject/third_party/lib -lresample-aarch64-unknown-linux-gnu -lspeex-aarch64-unknown-linux-gnu -lgsmcodec-aarch64-unknown-linux-gnu -lsrtp-aarch64-unknown-linux-gnu -lilbccodec-aarch64-unknown-linux-gnu -lg7221codec-aarch64-unknown-linux-gnu -lwebrtc-aarch64-unknown-linux-gnu -lyuv-aarch64-unknown-linux-gnu
// #cgo LDFLAGS: -lssl -lcrypto -lm -lpthread -luuid -lrt -lasound
// #cgo LDFLAGS: -lSDL2 -lv4l2 -lvpx -lopenh264
import "C"
