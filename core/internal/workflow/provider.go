package workflow

type ProviderFunc func(ctx *Context) error

type Provider interface {
	Parse(step Step) (StepDelegate[any], error)
}

type StepDelegate[T any] struct {
	Step    Step
	Execute ProviderFunc
	Parent  T
}
