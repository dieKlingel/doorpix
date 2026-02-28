package oplog

import "errors"

var ErrOplogWriteFailed = errors.New("oplog: failed to write oplog")

type Writer interface {
	Write(evt Event) error
}
