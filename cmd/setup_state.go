package cmd

import (
	"encoding/json"
	"os"
)

const stateFile = ".sdd-setup-state.json"

type setupState struct {
	Completed []string `json:"completed"`
}

func loadState() (*setupState, error) {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, err
	}
	var s setupState
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *setupState) save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(stateFile, data, 0644)
}

func (s *setupState) markDone(id string) {
	s.Completed = append(s.Completed, id)
}

func (s *setupState) isDone(id string) bool {
	for _, c := range s.Completed {
		if c == id {
			return true
		}
	}
	return false
}

func clearState() {
	os.Remove(stateFile) //nolint:errcheck
}
