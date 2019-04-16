package completion

import (
	"fmt"
	"strings"
)

func zsh(cmd *command) string {
	var zsh string
	zsh += fmt.Sprintf("#!/bin/zsh\n")
	zsh += zshFunctions(cmd)
	zsh += fmt.Sprintf("\ncompdef _%s %s\n", cmd.Name, cmd.Name)
	return zsh
}

func zshFunctions(cmd *command) string {
	var function string
	function += fmt.Sprintf("\n%s() {\n", zshFunctionsName(cmd))
	function += fmt.Sprintf("\tlocal line\n")
	if cmd.Name == "help" {
		function += zshHelpFunction(cmd)
	} else {
		function += zshSubCommandsInfo(cmd)
		function += zshArguments(cmd)
		function += zshSubCommandsLink(cmd)
	}
	function += fmt.Sprintf("}\n")
	if cmd.HasSubCommands() {
		for _, c := range cmd.SubCommands {
			function += zshFunctions(c)
		}
	}
	return function
}

func zshHelpFunction(cmd *command) string {
	var helpFunction string
	helpFunction += fmt.Sprintf("\n\t_arguments \\\n")
	helpFunction += fmt.Sprintf("\t\t\"1:command:(%s)\"\n", zshListAllCommands(cmd.Parent))
	return helpFunction
}

func zshFunctionsName(cmd *command) string {
	var parentName string
	if cmd.HasParent() {
		var parent = cmd.Parent
		for {
			parentName = fmt.Sprintf("_%s%s", parent.Name, parentName)
			if !parent.HasParent() {
				break
			}
			parent = parent.Parent
		}
	}
	return fmt.Sprintf("%s_%s", parentName, cmd.Name)
}

func zshListAllCommands(cmd *command) string {

	var list string
	if cmd.HasSubCommands() {
		for _, c := range cmd.SubCommands {
			list += fmt.Sprintf("%s ", c.Name)
		}
	}
	return strings.TrimSuffix(list, " ")
}

func zshSubCommandsInfo(cmd *command) string {
	var subCommandsInfo string
	if cmd.HasSubCommands() {
		subCommandsInfo += fmt.Sprintf("\n\tcmds=\"((\n")
		for _, c := range cmd.SubCommands {
			subCommandsInfo += fmt.Sprintf("\t\t%s\\:'%s'\n", c.Name, c.Description)
		}
		subCommandsInfo += fmt.Sprintf("\t))\"\n")
	}
	return subCommandsInfo
}

func zshSubCommandsLink(cmd *command) string {
	var subCommandsLink string
	if cmd.HasSubCommands() {
		subCommandsLink += fmt.Sprintf("\n\tcase $line[1] in\n")
		for _, c := range cmd.SubCommands {
			alias := func(alias []string) string {
				var res string
				for _, a := range alias {
					res += fmt.Sprintf("|%s", a)
				}
				return res
			}
			subCommandsLink += fmt.Sprintf("\t\t%s%s)\n\t\t\t%s\n\t\t\t;;\n", c.Name, alias(c.Alias), zshFunctionsName(c))
		}
		subCommandsLink += fmt.Sprintf("\tesac\n")
	}
	return subCommandsLink
}

func zshArguments(cmd *command) string {
	var arguments string
	if cmd.HasArguments() || cmd.HasFlags() || cmd.HasSubCommands() {
		arguments += fmt.Sprintf("\n\t_arguments -C \\\n")
		if cmd.HasFlags() {
			for _, f := range cmd.Flags {
				if f.HasShorthand() {
					arguments += fmt.Sprintf("\t\t{-%s,--%s}\"[%s]\"\\\n", f.Shorthand, f.Name, f.Description)
				} else {
					arguments += fmt.Sprintf("\t\t\"--%s[%s]\"\\\n", f.Name, f.Description)
				}
			}
		}
		if cmd.HasArguments() {
			for k, a := range cmd.Arguments {
				if a == "command" && cmd.HasSubCommands() {
					arguments += fmt.Sprintf("\t\t\"%d:%s:$cmds\"\\\n", k+1, a)
					arguments += fmt.Sprintf("\t\t\"*::arg:->args\"\n")
				} else {
					arguments += fmt.Sprintf("\t\t\"%d:%s\"\\\n", k+1, a)
				}
			}
		}
	}
	return arguments
}
