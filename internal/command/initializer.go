package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	menuFlufikPgpCommand := NewFlufikPgp()
	pushFlufikCommand := NewFlufikPushCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
		pushFlufikCommand.command,
		menuFlufikPgpCommand.commnad,
	)
	return rootFlufikCommand
}
