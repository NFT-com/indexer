package lambdas

var (
	ZeroAddress          = "0x0000000000000000000000000000000000000000"
	OneAddress           = "0x0000000000000000000000000000000000000001"
	DeadSuffixAddress    = "0x000000000000000000000000000000000000dEaD"
	DeadPrefixAddress    = "0xdEAD000000000000000042069420694206942069"
	DeadWeirdAddress     = "0x00000000000000000000045261D4Ee77acdb3286"
	AllOneAddress        = "0x1111111111111111111111111111111111111111"
	AllSixAddress        = "0x6666666666666666666666666666666666666666"
	AllFAddress          = "0xffffffffffffffffffffffffffffffffffffffff"
	CountFromZeroAddress = "0x0123456789012345678901234567890123456789"
	CountFromOneAddress  = "0x1234567890123456789012345678901234567890"
)

func IsMintAddress(address string) bool {
	switch address {
	case ZeroAddress:
		return true
	default:
		return false
	}
}

func IsBurnAddress(address string) bool {
	switch address {
	case ZeroAddress:
		return true
	case OneAddress:
		return true
	case DeadSuffixAddress:
		return true
	case DeadPrefixAddress:
		return true
	case DeadWeirdAddress:
		return true
	case AllOneAddress:
		return true
	case AllSixAddress:
		return true
	case AllFAddress:
		return true
	case CountFromZeroAddress:
		return true
	case CountFromOneAddress:
		return true
	default:
		return false
	}
}
