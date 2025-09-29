package main

type Message string

type OverallState struct {
	ThreadID     int
	UserInput    string
	Messages     []Message
	ExerciseList ListOfExercises
}
