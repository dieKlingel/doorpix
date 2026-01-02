package pjsua2

import (
	"github.com/dieklingel/go-pjproject/pjsua2"
)

func init() {
	config := pjsua2.NewEpConfig()
	endpoint := pjsua2.NewEndpoint()
	endpoint.LibCreate()
	endpoint.LibInit(config)
	endpoint.LibStart()
}
