package workflow

type Pipeline struct {
	Trigger string `yaml:"trigger"`
	Steps   []Step `yaml:"steps"`
}
