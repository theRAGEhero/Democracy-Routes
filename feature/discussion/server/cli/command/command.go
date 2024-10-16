package command

import "strings"

func IsHelp(cmd string) bool {
	return cmd == "-h" || strings.Contains(cmd, "help")
}
