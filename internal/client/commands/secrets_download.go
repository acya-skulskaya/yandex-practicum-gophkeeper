package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/handlers/grpc/secret"
	"github.com/acya-skulskaya/yandex-practicum-gophkeeper/internal/client/repository/secret_data"
	"github.com/spf13/cobra"
)

func GetDownloadCmd(secretsHandler *secret.Handler) *cobra.Command {
	var secretID, secretVersionID, path string

	cmd := &cobra.Command{
		GroupID: CommandGroupSecrets,
		Use:     `download --id=SECRET_ID --version-id=SECRET_VERSION_ID --path=LOCAL_PATH_TO_SAVE_FILE`,
		Short:   "download secret's binary by secret id",
		RunE: func(cmd *cobra.Command, args []string) error {
			for secretID == "" {
				fmt.Print("secret id is required, enter id: ")
				_, err := fmt.Scanln(&secretID)
				if err != nil {
					fmt.Printf("error: could not read name: %s\n", err.Error())
				}
			}

			if secretVersionID == "" {
				fmt.Print("secret id is empty, you can specify it or leave empty: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&secretVersionID)
			}

			var fileInfo os.FileInfo
			for path == "" {
				fmt.Print("specify path to save file ot leave empty to save file to current working directory: ")
				//nolint:errcheck,gosec // can be left empty
				fmt.Scanln(&path)
				if path == "" {
					path = "./"
				}
				var err error
				fileInfo, err = os.Stat(path)
				if err != nil {
					fmt.Printf("error: could not read path: %s\n", err.Error())
					path = ""
				} else if !fileInfo.IsDir() {
					fmt.Printf("error: path shuld ne a directory\n")
					path = ""
				}
			}

			versionID, text, bs, err := secretsHandler.Download(cmd.Context(), secretID, secretVersionID)
			if err != nil {
				return fmt.Errorf("could not dowload: %w", err)
			}

			var secretData secret_data.SecretTypeBinaryData
			err = json.Unmarshal([]byte(text), &secretData)
			if err != nil {
				return fmt.Errorf("could not unmarshal secret data: %w", err)
			}

			if !strings.HasSuffix(path, "/") {
				path += "/"
			}
			filename := fmt.Sprintf("%ssecret-%s-version-%s-%s", path, secretID, versionID, secretData.FileName)

			//nolint:gosec // ignore
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o666)
			if err != nil {
				return fmt.Errorf("error opening file %s: %w", filename, err)
			}
			//nolint:errcheck // ignore err
			defer file.Close()

			writer := bufio.NewWriter(file)

			_, err = writer.Write(bs)
			if err != nil {
				return fmt.Errorf("error writing file %s to buffer: %w", filename, err)
			}

			err = writer.Flush()
			if err != nil {
				return fmt.Errorf("error flushing writer to file %s: %w", filename, err)
			}

			fmt.Printf("saved file to %s\n", filename)

			return nil
		},
	}

	cmd.Flags().StringVar(&secretID, "id", "", "secret id")
	cmd.Flags().StringVar(&secretVersionID, "version-id", "", "secret version id")
	cmd.Flags().StringVar(&path, "path", "", "local path where to put downloaded file")

	return cmd
}
