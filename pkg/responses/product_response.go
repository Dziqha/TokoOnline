package responses

type Product struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Price int `json:"price"`
	Stock int `json:"stock"`
}