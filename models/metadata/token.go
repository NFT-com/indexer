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

	fieldAttributes           = "attributes"
	fieldAttributeDisplayType = "display_type"
	fieldAttributeTraitType   = "trait_type"
	fieldAttributeValue       = "value"
	fieldAttributeTraitValue  = "trait_value"
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
	kv := map[string]interface{}{}
	if err := json.Unmarshal(data, &kv); err != nil {
		return err
	}

	t.Extras = make(map[string]interface{})

	for k, v := range kv {
		var ok bool
		switch k {
		case fieldImage:
			t.Image, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldImage, v)
			}
		case fieldImageData:
			t.ImageData, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldImage, v)
			}
		case fieldExternalURL:
			t.ExternalURL, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldExternalURL, v)
			}
		case fieldDescription:
			t.Description, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldDescription, v)
			}
		case fieldName:
			t.Name, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldName, v)
			}
		case fieldBackgroundColor:
			t.BackgroundColor, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldBackgroundColor, v)
			}
		case fieldAnimationURL:
			t.AnimationURL, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldAnimationURL, v)
			}
		case fieldYoutubeURL:
			t.YoutubeURL, ok = v.(string)
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldYoutubeURL, v)
			}
		case fieldAttributes:
			// Attributes — used in the OpenSea Metadata Standard.
			// See https://docs.opensea.io/docs/metadata-standards#attributes
			attributes, ok := v.([]interface{})
			if !ok {
				return fmt.Errorf("wrong type for field %q: %T", fieldAttributes, v)
			}

			for _, value := range attributes {
				a, ok := value.(map[string]interface{})
				if !ok {
					return fmt.Errorf("wrong type for array element in %q: %T", fieldAttributes, value)
				}

				var attribute Attribute
				for name, value := range a {
					switch name {
					case fieldAttributeDisplayType:
						attribute.DisplayType, ok = value.(string)
						if !ok {
							return fmt.Errorf("wrong type for field %q: %T", "display_type", value)
						}
					case fieldAttributeTraitType:
						attribute.TraitType, ok = value.(string)
						if !ok {
							return fmt.Errorf("wrong type for field %q: %T", "trait_type", value)
						}
					case fieldAttributeValue:
						attribute.Value = value
					case fieldAttributeTraitValue:
						attribute.TraitValue = value
					}
				}

				t.Attributes = append(t.Attributes, attribute)
			}
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
			t.Extras[k] = v
		}
	}

	return nil
}
