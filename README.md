# terminaltask

[![CI](https://github.com/jacobdanielrose/terminaltask/actions/workflows/ci.yml/badge.svg)](https://github.com/jacobdanielrose/terminaltask/actions/workflows/ci.yml)
[![Release](https://github.com/jacobdanielrose/terminaltask/actions/workflows/release.yml/badge.svg)](https://github.com/jacobdanielrose/terminaltask/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/jacobdanielrose/terminaltask/branch/main/graph/badge.svg)](https://codecov.io/gh/jacobdanielrose/terminaltask)
[![Go Report Card](https://goreportcard.com/badge/github.com/jacobdanielrose/terminaltask)](https://goreportcard.com/report/github.com/jacobdanielrose/terminaltask)

terminaltask is a keyboard-driven task manager that runs directly in your terminal. It provides a clean, minimal interface for capturing, reviewing, and updating tasks without leaving the command line.

![](https://raw.githubusercontent.com/jacobdanielrose/terminaltask/main/assets/terminaltask-demo.gif)

## Features

- **Task management:** Create, edit, delete, and mark tasks as completed.
- **List view:** Navigate through tasks with a focused, scrollable list.
- **Keyboard first:** Drive everything with keys – no mouse required.
- **Date picker:** Set due dates via a keyboard-driven date picker.
- **Inline help:** Toggle contextual help with key bindings in both list and edit views.
- **Themed UI:** Uses `lipgloss` and other Charm libraries for a pleasant terminal UI.

## Installation

### Homebrew (macOS)

You can install `terminaltask` using Homebrew via the tap:

```/dev/null/sh#L1-2
brew tap jacobdanielrose/homebrew-tap
brew install --cask terminaltask
```

After installation, macOS may block the app the first time you try to run it because it is from an unidentified developer.  
If this happens, open **System Settings → Privacy & Security**, scroll down to the security section, and explicitly allow `terminaltask` to run. Once approved, you can launch it normally from your terminal.

### Building from Source

Ensure you have Go installed (version 1.20 or higher is recommended).

1. Clone the repository:

   ```/dev/null/sh#L1-3
   git clone https://github.com/jacobdanielrose/terminaltask.git
   cd terminaltask
   ```

2. Build the application using the Makefile:

   ```/dev/null/sh#L5-6
   make build
   ./bin/terminaltask
   ```

### Downloading Binaries

You can also download pre-built binaries from the GitHub releases page:

- https://github.com/jacobdanielrose/terminaltask/releases

Download the appropriate archive for your OS, extract it, and run the binary from your terminal.

## Usage

- **Navigation:**
  - Use arrow keys or `j/k` to move through the task list.
  - Press `n` to create a new task.
  - Press `e` to edit the currently selected task.
  - Press `space` to toggle a task as completed.
  - Press `esc` to exit edit mode.
  - Press `ctrl+c` at any time to quit.

- **Editing:**
  - Text input fields for title and description.
  - [bubble-datepicker](https://github.com/EthanEFung/bubble-datepicker) (thanks EthanEFung!) for selecting task due dates.
  - Use `enter` to move between fields.
  - When the date picker is focused, press `tab` to move focus from the year/month header down to the calendar, and `shift+tab` to move focus back up to the header.
  - Press `ctrl+s` to save changes.
  - Press `esc` to exit the edit menu without saving.

- **Shortcuts:**
  - `?` to toggle the help menu and view key bindings in the list view.
  - `ctrl+o` to toggle the help menu and view key bindings in the edit view.

## Testing

This project uses Go’s standard testing tools.

- Run all tests:

  ```/dev/null/sh#L1-1
  make test
  ```

- Run tests with coverage (-cover):

  ```/dev/null/sh#L1-1
  make cover
  ```

- Generate a coverage profile:

  ```/dev/null/sh#L1-3
  make coverage
  ```

Continuous integration runs the test suite and updates coverage on every push and pull request.

## Development / Roadmap

This project is actively developed. Some planned enhancements include:

- [ ] Expanding test coverage across core packages (`internal/task`, `internal/store`, `internal/app`).
- [ ] Filtering out completed tasks.
- [ ] Filtering by description.
- [ ] Synchronizing with an external CalDAV server.


If you have ideas for improvements or new features, feel free to open an issue or reach out at `jacobdanielrose@protonmail.com`.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.  
Tests are appreciated where practical, especially for core behavior in the task model, storage, and app update logic.

## License

This project is licensed under the GNU General Public License. See the [LICENSE](LICENSE) file for details.

## Contact

For questions or support, please feel free to open an issue on GitHub.
