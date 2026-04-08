package inventory

import (
	"context"
	"testing"

	"github.com/adafatya/wms-backend/internal/db/testutil"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	db, querier, cleanup := testutil.SetupTestDB(t)
	defer cleanup()

	repo := NewRepository(querier)
	ctx := context.Background()

	t.Run("CreateProduct", func(t *testing.T) {
		req := CreateProductRequest{
			Name:    "Test Product",
			SKUCode: "SKU-001",
			UOM:     "PCS",
		}

		p, err := repo.CreateProduct(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, p.ID)
		assert.Equal(t, req.Name, p.Name)
		assert.Equal(t, req.SKUCode, p.SKUCode)
	})

	t.Run("CreateLocation", func(t *testing.T) {
		req := CreateLocationRequest{
			Name: "Test Location",
			Code: "LOC-001",
		}

		l, err := repo.CreateLocation(ctx, req)
		assert.NoError(t, err)
		assert.NotEmpty(t, l.ID)
		assert.Equal(t, req.Name, l.Name)
		assert.Equal(t, req.Code, l.Code)
	})

	t.Run("BulkUpsertInventories", func(t *testing.T) {
		// Cleanup first
		testutil.TruncateTables(db)

		p1, _ := repo.CreateProduct(ctx, CreateProductRequest{Name: "P1", SKUCode: "S1", UOM: "PCS"})
		p2, _ := repo.CreateProduct(ctx, CreateProductRequest{Name: "P2", SKUCode: "S2", UOM: "PCS"})
		l1, _ := repo.CreateLocation(ctx, CreateLocationRequest{Name: "L1", Code: "C1"})
		l2, _ := repo.CreateLocation(ctx, CreateLocationRequest{Name: "L2", Code: "C2"})

		productIDs := []int64{p1.ID, p2.ID, p1.ID}
		locationIDs := []int64{l1.ID, l1.ID, l2.ID}
		quantities := []decimal.Decimal{
			decimal.NewFromFloat(10.5),
			decimal.NewFromFloat(20.0),
			decimal.NewFromFloat(5.25),
		}

		err := repo.BulkUpsertInventories(ctx, productIDs, locationIDs, quantities)
		assert.NoError(t, err)

		// Verify P1 in L1
		q1, err := repo.GetInventoryStock(ctx, p1.ID, l1.ID)
		assert.NoError(t, err)
		assert.True(t, decimal.NewFromFloat(10.5).Equal(q1))

		// Update P1 in L1
		err = repo.BulkUpsertInventories(ctx, []int64{p1.ID}, []int64{l1.ID}, []decimal.Decimal{decimal.NewFromFloat(15.0)})
		assert.NoError(t, err)

		q1Updated, _ := repo.GetInventoryStock(ctx, p1.ID, l1.ID)
		assert.True(t, decimal.NewFromFloat(15.0).Equal(q1Updated))
	})

	t.Run("ListInventories", func(t *testing.T) {
		res, pag, err := repo.ListInventories(ctx, 1, 10)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.NotNil(t, pag)
	})

	t.Run("SearchProducts", func(t *testing.T) {
		// P1 and P2 already created in previous test
		res, _, err := repo.ListProducts(ctx, "P1", 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "P1", res[0].Name)

		res2, _, err := repo.ListProducts(ctx, "non-existent", 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res2, 0)
	})

	t.Run("SearchLocations", func(t *testing.T) {
		res, _, err := repo.ListLocations(ctx, "L1", 1, 10)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "L1", res[0].Name)
	})
}
