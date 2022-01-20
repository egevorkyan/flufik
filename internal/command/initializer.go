package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	menuFlufikPgpCommand := NewFlufikPgp()
	pushFlufikCommand := NewFlufikPushCommand()
	serviceFlufikCommand := NewFlufikServiceCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
		pushFlufikCommand.command,
		menuFlufikPgpCommand.commnad,
		serviceFlufikCommand.command,
	)
	return rootFlufikCommand
}
