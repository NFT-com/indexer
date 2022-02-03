package nft

type NFT struct {
	ID       string
	Network  string
	Chain    string
	Contract string
	Owner    string
	Data     map[string]interface{}
}
