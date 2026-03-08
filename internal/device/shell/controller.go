package shell

import "os/exec"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Exec(name string, args ...string) ([]byte, error) {
	return exec.Command(name, args...).CombinedOutput()
}
