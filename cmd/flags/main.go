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
	Strategy          string
	Backup            string
	Deploy            string
	Exclude           string
	Type              string
	Ref               string
	SrcEnvironment    string
	DestEnvironment   string
	UseLatestBackup   bool
}

var Acc Accumulator
