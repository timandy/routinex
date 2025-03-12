package api

type CmdOptions interface {
	IsDebug() bool         // is debug mode enabled or not
	IsVerbose() bool       // is verbose mode enabled or not
	GetArgs() []string     // get remain args exclude the options of current program
	SetArgs(args []string) // set remain args
	GetPackage() string    // get package name
	GetWorkDir() string    // get work dir
}
