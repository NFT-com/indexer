package logs

type EventType int

const (
	TypeMint EventType = iota + 1
	TypeTransfer
	TypeBurn
	TypeSale
	TypeURI
)

func (e EventType) String() string {
	switch e {
	case TypeMint:
		return "mint"
	case TypeTransfer:
		return "transfer"
	case TypeBurn:
		return "burn"
	case TypeSale:
		return "sale"
	case TypeURI:
		return "uri"
	default:
		return ""
	}
}
