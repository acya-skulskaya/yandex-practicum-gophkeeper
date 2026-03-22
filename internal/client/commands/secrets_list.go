package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/spf13/cobra"
)

func GetListCmd(secretsHandler *secret.Handler) *cobra.Command {
	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use:     `list`,
		Short:   "list all secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			list, err := secretsHandler.List(cmd.Context())
			if err != nil {
				return fmt.Errorf("could list secrets: %w", err)
			}

			if len(list) == 0 {
				fmt.Println("no secrets found")
				return nil
			}

			emdash := "—"
			dash := "–"
			format := "%-5s | %-15s | %-20s | %-25s | %-25s | %s\n"

			fmt.Println(strings.Repeat(emdash, 115))
			fmt.Printf(format, "ID", "TYPE", "NAME", "CREATED", "UPDATED", "NUM.VER.")
			fmt.Println(strings.Repeat(dash, 115))

			for _, s := range list {
				fmt.Printf(format,
					strconv.Itoa(int(s.ID)),
					s.Type,
					s.Name,
					s.CreatedAt.Format(time.DateTime),
					s.UpdatedAt.Format(time.DateTime),
					strconv.Itoa(len(s.Versions)))
			}

			fmt.Println(strings.Repeat(emdash, 115))

			return nil
		},
	}

	return cmd
}
