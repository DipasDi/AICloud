package models

type EmbedRequest struct {
	Model      string `json:"model"`
	Input      string `json:"input"`
	Dimensions int    `json:"dimensions,omitempty"`
}

type EmbedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
	} `json:"data"`
}
