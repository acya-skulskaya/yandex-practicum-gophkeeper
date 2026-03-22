package commands

import (
	"fmt"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/config"
	"github.com/spf13/cobra"
)

func GetVersionCmd(info config.BuildInfo) *cobra.Command {
	infoStr := config.GetBuildInfoStr(info)

	cmd := &cobra.Command{
		GroupID: CommandGroupApp,
		Use:     "version",
		Short:   "Prints version info of the app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(infoStr)
		},
	}

	return cmd
}
