package completion

import (
	"fmt"
	"strings"
)

func bash(cmd *command) string {
	var bash string
	bash += fmt.Sprintf("#!/bin/bash\n")
	bash += bashFunctions(cmd)
	bash += fmt.Sprintf("\ncomplete -F _%s %s\n", cmd.Name, cmd.Name)
	return bash
}

func bashFunctions(cmd *command) string {
	var function string
	function += fmt.Sprintf("\n%s() {\n", cmd.FullName())
	if cmd.HasSubCommands() || cmd.HasFlags() {
		function += fmt.Sprintf("\tif [ $COMP_CWORD -eq %d ]; then\n", cmd.Level())
		function += bashRootCompReply(cmd)
		if cmd.HasSubCommands() {
			function += fmt.Sprintf("\telif [ $COMP_CWORD -ge %d ]; then\n", cmd.Level()+1)
			function += bashSubCommandsCompReply(cmd)
		}
		function += fmt.Sprintf("\tfi\n")
	}
	function += fmt.Sprintf("\treturn 0\n")
	function += fmt.Sprintf("}\n")
	if cmd.HasSubCommands() {
		for _, c := range cmd.SubCommands {
			function += bashFunctions(c)
		}
	}
	return function
}

func bashRootCompReply(cmd *command) string {
	var compReply string
	if cmd.HasSubCommands() || cmd.HasFlags() {
		if cmd.HasSubCommands() {
			for _, c := range cmd.SubCommands {
				compReply += fmt.Sprintf(" %s", c.Name)
			}
		}
		if cmd.HasFlags() {
			for _, f := range cmd.Flags {
				compReply += fmt.Sprintf(" --%s", f.Name)
				if f.HasShorthand() {
					compReply += fmt.Sprintf(" -%s", f.Shorthand)
				}
			}
		}
		compReply = fmt.Sprintf("\t\tCOMPREPLY=($(compgen -W \"%s\" \"${COMP_WORDS[%d]}\"))\n",
			strings.TrimPrefix(compReply, " "), cmd.Level(),
		)
	}
	return compReply
}

func bashSubCommandsCompReply(cmd *command) string {
	var subCommands string
	if cmd.HasSubCommands() {
		subCommands += fmt.Sprintf("\t\tcase \"${COMP_WORDS[%d]}\" in\n", cmd.Level())
		for _, c := range cmd.SubCommands {
			subCommands += fmt.Sprintf("\t\t\t\"%s\"%s)\n", c.Name, c.Alias.Format("|\"", "\""))
			subCommands += fmt.Sprintf("\t\t\t\t%s\n", c.FullName())
			subCommands += fmt.Sprintf("\t\t\t\t;;\n")
		}
		subCommands += fmt.Sprintf("\t\tesac\n")
	}
	return subCommands
}
