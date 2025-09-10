package main

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	// ...
)

func main() {
	_ = godotenv.Load()
	groq_key := os.Getenv("GROQ_API_KEY")
	client := openai.NewClient(
		option.WithAPIKey(groq_key),
		option.WithBaseURL("https://api.groq.com/openai/v1"),
	)
	chatCompletion, err := client.Chat.Completions.New(context.Background(),
		openai.ChatCompletionNewParams{
			Model: "llama-3.1-8b-instant",
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a helpful assistant that responds sucintly"),
				openai.UserMessage("Can I call tools using you?"),
			},
		})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
}
