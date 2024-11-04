package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/jadd/aurabase/daos"
	"github.com/jadd/aurabase/models"
)

// Resave all view collections to ensure that the proper id normalization is applied.

func init() {
	AppMigrations.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		collections, err := dao.FindCollectionsByType(models.CollectionTypeView)
		if err != nil {
			return nil
		}

		for _, collection := range collections {
			// ignore errors to allow users to adjust
			// the view queries after app start
			dao.SaveCollection(collection)
		}

		return nil
	}, nil)
}
