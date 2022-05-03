package processors

import (
	"fmt"
	"net/url"
	"path"
)

const (
	ipfsScheme       = "ipfs"
	resolverHostname = "ipfs.io"
	resolverRootPath = "ipfs"
)

// resolveURI is a temporary helper function to work around IPFS URI addresses.
// IPFS addresses will be replaced with an analogous `ipfs.io` address.
func resolveURI(uri string) (string, error) {

	// FIXME: There should be IPFS-aware parsing here, `net/url` is a little shoehorned.

	address, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("could not parse URI address: %w", err)
	}

	// If we're not working with and IPFS address, were all good.
	if address.Scheme != ipfsScheme {
		return uri, nil
	}

	resourcePath := path.Join(resolverRootPath, address.Host, address.Path)

	resolved := url.URL{
		Scheme: "https",
		Host:   resolverHostname,
		Path:   resourcePath,
	}

	return resolved.String(), nil
}
