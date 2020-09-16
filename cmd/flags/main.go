package flags

type Accumulator struct {
	AutoAccept        bool
	ApproveProdDeploy bool
	Login             string
	Password          string
	Output            string
	Subscription      string
	Environment       string
	Name              string
	Retention         string
	Component         []string
	Package           string
	Deploy            string
	Exclude           string
	Type              string
	Ref               string
}

var Acc Accumulator
