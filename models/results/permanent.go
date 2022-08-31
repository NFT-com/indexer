package results

import (
	"strings"
)

func Permanent(err error) bool {

	msg := err.Error()
	switch {

	// any kind of runtime error is a bug and should always hard fail
	case strings.Contains(msg, "runtime error"):
		return true

	// an invalid topic number indicates an error in parsing logic and should hard fail
	case strings.Contains(msg, "invalid number of topics"):
		return true

	// if the message response is too large once, it will always be, so we should hard fail
	// TODO: implement adaptive job size to split them down further
	// => https://github.com/NFT-com/indexer/issues/238
	case strings.Contains(msg, "the message response is too large"):
		return true

	// retrieval of tokens that have been deleted reverts the execution
	// TODO: implement proper handling of smart-contract level deletions
	// => https://github.com/NFT-com/indexer/issues/232
	// TODO: implement proper retrieval of ENS metadata edge case
	// => https://github.com/NFT-com/indexer/issues/226
	case strings.Contains(msg, "execution reverted"):
		return true

	// decoding of unsupported formats fails for now
	// TODO: implement support for object-based attribute pairs
	// => https://github.com/NFT-com/indexer/issues/225
	// TODO: implement support for Enjin metadata format
	// => https://github.com/NFT-com/indexer/issues/237
	case strings.Contains(msg, "cannot unmarshal object into Go struct field Token.attributes of type []metadata.Attribute"):
		return true

	// we can't retrieve metadata if we have no token ID
	// TODO: implement support for Decentraland Registrar metadata edge case
	// => https://github.com/NFT-com/indexer/issues/227
	case strings.Contains(msg, "token URI empty"):
		return true

	// if we run into an unknown token URI format, we should always hard fail
	case strings.Contains(msg, "unknown URI format"):
		return true

	// this should be retried, but we should first properly handle collections
	// where this means it will never work to avoid pipeline clogging
	// TODO: implement proper handling where 500 means deletion
	// => https://github.com/NFT-com/indexer/issues/228
	case strings.Contains(msg, "bad response code (500)"):
		return true

	// this should be retried, we we should first properly handle collections
	// where this means it will never work to avoid pipeline clogging
	// TODO: properly handle IPFS-based deletion of tokens
	// => https://github.com/NFT-com/indexer/issues/230
	case strings.Contains(msg, "bad response code (404)"):
		return true

	default:
		return false
	}
}
