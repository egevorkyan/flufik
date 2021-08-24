package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	//genFlufikPgpCommand := NewFlufikPgpGenCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
		//genFlufikPgpCommand.command,
	)
	return rootFlufikCommand
}
