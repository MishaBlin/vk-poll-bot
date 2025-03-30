package response

type UpdateResponse struct {
	Update Update `json:"update"`
}

type Update struct {
	Message string   `json:"message"`
	Props   struct{} `json:"props"`
}
