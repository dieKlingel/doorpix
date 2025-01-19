package linphone

// #include <bctoolbox/list.h>
import "C"

type BctbxList C.bctbx_list_t

func bctbxListToSlice[T any](list *BctbxList) []*T {
	size := C.bctbx_list_size((*C.bctbx_list_t)(list))
	items := make([]*T, int(size))

	item := list
	for item != nil {
		items = append(items, (*T)(item.data))
		item = (*BctbxList)(item.next)
	}

	return items
}
