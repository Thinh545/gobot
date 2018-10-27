package main

type (
	Request struct {
		Object string `json:"object,omitempty"`
		Entry  []struct {
			ID        string      `json:"id,omitempty"`
			Time      int64       `json:"time,omitempty"`
			Messaging []Messaging `json:"messaging,omitempty"`
		} `json:"entry,omitempty"`
	}
	Messaging struct {
		Sender    *User    `json:"sender,omitempty"`
		Recipient *User    `json:"recipient,omitempty"`
		Timestamp int      `json:"timestamp,omitempty"`
		Message   *Message `json:"message,omitempty"`
	}

	User struct {
		ID string `json:"id,omitempty"`
	}

	Message struct {
		MID  string `json:"mid,omitempty"`
		Text string `json:"text,omitempty"`
	}

	ResponseMessage struct {
		MessageType string      `json:"messaging_type"`
		Recipient   *User       `json:"recipient"`
		Message     *ResMessage `json:"message,omitempty"`
		Action      string      `json:"sender_action,omitempty"`
	}

	ResMessage struct {
		Text string `json:"text,omitempty"`
	}
)
