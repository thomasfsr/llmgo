package main

type Message string

type OverallState struct {
	user_id    int
	user_input string
	// messages   []Message
	task_list ListOfTasks
}
