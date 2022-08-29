package metadata

type Attribute struct {
	DisplayType string      `json:"display_type"`
	TraitType   string      `json:"trait_type"`
	Value       interface{} `json:"value"`
	TraitValue  interface{} `json:"trait_value"`
}
