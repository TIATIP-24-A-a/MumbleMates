package cmd

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

func TestNameModel(t *testing.T) {
	name := "de_dust2"

	m := initialNameModel()
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(300, 150))

	tm.Type(name)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	tm.Quit()

	fm := tm.FinalModel(t)

	switch fm := fm.(type) {
	case model:
		if fm.name != name {
			t.Errorf("Expected name to be '%s, got %q", name, fm.name)
		}
	default:
		t.Fatalf("Final model is not a model: %T", fm)
	}
}
