package models

type Configuration struct {
	MongoURI        string `json:"MONGO_URI"`
	MongoDB         string `json:"MONGO_DB"`
	MongoCollection string `json:"MONGO_COLLECTION"`
	ServerPort      string `json:"SERVER_PORT"`
	RightMoveAPIURL string `json:"RIGHTMOVE_API_URL"`
}
