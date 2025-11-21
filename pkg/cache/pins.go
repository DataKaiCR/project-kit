package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// PinRecord represents a pinned project in a slot
type PinRecord struct {
	Slot        int    `json:"slot"`
	ProjectID   string `json:"project_id"`
	ProjectPath string `json:"project_path"`
}

// GetPinsFile returns the path to the pins file
func GetPinsFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "pk")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "pins.json"), nil
}

// LoadPins reads all pinned projects
func LoadPins() (map[int]PinRecord, error) {
	pinsFile, err := GetPinsFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(pinsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[int]PinRecord), nil
		}
		return nil, err
	}

	var pins map[int]PinRecord
	if err := json.Unmarshal(data, &pins); err != nil {
		return nil, err
	}

	return pins, nil
}

// SavePins writes pinned projects to disk
func SavePins(pins map[int]PinRecord) error {
	pinsFile, err := GetPinsFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(pins, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(pinsFile, data, 0644)
}

// AddPin pins a project to a specific slot (1-5)
func AddPin(slot int, projectID, projectPath string) error {
	if slot < 1 || slot > 5 {
		return fmt.Errorf("slot must be between 1 and 5")
	}

	pins, err := LoadPins()
	if err != nil {
		return err
	}

	pins[slot] = PinRecord{
		Slot:        slot,
		ProjectID:   projectID,
		ProjectPath: projectPath,
	}

	return SavePins(pins)
}

// RemovePin removes a pin by slot number
func RemovePin(slot int) error {
	pins, err := LoadPins()
	if err != nil {
		return err
	}

	if _, exists := pins[slot]; !exists {
		return fmt.Errorf("no pin in slot %d", slot)
	}

	delete(pins, slot)
	return SavePins(pins)
}

// RemovePinByProject removes a pin by project ID
func RemovePinByProject(projectID string) error {
	pins, err := LoadPins()
	if err != nil {
		return err
	}

	found := false
	for slot, pin := range pins {
		if pin.ProjectID == projectID {
			delete(pins, slot)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("project '%s' is not pinned", projectID)
	}

	return SavePins(pins)
}

// GetPin retrieves a pin by slot number
func GetPin(slot int) (*PinRecord, error) {
	pins, err := LoadPins()
	if err != nil {
		return nil, err
	}

	pin, exists := pins[slot]
	if !exists {
		return nil, fmt.Errorf("no pin in slot %d", slot)
	}

	return &pin, nil
}

// ListPins returns all pins sorted by slot
func ListPins() ([]PinRecord, error) {
	pins, err := LoadPins()
	if err != nil {
		return nil, err
	}

	var pinList []PinRecord
	for _, pin := range pins {
		pinList = append(pinList, pin)
	}

	// Sort by slot number
	sort.Slice(pinList, func(i, j int) bool {
		return pinList[i].Slot < pinList[j].Slot
	})

	return pinList, nil
}

// ClearPins removes all pins
func ClearPins() error {
	return SavePins(make(map[int]PinRecord))
}

// IsPinned checks if a project is pinned (returns slot number or -1)
func IsPinned(projectID string) int {
	pins, err := LoadPins()
	if err != nil {
		return -1
	}

	for slot, pin := range pins {
		if pin.ProjectID == projectID {
			return slot
		}
	}

	return -1
}
