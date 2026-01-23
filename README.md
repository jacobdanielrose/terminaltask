# TerminalTask

TerminalTask is a task management application designed to run directly in your terminal. It provides a simple and intuitive interface for managing reminders and to-do lists using the minimalist, text-based charm of terminal applications.

## Features

- **Task Management:** Create, edit, delete, and mark tasks as completed.
- **Visual Interface:** Navigate through tasks with a comprehensive list view.
- **Styling:** Uses `lipgloss` for a visually appealing terminal UI.
- **Date Picker:** Easily set due dates for tasks with a keyboard-driven date picker.
- **Help Menu:** Quickly access help information and keyboard shortcuts.

## Installation

Ensure you have Go installed (version 1.16 or higher is recommended).

1. Clone the repository:
   ```sh
   git clone https://github.com/yourusername/terminaltask.git
   ```

2. Change to the project directory:
   ```sh
   cd terminaltask
   ```

3. Build the application:
   ```sh
   go build
   ```

4. Run the application:
   ```sh
   ./terminaltask
   ```

## Usage

- **Navigation:**
  - Use arrow keys and `j/k` to navigate through the task list.
  - Press `e` to edit a selected task.
  - Press `n` to create a new task.
  - Use `esc` to exit the edit mode.
  - `ctrl+c` at any time quits the program.
  - `space` to mark a task as completed

- **Editing:**
  - Text input fields for title, description
  - [Datepicker](https://github.com/EthanEFung/bubble-datepicker) (thanks EthanEFung!) for selecting the task's due date.
  - Save changes by navigating through fields `enter` and saving `ctrl+s` when done.
  - To exit the edit menu (without saving) press `esc` 

- **Shortcuts:**
  - `?` to toggle the help menu and view key bindings in the list view
  - `ctrl+o` to toggle the help menu and view key bindings in the edit view


## Development

This project is still very much being actively developed. There is currently no persistence (think of this as as pre-alpha), so I plan to add a DB backend to add persistence. Other smaller tweaks and big changes include:

- Ability to filter out already completed tasks
- Ability to filter by descriptions
- Ability to synchronize with an external calDav server (will definitely come last)

If you can think of any helpful tweaks and even completely new features, please feel free to reach out at `jacobdanielrose@protonmail.com`
## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your updates. I'm not the disciplined developer ever so I will probably only start writing tests when I'm finished with the main sections of the project. Feel free to create or include your own tests if you please.

## License

This project is licensed under the GNU General Public License. See the [LICENSE](LICENSE) file for details.

## Contact

For questions or support, please feel free to open an issue on GitHub.
