package secret_data

import (
	"fmt"
)

type SecretTypeTextData struct {
	Text     string `json:"text"`
	Metadata string `json:"metadata"`
}

func (b *SecretTypeTextData) PrintData() error {
	fmt.Printf("TEXT: %s\n", b.Text)
	fmt.Printf("METADATA: %s\n", b.Metadata)
	return nil
}

func (b *SecretTypeTextData) Fill(requireNotEmpty bool) (string, error) {
	var metadata, text string

	if metadata == "" {
		fmt.Print("enter new metadata or leave empty: ")
		//nolint:errcheck,gosec // can be left empty
		fmt.Scanln(&metadata)
	}

	if requireNotEmpty {
		for text == "" {
			fmt.Print("secret text is required, enter text: ")
			_, err := fmt.Scanln(&text)
			if err != nil {
				fmt.Printf("error: could not read text: %s\n", err.Error())
			}
		}
	} else {
		if text == "" {
			fmt.Print("enter new secret text if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&text)
		}
	}

	if text != "" {
		b.Text = text
	}
	if metadata != "" {
		b.Metadata = metadata
	}

	return "", nil
}
