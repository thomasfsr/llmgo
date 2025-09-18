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

// A struct that will be converted to a Structured Outputs response schema
// type HistoricalComputer struct {
// 	Origin       Origin   `json:"origin" jsonschema_description:"The origin of the computer"`
// 	Name         string   `json:"full_name" jsonschema_description:"The name of the device model"`
// 	Legacy       string   `json:"legacy" jsonschema:"enum=positive,enum=neutral,enum=negative" jsonschema_description:"Its influence on the field of computing"`
// 	NotableFacts []string `json:"notable_facts" jsonschema_description:"A few key facts about the computer"`
// }

// type Origin struct {
// 	YearBuilt    int64  `json:"year_of_construction" jsonschema_description:"The year it was made"`
// 	Organization string `json:"organization" jsonschema_description:"The organization that was in charge of its development"`
// }

type extractTask struct {
	LabelTask       string `json:"label_task" jsonschema:"enum=query_data,enum=update,enum=chat" jsonschema_description:"The label of the task if it is quering database, updating the database or neither just chatting."`
	TaskDescription string `json:"task_description" jsonschema_description:"The task with the main informations to execute the user request"`
}

type ListOfTasks struct {
	Tasks []extractTask `json:"tasks" jsonschema_description:"List of tasks extracted from the user input."`
}

func GenerateSchema[T any]() interface{} {
	// Structured Outputs uses a subset of JSON schema
	// These flags are necessary to comply with the subset
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return schema
}

// Generate the JSON schema at initialization time
var ListOfTasksSchema = GenerateSchema[ListOfTasks]()

func main() {
	_ = godotenv.Load()
	groq_key := os.Getenv("GROQ_API_KEY")

	client := openai.NewClient(
		option.WithAPIKey(groq_key),
		option.WithBaseURL("https://api.groq.com/openai/v1"),
	)
	ctx := context.Background()

	question := "Add 5 kilograms of rice to my data and retrieves how many eggs are there"

	print("> ")
	println(question)

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        "extract_tasks",
		Description: openai.String("Tasks extracted from users input"),
		Schema:      ListOfTasksSchema,
		Strict:      openai.Bool(true),
	}

	// Query the Chat Completions API
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(question),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{JSONSchema: schemaParam},
		},
		// Only certain models can perform structured outputs
		Model: "moonshotai/kimi-k2-instruct-0905",
	})

	if err != nil {
		panic(err.Error())
	}

	// The model responds with a JSON string, so parse it into a struct
	var listoftasks ListOfTasks
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &listoftasks)
	if err != nil {
		panic(err.Error())
	}

	// Use the model's structured response with a native Go struct
	fmt.Printf("tasks: %v\n", listoftasks.Tasks)

	for i, task := range listoftasks.Tasks {
		fmt.Printf("%v. %v\n", i+1, task.LabelTask)
		fmt.Printf("%v. %v\n", i+1, task.TaskDescription)

	}
}
