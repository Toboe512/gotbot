package telegram

type UpdatesResponse struct {
	Ok          bool               `json:"ok"`
	Description string             `json:"description"`
	ErrorCode   int                `json:"error_code"`
	Result      []Update           `json:"result"`
	Parameters  ResponseParameters `json:"parameters"`
}

type ResponseParameters struct {
	MigrateToChatId int `json:"migrate_to_chat_id"`
	RetryAfter      int `json:"retry_after"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"message"`
}

type Message struct {
	MessageID            int         `json:"message_id"`
	MessageThreadID      int         `json:"message_thread_id"`
	From                 User        `json:"from"`
	SenderChat           Chat        `json:"sender_chat"`
	SanderBusinessBot    User        `json:"sander_business_bot"`
	Date                 int         `json:"date"`
	BusinessConnectionID string      `json:"business_connection_id"`
	Chat                 Chat        `json:"chat"`
	Text                 string      `json:"text"`
	Photo                []PhotoSize `json:"photo"`
	Caption              string      `json:"caption"`
}

type PhotoSize struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     int    `json:"file_size"`
}

type User struct {
	ID        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"username"`
}

type Chat struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsForum   bool   `json:"is_forum"`
}
