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

var ListOfTasksSchema = GenerateSchema[ListOfTasks]()

func main() {
	question := "Add 5 kilograms of rice to my data and retrieves how many eggs are there"
	ovlstate := ExtractTask(question, 1)
	fmt.Println("STOP")
	fmt.Println(ovlstate)
}

func ExtractTask(user_input string, user_id int) OverallState {
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
		Name:        "extract_tasks",
		Description: openai.String("Tasks extracted from users input"),
		Schema:      ListOfTasksSchema,
		Strict:      openai.Bool(true),
	}

	// Query the Chat Completions API
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a assistant that parses the requests of the user for another agent to process them."),
			openai.UserMessage(user_input),
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

	for _, task := range listoftasks.Tasks {
		fmt.Printf("- %v\n", task.LabelTask)
		fmt.Printf("- %v\n", task.TaskDescription)

	}
	user_message := Message(user_input)
	return OverallState{user_id: user_id, user_input: user_input, messages: []Message{user_message}, task_list: listoftasks}
}
