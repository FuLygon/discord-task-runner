package tasker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type fcmPayload struct {
	ValidateOnly bool           `json:"validate_only"`
	Message      payloadMessage `json:"message"`
}

type payloadMessage struct {
	Token   string `json:"token"`
	Android struct {
		Priority string `json:"priority"`
	} `json:"android"`
	Data map[string]string `json:"data"`
}

func ExecuteTask(projectId, deviceToken, task string, variables map[string]string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// prepare payload
	payload := fcmPayload{
		ValidateOnly: false,
		Message: payloadMessage{
			Token: deviceToken,
		},
	}
	payload.Message.Android.Priority = "high"
	messageData := make(map[string]string)
	messageData["task"] = task
	if len(variables) > 0 {
		for variable, value := range variables {
			messageData["%"+variable] = value
		}
	}
	payload.Message.Data = messageData
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to prepare message payload: %w", err)
	}

	// prepare request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("https://fcm.googleapis.com/v1/projects/%s/messages:send", projectId),
		bytes.NewBuffer(payloadJSON))
	if err != nil {
		return fmt.Errorf("failed to prepare message request: %w", err)
	}

	accessToken, err := getAccessToken(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", accessToken))

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message request: %w", err)
	}
	defer resp.Body.Close()

	// Return error if status code is not 200
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message request with status: %s", resp.Status)
	}

	return nil
}
