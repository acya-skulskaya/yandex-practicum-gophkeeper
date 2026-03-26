package secret_data

import (
	"encoding/json"
	"fmt"
)

type Repo[T Model] interface {
	Print(data string) error
	GetJSON(data string, requireNotEmpty bool) ([]byte, string, error)
}

type Model interface {
	PrintData() error
	Fill(requireNotEmpty bool) (string, error)
}

type StandardInputOutputRepo[T Model] struct {
}

func NewStandardInputOutputRepo[T Model]() *StandardInputOutputRepo[T] {
	return &StandardInputOutputRepo[T]{}
}

func (i *StandardInputOutputRepo[T]) Print(data string) error {
	var dataModel T

	err := json.Unmarshal([]byte(data), &dataModel)
	if err != nil {
		return fmt.Errorf("could not unmarshall secretVersionData %w", err)
	}

	dataModel.PrintData()

	return nil
}

func (i *StandardInputOutputRepo[T]) GetJSON(data string, requireNotEmpty bool) ([]byte, string, error) {
	var dataModel T

	if data != "" {
		err := json.Unmarshal([]byte(data), &dataModel)
		if err != nil {
			return []byte{}, "", fmt.Errorf("could not unmarshall secretVersionData %w", err)
		}
	}

	binaryFilePath, err := dataModel.Fill(requireNotEmpty)

	secretDataJSON, err := json.Marshal(dataModel)
	if err != nil {
		return []byte{}, "", fmt.Errorf("could not marshal secretData: %w", err)
	}

	return secretDataJSON, binaryFilePath, nil
}
