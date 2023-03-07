package opstatus

// Case represents a specific error condition. For example: purchase_limit_exceeded, insufficient_inventory.
type Case interface {
	Identifier() string
}
