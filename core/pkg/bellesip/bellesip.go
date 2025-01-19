package bellesip

// #cgo CFLAGS: -I/usr/include/bellesip
// #cgo LDFLAGS: -lbelle-sip
// #include <belle-sip/belle-sip.h>
import "C"

type ObjectPool C.belle_sip_object_pool_t

func ObjectPoolPush() *ObjectPool {
	return (*ObjectPool)(C.belle_sip_object_pool_push())
}
