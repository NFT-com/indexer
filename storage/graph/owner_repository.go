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

	if len(additions) == 0 {
		return nil
	}

	set := make(map[string]*results.Addition, len(additions))
	for _, addition := range additions {
		key := fmt.Sprintf("%s-%s", addition.NFT.Owner, addition.NFT.ID)
		existing, ok := set[key]
		if ok {
			existing.NFT.Number += addition.NFT.Number
			continue
		}
		set[key] = addition
	}

	additions = make([]*results.Addition, 0, len(set))
	for _, addition := range set {
		additions = append(additions, addition)
	}

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

	if len(modifications) == 0 {
		return nil
	}

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
