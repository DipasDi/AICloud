package models

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Code    int    `json:"code"`
	} `json:"error,omitempty"`
}

type AudioChatRequest struct {
	Model     string         `json:"model"`
	Reasoning *Reasoning     `json:"reasoning,omitempty"`
	Messages  []AudioMessage `json:"messages"`
}

type Reasoning struct {
	Enabled bool `json:"enabled"`
}

type AudioMessage struct {
	Role    string        `json:"role"`
	Content []ContentPart `json:"content"`
}

type ContentPart struct {
	Type       string      `json:"type"`
	Text       string      `json:"text,omitempty"`
	InputAudio *InputAudio `json:"input_audio,omitempty"`
}

type InputAudio struct {
	Data   string `json:"data"`
	Format string `json:"format"`
}
