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

// var inputs inputs.OwnerChange
// err := json.Unmarshal(addition.InputData, &inputs)
// if err != nil {
// 	return fmt.Errorf("could not decode owner change inputs: %w", err)
// }

// collection, err := a.collections.One(addition.ChainID, addition.ContractAddress)
// if err != nil {
// 	return fmt.Errorf("could not retrieve collection: %w", err)
// }

// nftHash := sha3.Sum256([]byte(fmt.Sprintf("%d-%s-%s", addition.ChainID, addition.ContractAddress, addition.TokenID)))
// nftID := uuid.Must(uuid.FromBytes(nftHash[:16]))

// err = a.nfts.Touch(nftID.String(), collection.ID, addition.TokenID)
// if err != nil {
// 	return fmt.Errorf("could not touch NFT: %w", err)
// }

// err = a.owners.AddCount(nftID.String(), inputs.PrevOwner, -int(inputs.Number))
// if err != nil {
// 	return fmt.Errorf("could not decrease previous owner count: %w", err)
// }

// err = a.owners.AddCount(nftID.String(), inputs.NewOwner, int(inputs.Number))
// if err != nil {
// 	return fmt.Errorf("could not increase new owner count: %w", err)
// }

func (n *OwnerRepository) Add(addition *jobs.Addition) error {
	return fmt.Errorf("not implemented")
}

func (o *OwnerRepository) Change(modifications ...*jobs.Modification) error {
	return fmt.Errorf("not implemented")
}
