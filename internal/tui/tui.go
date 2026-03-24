// Package tui provides a Terminal User Interface (TUI) for real-time
// monitoring of health check results using the Bubble Tea framework.
package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gohealthchecker/internal/checker"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	table    table.Model
	checker  *checker.Checker
	results  <-chan checker.Result
	stateMap map[string]*checker.TargetState
}

// Msg used to update the TUI when a result arrives
type resultMsg checker.Result

func (m model) Init() tea.Cmd {
	return m.waitForResult()
}

func (m model) waitForResult() tea.Cmd {
	return func() tea.Msg {
		res, ok := <-m.results
		if !ok {
			return nil
		}
		return resultMsg(res)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case resultMsg:
		res := checker.Result(msg)
		// Process result for alerts (side effect)
		m.checker.ProcessResult(res, m.stateMap)

		// Update table rows
		rows := m.table.Rows()
		updated := false
		for i, row := range rows {
			if row[0] == res.Target.URL {
				rows[i] = m.formatRow(res)
				updated = true
				break
			}
		}
		if !updated {
			rows = append(rows, m.formatRow(res))
		}
		m.table.SetRows(rows)
		return m, m.waitForResult()
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) formatRow(res checker.Result) table.Row {
	status := "✅ OK"
	if res.Error != nil || res.Status != 200 {
		status = fmt.Sprintf("❌ FAIL (%d)", res.Status)
		if res.Error != nil {
			status = "❌ ERR"
		}
	}

	return table.Row{
		res.Target.URL,
		status,
		res.Duration.String(),
		res.LastCheck.Format("15:04:05"),
	}
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n  Press 'q' to quit.\n"
}

// StartUI initializes and runs the Bubble Tea TUI program.
// It requires a [checker.Checker] and a channel of [checker.Result]s
// to display and manage real-time updates.
func StartUI(c *checker.Checker, results <-chan checker.Result) error {
	columns := []table.Column{
		{Title: "Service URL", Width: 40},
		{Title: "Status", Width: 15},
		{Title: "Latency", Width: 10},
		{Title: "Last Check", Width: 12},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{
		table:    t,
		checker:  c,
		results:  results,
		stateMap: make(map[string]*checker.TargetState),
	}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}
	return nil
}
