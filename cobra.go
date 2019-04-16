package completion

import (
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Cobra adapter for completion
func Cobra(cmd *cobra.Command) ICommand {

	var cobraCommandConverterWithParent func(cmd *cobra.Command, parent *command) *command
	cobraCommandConverterWithParent = func(cmd *cobra.Command, parent *command) *command {

		use := strings.Split(cmd.Use, " ")

		var command = command{
			Name:        use[0],
			Description: cmd.Short,
			Alias:       cmd.Aliases,
			Parent:      parent,
		}

		if cmd.HasFlags() {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				if !f.Hidden {
					command.Flags = append(command.Flags, flag{
						Name:        f.Name,
						Shorthand:   f.Shorthand,
						Description: f.Usage,
					})
				}
			})
		}

		if len(use) > 1 {
			for _, a := range use[1:] {
				a := regexp.MustCompile("[^a-zA-Z0-9]+").ReplaceAllString(a, "")
				command.Arguments = append(command.Arguments, a)
			}
		}

		if len(cmd.Commands()) > 0 {
			for _, c := range cmd.Commands() {
				command.SubCommands = append(command.SubCommands, cobraCommandConverterWithParent(c, &command))
			}
		}

		return &command
	}

	return cobraCommandConverterWithParent(cmd, nil)
}
