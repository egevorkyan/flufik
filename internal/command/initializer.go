package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	//genFlufikPgpCommand := NewFlufikPgpGenCommand()
	pushFlufikCommand := NewFlufikPushCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
		pushFlufikCommand.command,
		//genFlufikPgpCommand.command,
	)
	return rootFlufikCommand
}
