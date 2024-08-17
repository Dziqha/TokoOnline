package responses

type Order struct {
	Id string `json:"id"`
	UserId string `json:"userId"`
	ProductId string `json:"productId"`
	Quantity int `json:"quantity"`
	TotalPrice int `json:"totalPrice"`
	Status  string `json:"status"`
}