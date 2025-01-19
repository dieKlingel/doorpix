package main

import (
	"fmt"

	"github.com/dieklingel/doorpix/core/pkg/linphone"
)

func main() {
	fmt.Printf("doorpix, linphone version: %s", linphone.GetVersion())
}
