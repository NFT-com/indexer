package graph

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/NFT-com/indexer/models/jobs"
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

func (n *OwnerRepository) Add(additions ...*jobs.Addition) error {

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
			addition.OwnerAddress,
			addition.NFTID(),
			addition.TokenCount,
		)
	}

	_, err := query.Exec()
	if err != nil {
		return fmt.Errorf("could not execute query: %w", err)
	}

	return nil
}

func (o *OwnerRepository) Change(modifications ...*jobs.Modification) error {

	additions := make([]*jobs.Addition, 0, 2*len(modifications))
	for _, modification := range modifications {

		// First we add the count to the new owner.
		addition := jobs.Addition{
			// ID not needed
			ChainID:         modification.ChainID,
			ContractAddress: modification.ContractAddress,
			TokenID:         modification.TokenID,
			// TokenStandard not needed
			OwnerAddress: modification.ReceiverAddress,
			TokenCount:   modification.TokenCount,
		}
		additions = append(additions, &addition)

		// Then we remove it from the old owner.
		removal := addition
		removal.OwnerAddress = modification.SenderAddress
		removal.TokenCount = -removal.TokenCount
		additions = append(additions, &removal)
	}

	err := o.Add(additions...)
	if err != nil {
		return fmt.Errorf("could not add counts: %w", err)
	}

	return nil
}
