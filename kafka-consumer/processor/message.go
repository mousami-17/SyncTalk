package processor

type Message struct {
	Sender     string `json:"sender"`
	SenderName string `json:"senderName"`
	Room       string `json:"room"`
	Type       string `json:"type"`
	Content    string `json:"content"`
	Server     string `json:"server"`
	
	// Media attachment fields
	FileURL   string `json:"fileUrl"`
	FileType  string `json:"fileType"`
	FileName  string `json:"fileName"`
	FileSize  int64  `json:"fileSize"`
	Thumbnail string `json:"thumbnail"`
	Duration  int    `json:"duration"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
}
