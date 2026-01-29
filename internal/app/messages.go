package app

import task "github.com/jacobdanielrose/terminaltask/internal/task"

type TasksSavedMsg struct{ msg string }
type TasksSaveErrorMsg struct{ Err error }

type TasksLoadedMsg struct{ Tasks []task.Task }
type TasksLoadErrorMsg struct{ Err error }
