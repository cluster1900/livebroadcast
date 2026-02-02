package centrifugo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	URL     string
	APIKey  string
	Secret  string
	Timeout time.Duration
}

type PublishRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type PublishResponse struct {
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

type Client struct {
	config Config
	client *http.Client
}

func NewClient(url, apiKey string) *Client {
	if apiKey == "" {
		apiKey = "your_centrifugo_api_key"
	}
	return &Client{
		config: Config{
			URL:     url,
			APIKey:  apiKey,
			Timeout: 5 * time.Second,
		},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Publish(channel string, data interface{}) error {
	reqBody := PublishRequest{
		Method: "publish",
		Params: map[string]interface{}{
			"channel": channel,
			"data":    data,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.config.URL+"/api", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" && c.config.APIKey != "your_centrifugo_api_key" {
		req.Header.Set("Authorization", "apikey "+c.config.APIKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("centrifugo returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result PublishResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Error != nil {
		return fmt.Errorf("centrifugo error: %v", result.Error)
	}

	return nil
}

func (c *Client) GenerateConnectionToken(userID string, expireAt int64) string {
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}
	payload := map[string]interface{}{
		"sub": userID,
		"exp": expireAt,
	}

	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)

	headerB64 := base64URLEncode(headerJSON)
	payloadB64 := base64URLEncode(payloadJSON)

	signature := c.sign(headerB64 + "." + payloadB64)

	return headerB64 + "." + payloadB64 + "." + signature
}

func (c *Client) sign(message string) string {
	h := hmac.New(sha256.New, []byte(c.config.Secret))
	h.Write([]byte(message))
	return base64URLEncode(h.Sum(nil))
}

func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func GetChannels(roomID string) []string {
	return []string{
		fmt.Sprintf("room_%s", roomID),
		fmt.Sprintf("room_%s:streamer", roomID),
	}
}

type DanmuMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		ID       string   `json:"id"`
		UserID   string   `json:"user_id"`
		Nickname string   `json:"nickname"`
		Level    int      `json:"level"`
		Avatar   string   `json:"avatar"`
		Content  string   `json:"content"`
		Color    string   `json:"color"`
		Badges   []string `json:"badges"`
	} `json:"data"`
}

type GiftMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		Sender struct {
			ID       string `json:"id"`
			Nickname string `json:"nickname"`
			Level    int    `json:"level"`
		} `json:"sender"`
		Gift struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Icon      string `json:"icon"`
			Animation string `json:"animation"`
		} `json:"gift"`
		Count      int `json:"count"`
		Combo      int `json:"combo"`
		TotalValue int `json:"total_value"`
	} `json:"data"`
}

type OnlineCountMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		Count int `json:"count"`
	} `json:"data"`
}

type StreamStatusMessage struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Data      struct {
		Status string `json:"status"`
		Reason string `json:"reason,omitempty"`
	} `json:"data"`
}
