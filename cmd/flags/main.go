package flags

type Accumulator struct {
	AutoAccept   bool
	Login        string
	Password     string
	Output       string
	Subscription string
	Environment  string
	Package      string
	Deploy       string
	Exclude      string
}

var Acc Accumulator
