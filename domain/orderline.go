package domain

type OrderLine struct {
	OrderID  Reference
	Sku      Sku
	Quantity int
}
