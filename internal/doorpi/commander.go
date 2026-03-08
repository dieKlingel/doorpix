package doorpi

type Commander interface {
	Exec(name string, args ...string) ([]byte, error)
}
