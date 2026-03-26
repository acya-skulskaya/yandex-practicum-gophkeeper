package secret_data

import (
	"fmt"
)

type SecretTypeBankCardData struct {
	Number     string `json:"number"`
	Expiry     string `json:"expiry"`
	HolderName string `json:"holder_name"`
	CVV        string `json:"cvv"`
	Metadata   string `json:"metadata"`
}

func (b *SecretTypeBankCardData) PrintData() error {
	fmt.Printf("CARD NUMBER: %s\n", b.Number)
	fmt.Printf("CARD CVV: %s\n", b.CVV)
	fmt.Printf("CARD EXPIRY: %s\n", b.Expiry)
	fmt.Printf("CARD HOLDER NAME: %s\n", b.HolderName)
	fmt.Printf("METADATA: %s\n", b.Metadata)
	return nil
}

func (b *SecretTypeBankCardData) Fill(requireNotEmpty bool) (string, error) {
	var metadata, bcNumber, bcHolderName, bcCVV, bcExpiry string

	if metadata == "" {
		fmt.Print("enter new metadata or leave empty: ")
		//nolint:errcheck,gosec // can be left empty
		fmt.Scanln(&metadata)
	}

	if requireNotEmpty {
		for bcNumber == "" {
			fmt.Print("card number is required, enter card number: ")
			_, err := fmt.Scanln(&bcNumber)
			if err != nil {
				fmt.Printf("error: could not read card number: %s\n", err.Error())
			}
		}
		for bcHolderName == "" {
			fmt.Print("card holder name is required, enter holder name: ")
			_, err := fmt.Scanln(&bcHolderName)
			if err != nil {
				fmt.Printf("error: could not read holder name: %s\n", err.Error())
			}
		}
		for bcCVV == "" {
			fmt.Print("CVV is required, enter CVV: ")
			_, err := fmt.Scanln(&bcCVV)
			if err != nil {
				fmt.Printf("error: could not read CVV: %s\n", err.Error())
			}
		}

		for bcExpiry == "" {
			fmt.Print("card expiry date is required, enter expiry date: ")
			_, err := fmt.Scanln(&bcExpiry)
			if err != nil {
				fmt.Printf("error: could not read expiry date: %s\n", err.Error())
			}
		}
	} else {
		if bcNumber == "" {
			fmt.Print("enter new card number if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&bcNumber)
		}
		if bcHolderName == "" {
			fmt.Print("enter new card holder name if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&bcHolderName)
		}
		if bcCVV == "" {
			fmt.Print("enter new card CVV if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&bcCVV)
		}
		if bcExpiry == "" {
			fmt.Print("enter new card expiry date if you want to update it or leave empty: ")
			//nolint:errcheck,gosec // can be left empty
			fmt.Scanln(&bcExpiry)
		}
	}

	if bcNumber != "" {
		b.Number = bcNumber
	}
	if bcHolderName != "" {
		b.HolderName = bcHolderName
	}
	if bcCVV != "" {
		b.CVV = bcCVV
	}
	if bcExpiry != "" {
		b.Expiry = bcExpiry
	}
	if metadata != "" {
		b.Metadata = metadata
	}

	return "", nil
}
