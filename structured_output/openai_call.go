package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

func main() {
	question := "Chest Press, 5 sets of 10 reps with 22 kgs."
	ovlstate := ExtractTask(question, 1)
	fmt.Println(ovlstate)
}

type extractTask struct {
	LabelTask       string `json:"label_task" jsonschema:"enum=query_data,enum=update,enum=chat" jsonschema_description:"The label of the task if it is quering database, updating the database or neither just chatting."`
	TaskDescription string `json:"task_description" jsonschema_description:"The task with the main informations to execute the user request"`
}

type ListOfTasks struct {
	Tasks []extractTask `json:"tasks" jsonschema_description:"List of tasks extracted from the user input."`
}

type ListOfExercises struct {
	Exercises []ExerciseData `json:"exercises" jsonschema_description:"List of exercises, each exercise with its on sets."`
}

type ExerciseData struct {
	Exercise     string        `json:"exercise" jsonschema_description:"Exercise name"`
	ExerciseSets []ExerciseSet `json:"exercise_sets" jsonschema_description:"the sets of the exercise."`
}

type ExerciseSet struct {
	NReps  uint8   `json:"n_reps" jsonschema_description:"number of reps of the exercise set"`
	Weight float32 `json:"weight" jsonschema_description:"weight of the exercise set in kilograms (kg)"`
}

func GenerateSchema[T any]() *jsonschema.Schema {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

var ListOfExercisesSchema = GenerateSchema[ListOfExercises]()

func ExtractTask(user_input string, thread_id int) OverallState {
	_ = godotenv.Load()
	groq_key := os.Getenv("GROQ_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(groq_key),
		option.WithBaseURL("https://api.groq.com/openai/v1"),
	)
	ctx := context.Background()

	print("> ")
	println(user_input)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "exercises",
		Description: openai.String("Exercises extracted from users input"),
		Schema:      ListOfExercisesSchema,
		Strict:      openai.Bool(true),
	}

	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You should parse the user input to extract information about workout session. You should indentify the exercise(s) sets, each set has its own reps and weight."),
			openai.UserMessage(user_input),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		Model: "moonshotai/kimi-k2-instruct-0905",
	})

	if err != nil {
		panic(err.Error())
	}
	chat_response_content := &chat.Choices[0].Message.Content
	fmt.Println(*chat_response_content)

	listofexercises := ListOfExercises{}
	err = json.Unmarshal([]byte(*chat_response_content), &listofexercises)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("tasks: %v\n", listofexercises.Exercises)

	for _, exercise := range listofexercises.Exercises {
		fmt.Printf("- %v\n", exercise.Exercise)
		fmt.Printf("- %v\n", exercise.ExerciseSets)
	}

	chat2, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("Confirm to the user his training set was written in the database."),
			openai.UserMessage(user_input),
		},
		Model: "moonshotai/kimi-k2-instruct-0905",
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("->"+string(*chat_response_content))
	fmt.Println("->"+string(chat2.Choices[0].Message.Content))

	user_message := Message(user_input)
	return OverallState{ThreadID: thread_id, UserInput: user_input, Messages: []Message{user_message}, ExerciseList: listofexercises}
}
