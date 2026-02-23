package dto

type AttachmentRequest struct {
	File string `json:"file"`
}

type AttachmentResponse struct {
	URL string `json:"url"`
}
