package logs

type IoWriterFunc func(string, ...any)

func (f IoWriterFunc) Write(p []byte) (n int, err error) {
	f(string(p))
	return len(p), nil
}
