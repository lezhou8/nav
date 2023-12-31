# Nav

Minimal terminal file manager/navigator written in `go`

## Installation

*Does not work yet, this project is a work in progress*

```{go}
go install github.com/lezhou8/nav@latest
```

## Features
- Searching
- Save cursor locations
- Vi key bindings
- Toggle the display of hidden files
- Shortcut to home directory

## Key binds

| Key | Description |
| :-: | :---------: |
| `hjkl or arrow keys` | Basic navigation |
| `g, G` | Go to top or bottom |
| `~` | Go to home directory |
| `.` | Toggle hidden files |
| `/` | Filter search |
| `q` | Quit |

## Built with

- [Go](https://golang.org/)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
