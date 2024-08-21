package main

type Message struct {
	Kind string `json:"kind"` // txt, png or jpg
	Data string `json:"data"` // txt or base64 encoded image
}

type ModRequest struct {
	ID        string  `json:"id"`
	ClientID  string  `json:"client_id"`
	Kind      string  `json:"kind"`
	Message   Message `json:"message"`
	Approved  bool    `json:"approved"`
	Moderated bool    `json:"moderated"`
}
