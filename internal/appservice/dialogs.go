package appservice

import "fmt"

// PickAnkiAPKG opens a file dialog to select an Anki APKG file.
func (s *Service) PickAnkiAPKG() (string, error) {
	if s.wailsApp == nil {
		return "", fmt.Errorf("file dialog unavailable")
	}
	path, err := s.wailsApp.Dialog.OpenFile().
		SetTitle("Select Anki deck (.apkg)").
		AddFilter("Anki deck", "*.apkg").
		AddFilter("All files", "*.*").
		PromptForSingleSelection()
	if err != nil {
		return "", fmt.Errorf("open file dialog: %w", err)
	}
	if path == "" {
		return "", nil
	}
	return path, nil
}
