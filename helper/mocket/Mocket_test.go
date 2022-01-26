package mocket

import (
	"os"
	"testing"

	structHelper "two-in-one/helper/struct"

	mocket "github.com/selvatico/go-mocket"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type TestModel struct {
	ID           uint        `gorm:"column:m_id;PRIMARY_KEY"`
	Key          string      `gorm:"column:m_key"`
	Status       bool        `gorm:"column:m_status"`
	ChildModelId uint32      `gorm:"column:fk_child_model_id"`
	ChildModel   SubModel    `gorm:"foriegnKey:ChildModelId;references:ID"`
	SubModels    []*SubModel `gorm:"many2many:model_sub_models;foreignKey:m_id;references:s_id;joinForeignKey:fk_model_id;joinReferences:fk_join_id"`
}

func (t *TestModel) TableName() string {
	return "test_model"
}

type SubModel struct {
	ID uint `gorm:"column:s_id;primary_key"`
}

func (t *SubModel) TableName() string {
	return "test_sub_model"
}

var db *gorm.DB

func TestMain(m *testing.M) {

	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	mocket.Catcher.PanicOnEmptyResponse = true

	// GORM
	db, _ = gorm.Open(Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	// Run the tests
	m.Run()

	// Pass an exit code back
	os.Exit(0)
}

func TestMocketHelper_Select(t *testing.T) {

	// Create our helper object
	mocketHelper := New(db)

	t.Run("Model.First with WHERE", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: uint32(1),
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
			},
			Limit: 1,
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			// Get the first one
			db.First(data.Model, uint32(1))
		})

		mocketHelper.Reset()
	})

	t.Run("Model.Find with WHERE = ?", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: uint32(1),
				},
				{
					Field: "m_id",
					Value: uint32(2),
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
			},
			WrapQuotes: false,
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			// Get the first one
			db.Where("m_id = ?", uint32(1)).
				Where("m_id = ?", uint32(2)).
				Find(data.Model)
		})

		mocketHelper.Reset()
	})

	// Find by key (string)
	t.Run("Model.Find with WHERE IN v1", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_key",
					Value: []interface{}{"One", "Two"},
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&TestModel{
					ID: 2,
				}),
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			// Do a WHERE IN (string, string)
			db.Where("m_key IN (?)", data.Where[0].Value).
				Find(data.Model)
		})

		mocketHelper.Reset()
	})

	// Find by ID (int)
	t.Run("Model.Find with WHERE IN v2", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: []interface{}{1, 2},
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&TestModel{
					ID: 2,
				}),
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			// Do a WHERE IN (int, int)
			db.Where("m_id IN (?)", data.Where[0].Value).
				Find(data.Model)
		})

		mocketHelper.Reset()
	})

	// Throws an exception
	t.Run("Model.Find with Exception", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: []interface{}{10, 20},
				},
			},
			Response: nil,
		}

		// Basic SELECT statement
		mocketHelper.SelectWithException(data)

		// Do a WHERE IN (int, int)
		exception := db.Where("m_id IN (?)", data.Where[0].Value).
			Find(data.Model).
			Error

		// We should get an error back
		assert.NotNil(t, exception)
		assert.Equal(t, "sql error", exception.Error())

		mocketHelper.Reset()
	})

	// Custom SELECT query using a SELECT and WHERE IN (x)
	t.Run("Model.Find with SELECT and WHERE", func(t *testing.T) {

		selectColumns := []string{"m_id", "m_key"}

		// Build the Mocket request
		data := &Data{
			Model:  &TestModel{},
			Select: selectColumns,
			Where: []Where{
				{
					Field: "m_id",
					Value: []interface{}{1, 2},
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&TestModel{
					ID: 2,
				}),
			},
			WrapQuotes: false,
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			// Do a WHERE IN (int, int)
			db.Select("`m_id`,`m_key`").
				Where("m_id IN (?)", data.Where[0].Value).
				Find(data.Model)
		})

		mocketHelper.Reset()
	})

	// Custom SELECT query using a SELECT with JOIN
	t.Run("Model.Find with SELECT and JOIN", func(t *testing.T) {

		modelIds := []interface{}{1, 2}
		subModelIds := []interface{}{5, 10}

		// Build the Mocket request
		modelData := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: modelIds,
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&TestModel{
					ID: 2,
				}),
			},
		}

		// SELECT * FROM `test_model` WHERE m_id IN (1,2)
		mocketHelper.Select(modelData)

		// Build the Mocket request
		subModelData := &Data{
			Model: &SubModel{},
			ManyToMany: &ManyToMany{
				TableName:     "model_sub_models",
				ForeignKey:    "fk_model_id",
				ForeignValues: modelIds,
				Response: []map[string]interface{}{
					{
						"fk_model_id": modelIds[0],
						"fk_join_id":  subModelIds[0],
					},
					{
						"fk_model_id": modelIds[1],
						"fk_join_id":  subModelIds[1],
					},
				},
			},
			Where: []Where{
				{
					Field: "s_id",
					Value: subModelIds,
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&SubModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&SubModel{
					ID: 2,
				}),
			},
			WrapQuotes: true,
		}

		// SELECT * FROM `model_sub_models` WHERE `model_sub_models`.`fk_model_id` IN (1,2)
		// SELECT * FROM `test_sub_model` WHERE s_id IN (5,10)
		mocketHelper.Select(subModelData)

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			var response []*TestModel

			// Do a WHERE IN (int, int)
			db.Preload("SubModels").
				Where("m_id IN (?)", modelData.Where[0].Value).
				Find(&response)
		})

		mocketHelper.Reset()
	})

	// Custom SELECT query with a many2many join
	t.Run("Model.Find with many2many Join", func(t *testing.T) {

		modelIds := []interface{}{1}
		subModelIds := []interface{}{5, 10}

		// Build the Mocket request
		modelData := &Data{
			Model: &TestModel{},
			Where: []Where{
				{
					Field: "m_id",
					Value: modelIds,
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&TestModel{
					ID: 1,
				}),
			},
		}

		// SELECT * FROM `test_model` WHERE m_id = 1
		mocketHelper.Select(modelData)

		// Build the Mocket request
		subModelData := &Data{
			Model: &SubModel{},
			ManyToMany: &ManyToMany{
				TableName:     "model_sub_models",
				ForeignKey:    "fk_model_id",
				ForeignValues: modelIds,
				Response: []map[string]interface{}{
					{
						"fk_model_id": modelIds[0],
						"fk_join_id":  subModelIds[0],
					},
					{
						"fk_model_id": modelIds[0],
						"fk_join_id":  subModelIds[1],
					},
				},
			},
			Where: []Where{
				{
					Field: "s_id",
					Value: subModelIds,
				},
			},
			Response: []map[string]interface{}{
				structHelper.MapAsGorm(&SubModel{
					ID: 1,
				}),
				structHelper.MapAsGorm(&SubModel{
					ID: 2,
				}),
			},
			WrapQuotes: true,
		}

		// SELECT * FROM `model_sub_models` WHERE `model_sub_models`.`fk_model_id` IN (1)
		// SELECT * FROM `test_sub_model` WHERE s_id IN (5,10)
		mocketHelper.Select(subModelData)

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			var response []*TestModel

			// Do a WHERE IN (int, int)
			db.Preload("SubModels").
				Where("m_id IN (?)", modelData.Where[0].Value).
				Find(&response)
		})

		mocketHelper.Reset()
	})

	// Custom SELECT COUNT
	t.Run("Model.Count", func(t *testing.T) {

		var expectedCount int64 = 5

		data := &Data{
			Model:  &TestModel{},
			Select: []string{"count(*)"},
			Response: []map[string]interface{}{
				{
					"count": expectedCount,
				},
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			var value int64

			// Do a WHERE IN (int, int)
			db.Model(data.Model).Count(&value)

			// Did we get our number back?
			assert.Equal(t, expectedCount, value)
		})
	})

	// Custom SELECT COUNT with JOIN
	t.Run("Model.Count with JOIN", func(t *testing.T) {

		var expectedCount int64 = 5

		data := &Data{
			Model:  &TestModel{},
			Select: []string{"count(*)"},
			Response: []map[string]interface{}{
				{
					"count": expectedCount,
				},
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic SELECT statement
			mocketHelper.Select(data)

			var value int64

			// Do a WHERE IN (int, int)
			db.Model(data.Model).
				Joins("ChildModel").
				Joins("JOIN my_table ON my_table.id = fk_table_id").
				Count(&value)

			// Did we get our number back?
			assert.Equal(t, expectedCount, value)
		})
	})

	// Custom INSERT INTO _
	t.Run("Model.Insert", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
			Response: []map[string]interface{}{
				{
					"id": int64(150),
				},
			},
		}

		// What we're going to insert below
		modelObject := &TestModel{
			Key: "test",
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic INSERT statement
			mocketHelper.Insert(data)

			// Do our INSERT statement
			db.Model(data.Model).Create(modelObject)

			// Check the new ID matches
			assert.Equal(t, uint(150), modelObject.ID)
		})

		mocketHelper.Reset()
	})

	// Custom INSERT INTO with exception
	t.Run("Model.Insert with Exception", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{},
		}

		// What we're going to insert below
		modelObject := &TestModel{
			ID:  1,
			Key: "test",
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic INSERT statement
			mocketHelper.InsertWithException(data)

			// Do our INSERT statement
			assert.Error(t, db.Model(data.Model).Create(modelObject).Error)

			mocketHelper.Reset()
		})
	})

	// Custom UPDATE x SET key=value
	t.Run("Model.Save", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{
				ID:     1,
				Key:    "test",
				Status: true,
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic INSERT statement
			mocketHelper.Update(data)

			// Do our UPDATE statement
			assert.NoError(t, db.Save(data.Model).Error)

			mocketHelper.Reset()
		})
	})

	// Custom UPDATE x SET key=value
	t.Run("Model.Save with Update Fields", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{
				ID:     1,
				Key:    "test",
				Status: true,
			},
			Update: []string{"m_key", "m_status"},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic INSERT statement
			mocketHelper.Update(data)

			// Do our UPDATE statement
			assert.NoError(t, db.Save(data.Model).Error)

			mocketHelper.Reset()
		})
	})

	// Custom UPDATE x SET key=value
	t.Run("Model.Save with Exception", func(t *testing.T) {

		data := &Data{
			Model: &TestModel{
				ID:     1,
				Key:    "test",
				Status: true,
			},
		}

		// Make sure Mocket doesn't throw a hissy
		assert.NotPanics(t, func() {

			// Basic update statement
			mocketHelper.UpdateWithException(data)

			// Do our UPDATE statement
			assert.Error(t, db.Model(data.Model).Save(data.Model).Error)

			mocketHelper.Reset()
		})
	})

	t.Run("Reset OK", func(t *testing.T) {

		assert.NotPanics(t, func() {
			mocketHelper.Reset()
		})
	})
}
