package llm

import (
	"bufio"
	"encoding/json"
	"gentest/internal/llm/models"
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

type Body struct {
	Prompt string `json:"prompt"`
}



func Do(c *fiber.Ctx) error {
	// Create Resty client
	client := resty.New()
	incRequestBody := &Body{}
	if err := c.BodyParser(&incRequestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("Received request with prompt: %s", incRequestBody.Prompt)

	// Set up Ollama request
	ollamaReq := models.OllamaRequest{
		Model:  "deepseek-r1:14b", // Replace with your model
		Prompt: incRequestBody.Prompt,
		Stream: true,
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		log.Printf("Error marshaling Ollama request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to marshal request",
		})
	}

	// Make streaming request to Ollama
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		SetDoNotParseResponse(true).
		Post("http://localhost:11434/api/generate")

	if err != nil {
		log.Printf("Error connecting to Ollama: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to connect to Ollama",
		})
	}
	defer resp.RawBody().Close()

	// Log the response status
	log.Printf("Ollama response status: %s", resp.Status())

	c.Set("Content-Type", "text/event-stream")
	c.Context().Response.Header.Set("Cache-Control", "no-cache")
	c.Context().Response.Header.Set("Connection", "keep-alive")


// Use c.SendStream to stream response to client
c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
	scanner := bufio.NewScanner(bufio.NewReader(resp.RawBody()))
	for scanner.Scan() {
		line := scanner.Bytes()
		log.Printf("Received chunk from Ollama: %s", line) // Log each chunk

		var ollamaResp models.OllamaResponse
		if err := json.Unmarshal(line, &ollamaResp); err != nil {
			log.Printf("Error parsing Ollama response: %v", err)
			continue
		}

		// Write chunk to client
		if _, err := w.Write([]byte("data: " + ollamaResp.Response + "\n\n")); err != nil {
			log.Printf("Error writing to client: %v", err)
			return
		}

		// Flush the writer to ensure real-time streaming
		if err := w.Flush(); err != nil {
			log.Printf("Error flushing response: %v", err)
			return
		}

		// Exit if streaming is done
		if ollamaResp.Done {
			break
		}
	}

	if err := scanner.Err(); err != nil && err != http.ErrBodyReadAfterClose {
		log.Printf("Error scanning Ollama response: %v", err)
	}
})

	return nil
}