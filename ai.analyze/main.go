package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/firebase/genkit/go/ai"
	"github.com/firebase/genkit/go/genkit"
	"github.com/firebase/genkit/go/plugins/compat_oai/openai"
	"github.com/openai/openai-go/option"
)

func main() {
	ctx := context.Background()

	engineURL := os.Getenv("MODEL_RUNNER_BASE_URL")
	chatModelId := "openai/" + os.Getenv("CHAT_MODEL_ID")

	fmt.Println("🌍 LLM URL:", engineURL)
	fmt.Println("🤖 Chat Model:", chatModelId)

	file, err := os.ReadFile(os.Getenv("FILE_PATH_TO_ANALYZE"))
	if err != nil {
		log.Fatalln("😡:", err)
	}
	sourceCode := string(file)

	oaiPlugin := &openai.OpenAI{
		APIKey: "I💙DockerModelRunner",
		Opts: []option.RequestOption{
			option.WithBaseURL(engineURL),
		},
	}
	g := genkit.Init(ctx, genkit.WithPlugins(oaiPlugin))

	messages := []*ai.Message{}
	messages = append(messages, ai.NewSystemTextMessage(sourceCode))

	response, err := genkit.Generate(ctx, g,
		ai.WithModelName(chatModelId),
		ai.WithSystem(os.Getenv("SYSTEM_INSTRUCTIONS")),
		ai.WithMessages(
			messages...,
		),
		ai.WithPrompt(os.Getenv("USER_MESSAGE")),

		ai.WithStreaming(func(ctx context.Context, chunk *ai.ModelResponseChunk) error {
			fmt.Print(chunk.Text())
			return nil
		}),
	)
	if err != nil {
		fmt.Println("😡 Error during generation:", err)
	}

	// Save full response to a file
	outputFile := os.Getenv("REPORT_FILE_PATH")
	err = os.WriteFile(outputFile, []byte(response.Text()), 0644)
	if err != nil {
		fmt.Println("😡 Error writing output file:", err)
	} else {
		fmt.Println("✅ Full response saved to", outputFile)
	}

	fmt.Println()
	fmt.Println()

}
