package linphone

import "C"

func boolToCBool(b bool) C.uchar {
	if b {
		return 1
	}
	return 0
}
