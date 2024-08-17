package responses

type Carts struct {
	Id string `json:"id"`
	UserId string `json:"userId"`
	ProductId string `json:"productId"`
	Quantity int `json:"quantity"`
}