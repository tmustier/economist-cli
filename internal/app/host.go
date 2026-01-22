package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type ScreenID string

const (
	ScreenBrowse ScreenID = "browse"
	ScreenAll    ScreenID = "all-sections"
)

type ScreenBuilder func() tea.Model

type SwitchScreenMsg struct {
	ID    ScreenID
	Reset bool
}

type Host struct {
	current  ScreenID
	builders map[ScreenID]ScreenBuilder
	screens  map[ScreenID]tea.Model
	err      error
}

func NewHost(initial ScreenID, builders map[ScreenID]ScreenBuilder) (*Host, error) {
	builder := builders[initial]
	if builder == nil {
		return nil, fmt.Errorf("missing screen builder for %q", initial)
	}
	screens := map[ScreenID]tea.Model{
		initial: builder(),
	}
	return &Host{current: initial, builders: builders, screens: screens}, nil
}

func (h Host) Init() tea.Cmd {
	if model := h.currentModel(); model != nil {
		return model.Init()
	}
	return nil
}

func (h Host) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SwitchScreenMsg:
		return h.switchTo(msg)
	}

	model := h.currentModel()
	if model == nil {
		return h, nil
	}
	updated, cmd := model.Update(msg)
	h.screens[h.current] = updated
	return h, cmd
}

func (h Host) View() string {
	if h.err != nil {
		return h.err.Error()
	}
	if model := h.currentModel(); model != nil {
		return model.View()
	}
	return ""
}

func (h Host) currentModel() tea.Model {
	if h.screens == nil {
		return nil
	}
	return h.screens[h.current]
}

func (h Host) switchTo(msg SwitchScreenMsg) (tea.Model, tea.Cmd) {
	if msg.ID == "" {
		return h, nil
	}
	if msg.Reset {
		delete(h.screens, msg.ID)
	}
	if msg.ID != h.current {
		model, ok := h.screens[msg.ID]
		if !ok {
			builder := h.builders[msg.ID]
			if builder == nil {
				h.err = fmt.Errorf("unknown screen %q", msg.ID)
				return h, nil
			}
			model = builder()
			h.screens[msg.ID] = model
		}
		h.current = msg.ID
		h.err = nil
		return h, model.Init()
	}
	return h, nil
}
