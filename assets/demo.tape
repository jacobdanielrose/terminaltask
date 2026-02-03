# Set up the terminal window
Set Shell "bash"
Set Width 800
Set Height 800
Set FontSize 14

# Start recording
Output "terminaltask-demo.gif"

# Clear the screen and show a prompt
Type "clear"
Enter
Sleep 500ms

# Run terminaltask
Type "terminaltask"
Sleep 500ms
Enter
Sleep 1s

# side scroll
Sleep 500ms
Right
Sleep 500ms
Right
Sleep 500ms

# Create a new task (press 'n')
Type "n"
Sleep 500ms

# Type title
Type "Write VHS demo"
Sleep 500ms
Enter

# Type description
Type "Record a terminal session using VHS for the README."
Sleep 500ms
Enter

# At this point the date picker is focused (per your README),
# we'll just accept the default date and save.
Sleep 500ms
Down
Sleep 500ms
Right
Sleep 500ms
Right
Sleep 1s

# Save the task (ctrl+s)
Ctrl+S
Sleep 1s

# toggle app done
Space
Sleep 1s

# delete task
Type "r"
Sleep 1s

# Quit the app (ctrl+c)
Ctrl+c
Sleep 500ms
