package flags

type Accumulator struct {
	Login        string
	Password     string
	Output       string
	Subscription string
	Environment  string
	Package      string
	Deploy       string
}

var Acc Accumulator
