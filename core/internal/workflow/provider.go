package workflow

type ProviderFunc func(ctx *Context) error

type Provider interface {
	Parse(step Step) (StepDelegate, error)
}

type StepDelegate struct {
	Step    Step
	Execute ProviderFunc
	Parent  any
}
