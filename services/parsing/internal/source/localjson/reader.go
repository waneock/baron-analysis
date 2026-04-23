package localjson

import (
	"encoding/json"
	"fmt"
	"os"
	"skinbaron-analyzer/services/parsing/internal/domain"
)

var (
	ErrEmptyFilePath     = fmt.Errorf("file path is empty")
	ErrFileDoesNotExists = fmt.Errorf("file does not exists")
)

type ItemJSON struct {
	ID    string     `json:"id"`
	Name  string     `json:"name"`
	Wears []WearJSON `json:"wears"`
}

type WearJSON struct {
	Name string `json:"name"`
}

func ReadItemJSON(path string) (*[]domain.Item, error) {
	if path == "" {
		return nil, ErrEmptyFilePath
	}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, ErrFileDoesNotExists
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var items []ItemJSON
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}

	domainItems := make([]domain.Item, 0, len(items))
	jsonItemsToDomainItems(items, &domainItems)

	return &domainItems, nil
}

func jsonItemsToDomainItems(input []ItemJSON, domainItems *[]domain.Item) {
	for _, item := range input {
		addJsonItemToDomainItems(item, domainItems)
	}
}

func addJsonItemToDomainItems(input ItemJSON, domainItems *[]domain.Item) {
	wears := make([]string, 0, len(input.Wears))
	for _, wear := range input.Wears {
		wears = append(wears, wear.Name)
	}
	domainItem := domain.Item{
		ID:    input.ID,
		Name:  input.Name,
		Wears: wears,
	}

	*domainItems = append(*domainItems, domainItem)
}
