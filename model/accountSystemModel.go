package model

type Account struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	CreateTime string `json:create_time`
}
