package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/graph"
	"github.com/NFT-com/indexer/models/jobs"
	"github.com/NFT-com/indexer/models/results"
)

type OwnerRepository struct {
	build squirrel.StatementBuilderType
}

func NewOwnerRepository(db *sql.DB) *OwnerRepository {

	cache := squirrel.NewStmtCache(db)
	n := OwnerRepository{
		build: squirrel.StatementBuilder.RunWith(cache).PlaceholderFormat(squirrel.Dollar),
	}

	return &n
}

func (n *OwnerRepository) Add(additions ...*results.Addition) error {

	query := n.build.
		Insert("owners").
		Columns(
			"owner",
			"nft_id",
			"number",
		).
		Suffix("ON CONFLICT (nft_id, owner) DO UPDATE SET " +
			"number = (owners.number + EXCLUDED.number)")

	for _, addition := range additions {
		query = query.Values(
			addition.NFT.Owner,
			addition.NFT.ID,
			addition.NFT.Number,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}

func (o *OwnerRepository) Change(modifications ...*jobs.Modification) error {

	additions := make([]*results.Addition, 0, 2*len(modifications))
	for _, modification := range modifications {

		// First we add the count to the new owner.
		addition := graph.NFT{
			ID:     modification.NFTID(),
			Owner:  modification.ReceiverAddress,
			Number: modification.TokenCount,
		}
		additions = append(additions, &results.Addition{NFT: &addition})

		// Then we remove it from the old owner.
		removal := addition
		removal.Owner = modification.SenderAddress
		removal.Number = -removal.Number
		additions = append(additions, &results.Addition{NFT: &removal})
	}

	err := o.Add(additions...)
	if err != nil {
		return fmt.Errorf("could not add counts: %w", err)
	}

	return nil
}
