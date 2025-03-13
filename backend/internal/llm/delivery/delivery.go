package delivery

import (
	"bufio"
	"encoding/json"
	"gentest/internal/llm"
	"gentest/internal/llm/models"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const PromptStart = "Write ONLY code. You are given next piece of code, Output should consist solely of the test code without any explanations, comments, or formatting. The code down below: "

type Delivery struct {
	Usecase llm.Usecase
}

func NewDelivery(usecase llm.Usecase) *Delivery {
	return &Delivery{
		Usecase: usecase,
	}
}

func (d *Delivery) GenerateFirstTest(c *gin.Context) {
	var incRequestBody models.GenTestRequest
	if err := c.ShouldBindJSON(&incRequestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	Prompt := PromptStart + incRequestBody.Prompt

	genTestReq := &models.OllamaRequest{
		Model: "deepseek-r1:14b",
		Prompt: Prompt,
		Stream: true,
	}

	resp, err := d.Usecase.GenerateResponse(genTestReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get response from Ollama"})
	}
	defer resp.RawBody().Close()

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	c.Stream(func(w io.Writer) bool {
		scanner := bufio.NewScanner(resp.RawBody())
		for scanner.Scan() {
			line := scanner.Bytes()
			log.Printf("Received chunk from Ollama: %s", line)

			var ollamaResp models.OllamaResponse
			if err := json.Unmarshal(line, &ollamaResp); err != nil {
				log.Printf("Error parsing Ollama response: %v", err)
				continue
			}

			if _, err := w.Write([]byte(ollamaResp.Response)); err != nil {
				log.Printf("Error writing to client: %v", err)
				return false
			}

			c.Writer.Flush()

			if ollamaResp.Done {
				return false
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error scanning Ollama response: %v", err)
		}

		return false
	})	
}

func GenerateAdditionalTests(c *gin.Context) {
	
}