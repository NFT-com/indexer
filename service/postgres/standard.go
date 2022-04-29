package postgres

import (
	"fmt"

	"github.com/NFT-com/indexer/models/chain"
)

func (s *Store) Standards(collectionID string) ([]chain.Standard, error) {

	result, err := s.build.
		Select("standards.id, standards.name").
		From("standards_collections, standards").
		Where("standards_collections.collection = ?", collectionID).
		Where("standards_collections.standard = standards.id").
		Query()
	if err != nil {
		return nil, fmt.Errorf("could not query collections: %w", err)
	}
	defer result.Close()

	standards := make([]chain.Standard, 0)
	for result.Next() && result.Err() == nil {
		var standard chain.Standard
		err = result.Scan(
			&standard.ID,
			&standard.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve standards list: %w", err)
		}

		standards = append(standards, standard)
	}

	return standards, nil
}
