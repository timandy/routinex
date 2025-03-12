package opt

type AppOptions struct {
	Debug   bool `name:"debug" shorthand:"d" usage:"enable debug mode for detailed output."`
	Verbose bool `name:"verbose" shorthand:"v" usage:"enable verbose mode for additional information."`
	Help    bool `name:"help" shorthand:"h" usage:"display help information for this command."`
}
