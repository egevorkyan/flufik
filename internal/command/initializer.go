package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	genFlufikPgpCommand := NewFlufikPgpCommand()
	pushFlufikCommand := NewFlufikPushCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
		pushFlufikCommand.command,
		genFlufikPgpCommand.command,
	)
	return rootFlufikCommand
}
