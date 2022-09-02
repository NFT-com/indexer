package metadata

import (
	"encoding/json"
	"fmt"
)

const (
	fieldImage           = "image"
	fieldImageData       = "image_data"
	fieldExternalURL     = "external_url"
	fieldDescription     = "description"
	fieldName            = "name"
	fieldBackgroundColor = "background_color"
	fieldAnimationURL    = "animation_url"
	fieldYoutubeURL      = "youtube_url"
	fieldProperties      = "properties"
	fieldAttributes      = "attributes"
)

type Token struct {
	Image           string `json:"image"`
	ImageData       string `json:"image_data"`
	ExternalURL     string `json:"external_url"`
	Description     string `json:"description"`
	Name            string `json:"name"`
	BackgroundColor string `json:"background_color"`
	AnimationURL    string `json:"animation_url"`
	YoutubeURL      string `json:"youtube_url"`

	// Combination of Attributes from the OpenSea standard and
	// converted Properties from the ERC-1155 standard schema.
	Attributes []Attribute `json:"attributes"`

	// The Extras map contains any arbitrary metadata that is
	// not covered by the OpenSea standard at root level.
	Extras map[string]interface{} `json:"-"`
}

func (t *Token) UnmarshalJSON(data []byte) error {

	// Since calling UnmarshalJSON on a Token would cause an infinite loop,
	// use a type alias to automatically unmarshal fields that do not require
	// any special attention.
	type Alias Token
	var auxiliary Alias
	if err := json.Unmarshal(data, &auxiliary); err != nil {
		return err
	}
	t.Image = auxiliary.Image
	t.ImageData = auxiliary.ImageData
	t.ExternalURL = auxiliary.ExternalURL
	t.Description = auxiliary.Description
	t.Name = auxiliary.Name
	t.BackgroundColor = auxiliary.BackgroundColor
	t.AnimationURL = auxiliary.AnimationURL
	t.YoutubeURL = auxiliary.YoutubeURL
	t.Attributes = auxiliary.Attributes
	t.Extras = make(map[string]interface{})

	// Unmarshal into a map[string]interface to look for the properties field if it exists,
	// as well as any extra fields at the root level that are not part of the supported
	// standards.
	kv := map[string]interface{}{}
	if err := json.Unmarshal(data, &kv); err != nil {
		return err
	}

	for k, v := range kv {
		switch k {
		case fieldImage, fieldImageData, fieldExternalURL, fieldDescription, fieldName, fieldBackgroundColor, fieldAnimationURL, fieldYoutubeURL, fieldAttributes:
			// Nothing to do here, those fields are handled by the auxiliary struct automatically.
			continue

		case fieldProperties:
			// Properties — used by the ERC-1155 metadata schema.
			// See https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1155.md#erc-1155-metadata-uri-json-schema
			// This schema is a mere *suggestion*, so it is not thoroughly followed by implementations.
			properties, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldProperties, v)
			}

			for name, value := range properties {
				attribute := Attribute{
					TraitType: name,
					Value:     value,
				}
				t.Attributes = append(t.Attributes, attribute)
			}

		default:
			// This field is not part of the supported standards, therefore it is added
			// to the Extras map.
			t.Extras[k] = v
		}
	}

	return nil
}
