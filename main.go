package main

import (
    "github.com/charmbracelet/bubbles/list"
)

type status int

const (
    todo status = iota
    inProgress
    done
)

type Task struct {
    status  status
    title   string
    description string
}

// Implement the list.Item Interface
func (t Task) FilterValue() string {
    return t.title
}

func (t Task) Title() string {
    return t.title
}

func (t Task) Description() string {
    return t.description
}


/* Main Model */

type Model struct {
    list    list.Model
    err     error
}

func (m *Model) initList() {
    m.list = list.New([]list.Item, list.NewDefaultDelegate())
}


