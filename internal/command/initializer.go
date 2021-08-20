package command

func InitialiseRootCmd() *RootFlufikCommand {
	rootFlufikCommand := NewFlufikRootCommand()
	buildFlufikCommand := NewFlufikBuildCommand()
	rootFlufikCommand.Command.AddCommand(
		buildFlufikCommand.command,
	)
	return rootFlufikCommand
}
