package model

type Config struct {
	Token   string `json:"token"`
	WSToken string `json:"ws_token"`
	Lang    string `json:"lang"`
}
