package linphone

// #cgo linux CFLAGS: -I/usr/include/linphone
// #cgo linux LDFLAGS: -llinphone -lmediastreamer2 -lortp -lbctoolbox
import "C"
