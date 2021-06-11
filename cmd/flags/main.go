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
	SavePath          string
	LockTables        bool
	SkipHooks         bool
	PreventRollback   bool
	Component         []string
	Package           string
	Strategy          string
	Key               string
	Value             string
	VarType           string
	Backup            string
	BackupType        string
	Limit             string
	Offset            string
	Start             int
	End               int
	LogStreams        []string
	Stream            bool
	Search            string
	Deploy            string
	Exclude           string
	Type              string
	Ref               string
	Tag               string
	Branch            string
	Checksum          string
	CommitSHA         string
	CustomPackage     string
	SrcEnvironment    string
	DestEnvironment   string
	UseLatestBackup   bool
	LockSessionToIP   bool
}

var Acc Accumulator
