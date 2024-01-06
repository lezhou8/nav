# Nav

Minimal terminal file manager/navigator written in `go`

## Installation

*Does not work yet, this project is a work in progress*

```{go}
go install github.com/lezhou8/nav@latest
```

## Features
- Fuzzy filter searching
- Vi key bindings
- Copy, cut and pasting
- [`cd` on exit](#cd-on-exit)
- Selected files copied to clipboard on exit

## Key binds

| Key | Description |
| :-: | :---------: |
| `hjkl or arrow keys` | Basic navigation |
| `g, G` | Go to top or bottom |
| `~` | Go to home directory |
| `.` | Toggle hidden files |
| `/` | Filter search |
| `esc` | Exit filter search |
| `enter` | Accept filter search |
| `space` | Select |
| `y` | Copy/yank |
| `d` | Cut |
| `p` | Paste |
| `q` | Quit |

## `cd` on exit

Zsh/Bash

```{sh}
# Call this function whatever you like
# Add to .zshrc, .bashrc, or equivalent

function navcd() {
	nav "$@"
	cd "$(cat "${XDG_CACHE_HOME:=${HOME}/.cache}/nav/.nav_d")"
}
```

```{sh}
# You can bind it to a key

bindkey -s "^n" "navcd\n"
```

## Built with

- [Go](https://golang.org/)
- [bubbletea](https://github.com/charmbracelet/bubbletea)
- [bubbles](https://github.com/charmbracelet/bubbles)
- [lipgloss](https://github.com/charmbracelet/lipgloss)
- [fuzzy](https://github.com/sahilm/fuzzy)
- [golang-set](https://github.com/deckarep/golang-set)

## Acknowledgement

- `cd` on exit [fff](https://github.com/dylanaraps/fff/tree/master)
- fuzzy filter search [bubbles list](https://github.com/charmbracelet/bubbles)
