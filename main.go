package main

import (
	"fmt"
	"os"
    "encoding/json"
    "io/ioutil"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type status int

const divisor = 4

const (
    todo status = iota
    inProgress
    done
)

/* MODEL MANAGEMENT*/
var models []tea.Model
const (
    model status = iota
    form
)

/* STYLING */
var (
    columnStyle = lipgloss.NewStyle().Padding(1, 2)
    focusedStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("62"))
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241"))

)
     

/* Custom Item */
type JsonTask struct {
    Status  status
    Title   string
    Description string
}

func (t *JsonTask) GetTitle() string {
    return t.Title
}

func (t *JsonTask) GetDescription() string {
    return t.Description
}



func getTasksJson() []JsonTask {
    var tasks []JsonTask

    file, err := os.Open("tasks.json")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    data, err := ioutil.ReadAll(file)
    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(data, &tasks)
    if err != nil {
        panic(err)
    }


    return tasks

}



type Task struct {
    status  status
    title   string
    description string
}

func NewTask(status status, title, description string) Task {
    return Task{status: status, title: title, description: description} 
}

func (t *Task) Next() {
    if t.status == done {
        t.status = todo
    } else {
        t.status++
    }
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
    focused status
    lists    []list.Model
    err     error
    loaded  bool
    quitting    bool
}

func New() *Model {
    return &Model{}
}

func (m *Model) MoveToNext() tea.Msg {
    selectedItem := m.lists[m.focused].SelectedItem()
    selectedTask := selectedItem.(Task)
    m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
    selectedTask.Next()
    m.lists[selectedTask.status].InsertItem(len(m.lists[selectedTask.status].Items())-1, list.Item(selectedTask))
    return nil



}

//TODO: Delete selected task
func (m *Model) DeleteTask() tea.Msg{
    selectedItem := m.lists[m.focused].SelectedItem()
    selectedTask := selectedItem.(Task)
    m.lists[selectedTask.status].RemoveItem(m.lists[m.focused].Index())
    return nil

}


// TODO: Go to next list
func (m *Model) Next() {
    if m.focused == done {
        m.focused = todo
    } else {
        m.focused++
    }
}


// TODO: Go to prev list
func (m *Model) Prev() {
    if m.focused == todo {
        m.focused = done
    } else {
        m.focused--
    }
}



func (m *Model) initLists(width, height int) {
    js := getTasksJson()
    var tasks []Task
    for _, j := range js {
        tasks = append(tasks, NewTask(j.Status, j.Title, j.Description))
    }

    defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), width / divisor, height/2)
    defaultList.SetShowHelp(false)
    m.lists = []list.Model{defaultList, defaultList, defaultList}
    //Init To Do
    m.lists[todo].Title = "To Do"
    m.lists[todo].SetItems([]list.Item{
        Task{status: todo, title: "buy milk", description: "Chocolate Milk"},
        Task{status: todo, title: "eat sushi", description: "negitoro roll, miso soup, rice"},
        Task{status: todo, title: "fold laundry", description: "or wear wrinkly t-shirts"},
        Task{status: tasks[0].status, title: tasks[0].title, description: tasks[0].description},
    })
    //Init in progress
    m.lists[inProgress].Title = "In Progress"
    m.lists[inProgress].SetItems([]list.Item{
        Task{status: inProgress, title: "write code", description: "don't worry, it's Go"},
    })
    //Init done
    m.lists[done].Title = "Done"
    m.lists[done].SetItems([]list.Item{
        Task{status: done, title: "stay cool", description: "as a cucumber"},
    })

}

func (m Model) Init() tea.Cmd{
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    case tea.WindowSizeMsg:
        if !m.loaded {
        columnStyle.Width(msg.Width / divisor)
        focusedStyle.Width(msg.Width / divisor)
        columnStyle.Height(msg.Height - divisor)
        focusedStyle.Height(msg.Height - divisor)
        m.initLists(msg.Width, msg.Height)
        m.loaded = true
    }
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.quitting = true
            return m, tea.Quit
        case "left", "h":
            m.Prev()
        case "right", "l":
            m.Next()
        case "enter":
            return m, m.MoveToNext
        case "n":
            models[model] = m
            models[form] = NewForm(m.focused)
            return models[form].Update(nil)
        case "backspace":
            return m, m.DeleteTask

        }
    case Task:
        task := msg
        return m, m.lists[task.status].InsertItem(len(m.lists[task.status].Items()), task)

    }
    var cmd tea.Cmd
    m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
    return m, cmd

}

func (m Model) View() string {
    if m.quitting {
        return "" }
    if m.loaded {
        todoView := m.lists[todo].View()
        inProgView := m.lists[inProgress].View()
        doneView := m.lists[done].View()
        switch m.focused {
        case inProgress:
            return lipgloss.JoinHorizontal(
            lipgloss.Left, 
            columnStyle.Render(todoView),
            focusedStyle.Render(inProgView),
            columnStyle.Render(doneView),
            )

        case done:
            return lipgloss.JoinHorizontal(
            lipgloss.Left, 
            columnStyle.Render(todoView),
            columnStyle.Render(inProgView),
            focusedStyle.Render(doneView),
            )
        default:
            return lipgloss.JoinHorizontal(
            lipgloss.Left, 
            focusedStyle.Render(todoView),
            columnStyle.Render(inProgView),
            columnStyle.Render(doneView),
            )
        }
    
    } else {
        return "loading..."
}
}

/* FORM MODEL */
type Form struct{
    focused status
    title textinput.Model
    description textarea.Model
}

func NewForm(focused status) *Form {
    form := &Form{focused: focused}
    form.title = textinput.New()
    form.title.Focus()
    form.description = textarea.New()
    return form
}

func (m Form) CreateTask() tea.Msg {
    task := NewTask(m.focused, m.title.Value(), m.description.Value())
    return task

}

func (m Form) Init() tea.Cmd {
    return nil
}

func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "enter":
            if m.title.Focused() {
                m.title.Blur()
                m.description.Focus()
                return m, textarea.Blink
            } else {
                models[form] = m
                return models[model], m.CreateTask
            }
        }
    }

    if m.title.Focused() {
        m.title, cmd = m.title.Update(msg)
        return m, cmd
    } else {
        m.description, cmd = m.description.Update(msg)
        return m, cmd
    }

}

func (m Form) View() string {
    return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}

func main() {
    models = []tea.Model{New(), NewForm(todo)}
    m := models[model]
    p := tea.NewProgram(m)

    
    if err := p.Start(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

}
