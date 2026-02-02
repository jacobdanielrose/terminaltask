package app

import task "github.com/jacobdanielrose/terminaltask/internal/task"

// TasksSavedMsg indicates tasks were saved successfully.
type TasksSavedMsg struct{ msg string }

// TasksSaveErrorMsg indicates an error occurred while saving tasks.
type TasksSaveErrorMsg struct{ Err error }

// TasksLoadedMsg carries tasks loaded from the service.
type TasksLoadedMsg struct{ Tasks []task.Task }

// TasksLoadErrorMsg indicates an error occurred while loading tasks.
type TasksLoadErrorMsg struct{ Err error }
