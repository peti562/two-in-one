package mocket

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	mocket "github.com/selvatico/go-mocket"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type Helper struct {
	gormDb *gorm.DB

	queryType    string
	waitingMocks map[string]int
}

// New returns an instance of our mocket Helper
func New(gormDb *gorm.DB) *Helper {
	return &Helper{
		gormDb:       gormDb,
		waitingMocks: make(map[string]int),
	}
}

func (mh *Helper) Reset() {
	output := ""
	for query, count := range mh.waitingMocks {
		if count <= 0 {
			continue
		}
		output += fmt.Sprintf("%s was expected %d more time", query, count)
		if count > 1 {
			output += "s"
		}
		output += "\n"
	}
	if len(output) > 0 {
		if os.Getenv("MOCKET_PANIC") == "false" {
			fmt.Println("[Mocket] Error:", output)
		} else {
			panic("[Mocket] Error: " + output)
		}
	}
	mh.waitingMocks = make(map[string]int)
	mocket.Catcher.Reset()
}

func (mh *Helper) getHookName() string {

	hookName := "gorm:query"

	switch mh.queryType {
	case "insert":
		hookName = "gorm:create"
	case "update":
		hookName = "gorm:update"
	}

	return hookName
}

func (mh *Helper) catchNextInsert(withResponse []map[string]interface{}, withException bool) {

	// Bind onto gorm:create to create the INSERT statement
	_ = mh.gormDb.Callback().
		Create().
		Before(mh.getHookName()).
		Register("test:query", func(scope *gorm.DB) {

			// Get the SQL string
			scope.Statement.SQL.Grow(180)
			scope.Statement.AddClauseIfNotExists(clause.Insert{})
			if set := callbacks.ConvertToAssignments(scope.Statement); len(set) != 0 {
				scope.Statement.AddClause(set)
			}
			scope.Statement.Build("INSERT")

			// Build the SQL string
			// LT 08/07/2021, this no longer works because ... well ... Gorm?
			// callbacks.BuildQuerySQL(scope)

			// Get the prepared query string
			queryString := scope.Statement.SQL.String()

			// Find all count statements
			queryString = strings.ReplaceAll(queryString, "count(1)", "count(*)")
			queryString = strings.ReplaceAll(queryString, "SELECT * FROM", "INSERT INTO")

			// Handle the query payload
			mh.handleQuery(queryString, withResponse, withException)

			// Prevent the query running Gorm v1
			scope.InstanceSet("gorm:skip_query_callback", true)

			// Prevent the query running Gorm v2
			scope.DryRun = true
		})
}

func (mh *Helper) catchNextQuery(withResponse []map[string]interface{}, withException bool) {

	// Bind onto gorm:query to create the SELECT statement
	_ = mh.gormDb.Callback().
		Query().
		Before(mh.getHookName()).
		Register("test:query", func(scope *gorm.DB) {

			// Build the SQL string
			callbacks.BuildQuerySQL(scope)

			// Get the prepared query string
			queryString := scope.Statement.SQL.String()

			// Find all count statements
			queryString = strings.ReplaceAll(queryString, "count(1)", "count(*)")

			// Handle the query payload
			mh.handleQuery(queryString, withResponse, withException)

			// Prevent the query running
			scope.InstanceSet("gorm:skip_query_callback", true)

			// Prevent the query running Gorm v2
			scope.DryRun = true
		})
}

func (mh *Helper) catchNextUpdate(withResponse []map[string]interface{}, withException bool) {

	// Bind onto gorm:update to create the INSERT statement
	_ = mh.gormDb.Callback().
		Update().
		Before(mh.getHookName()).
		Register("test:query", func(scope *gorm.DB) {

			// Get the SQL string
			scope.Statement.SQL.Grow(180)
			scope.Statement.AddClauseIfNotExists(clause.Update{})
			if set := callbacks.ConvertToAssignments(scope.Statement); len(set) != 0 {
				scope.Statement.AddClause(set)
			}
			scope.Statement.Build("UPDATE")

			// Get the prepared query string
			queryString := scope.Statement.SQL.String()

			// Find all count statements
			queryString = strings.ReplaceAll(queryString, "count(1)", "count(*)")
			queryString = strings.ReplaceAll(queryString, "SELECT * FROM", "UPDATE")

			// Handle the query payload
			mh.handleQuery(queryString, withResponse, withException)

			// Prevent the query running Gorm v1
			scope.InstanceSet("gorm:skip_query_callback", true)

			// Prevent the query running Gorm v2
			scope.DryRun = true

			// Meant to skip the SELECT callbacks
			scope.Set("gorm:update_column", true)
			scope.Set("gorm:started_transaction", false)
		})
}

func (mh *Helper) handleQuery(queryString string, withResponse []map[string]interface{}, withException bool) {

	// Build our query string
	fmt.Printf("Handling Query: '%s' \n", queryString)

	// Handle the mocket select
	mockObject := mocket.Catcher.
		NewMock().
		OneTime().
		WithQuery(queryString)

	// If we want an exception
	if withException {
		mockObject.WithError(errors.New("sql error"))
		fmt.Println("- Returning an exception")
	} else {
		mh.waitingMocks[queryString]++
		mockObject.WithCallback(func(query string, values []driver.NamedValue) {
			fmt.Println("Handled Query ", query)
			mh.waitingMocks[queryString]--
		})
	}

	// INSERT statements, set the ID
	if mh.queryType == "insert" {

		// Make sure we have a response object
		if len(withResponse) > 0 {
			if idValue, hasKey := withResponse[0]["id"]; hasKey {
				mockObject.WithID(idValue.(int64))
				fmt.Println("- Return the insert ID:", idValue)
			}
		} else {
			mockObject.WithID(int64(1))
		}
	}

	// UPDATE statements, handle the rows affected
	if mh.queryType == "update" {

		// Make sure we have a response object
		if len(withResponse) > 0 {
			if idValue, hasKey := withResponse[0]["count"]; hasKey {
				mockObject.WithRowsNum(idValue.(int64))
				fmt.Println("- Return the affected rows:", idValue)
			}
		} else {
			mockObject.WithRowsNum(int64(1))
		}
	}

	// We have a custom response payload to give
	if withResponse != nil {
		mockObject.WithReply(withResponse)
		fmt.Println("- Returning a payload")
	}
}

func (mh *Helper) unregister() {

	switch mh.queryType {
	case "insert":
		_ = mh.gormDb.Callback().Create().Before(mh.getHookName()).Remove("test:query")
	case "update":
		_ = mh.gormDb.Callback().Update().Before(mh.getHookName()).Remove("test:query")
	case "select":
		_ = mh.gormDb.Callback().Query().Before(mh.getHookName()).Remove("test:query")
	}
}

func (mh *Helper) setupManyToManyQuery(data *ManyToMany) {

	// Create a fake model so that we can call .Find(x)
	type fakeModel struct {
		ID int
	}

	var response []*fakeModel

	// We're going to abstract the slice away
	var values interface{} = data.ForeignValues

	// Turn our values into a single value
	if len(data.ForeignValues) == 1 {
		values = data.ForeignValues[0]
	}

	tx := mh.gormDb.Begin()

	// Select our table with the right data
	tx = tx.Table(data.TableName)

	// Parse our model, get the schema definition etc.
	_ = tx.Statement.Parse(tx.Statement.Model)

	// Handle the WHERE query
	// We use a limit of -1 to mean ALL
	tx = generateWhere(tx, data.TableName, data.ForeignKey, values, true, -1)

	// Get the response payload
	tx.Find(&response)
}

func (mh *Helper) setupInsertQuery(data *Data) {

	tx := mh.gormDb.Begin()

	// Simulate the query
	tx.Model(data.Model).
		Create(data.Model)
}

func (mh *Helper) setupSelectQuery(data *Data) {

	tx := mh.gormDb.Begin()

	// Load in our model
	tx = tx.Model(data.Model)

	// Parse our model, get the schema definition etc.
	_ = tx.Statement.Parse(tx.Statement.Model)

	// Custom SELECT columns?
	if len(data.Select) > 0 {
		tx = tx.Select(data.Select)
	}

	// Loop over the where parts
	for _, wherePart := range data.Where {
		tx = generateWhere(tx, data.Model.TableName(), wherePart.Field, wherePart.Value, data.WrapQuotes, data.Limit)
	}

	// For testing the cache helper
	if len(data.CacheKey) > 0 {
		tx = tx.Set("cache:key", data.CacheKey)
	}

	// Simulate the query
	tx.Find(data.Model)
}

func (mh *Helper) setupUpdateQuery(data *Data) {

	tx := mh.gormDb.Begin()

	// Make sure we load in our model
	tx = tx.Model(data.Model)

	// Parse our model, get the schema definition etc.
	_ = tx.Statement.Parse(tx.Statement.Model)

	// Loop over the where parts
	for _, wherePart := range data.Where {
		tx = generateWhere(tx, data.Model.TableName(), wherePart.Field, wherePart.Value, data.WrapQuotes, data.Limit)
	}

	// We have custom fields to update
	if len(data.Update) > 0 {

		// Build our update map
		updateData := make(map[string]interface{})

		// The fields we want to update
		for _, field := range data.Update {
			updateData[field] = 1
		}

		tx.Model(data.Model).Updates(updateData)
		return
	}

	// Simulate the query
	tx.Save(data.Model)
}

func generateWhere(tx *gorm.DB, tableName, column string, wherePart interface{}, wrapQuotes bool, limit int) *gorm.DB {

	// By default, we're just the column name
	columnKey := column
	operator := "= "

	// What value do we have?
	kindOf := reflect.TypeOf(wherePart).Kind()

	if kindOf == reflect.Slice || kindOf == reflect.Array {
		operator = ""
		whereString := whereConvertString(wherePart.([]interface{}), false)
		wherePart = fmt.Sprintf("IN (%v)", whereString)
	}

	// If we're the PK and singular, we use the column wrapping
	if tx.Statement.Schema != nil {

		// Get the limit
		limitValue := getLimit(tx)

		// Use the one passed in
		if limitValue == 0 {
			limitValue = limit
		}

		// If we've called 'First' we'll have `LIMIT 1` against our PK field - wrap it
		if inArray(column, tx.Statement.Schema.PrimaryFieldDBNames) {

			// Determine if we're a slice, if we aren't and we're a single value, we use the wrapped version
			// as if it was generated using .First(&x, X)
			if kindOf != reflect.Slice && kindOf != reflect.Array && limitValue == 1 {
				wrapQuotes = true
			}
		}
	}

	// We want to wrap the column in quotes
	if wrapQuotes {
		columnKey = fmt.Sprintf("`%s`.`%s`", tableName, column)
	}

	// If we're a custom passed in .Where, we don't?!
	tx = tx.Where(fmt.Sprintf(`%s %s%v`, columnKey, operator, wherePart))

	return tx
}

func whereConvertString(a []interface{}, withSpace bool) string {
	str := ""
	for index := 0; index < len(a); index++ {

		// Get the objects underlying value
		str += fmt.Sprintf("%v", a[index])

		// Anything but the last one
		if index != (len(a) - 1) {
			str += ","
			if withSpace {
				str += " "
			}
		}
	}
	return str
}

func inArray(key string, data []string) bool {
	for _, checkKey := range data {
		if key == checkKey {
			return true
		}
	}
	return false
}

func getLimit(scope *gorm.DB) int {

	limitValue := 0

	// Read the LIMIT clause, Gorm v2
	if statementClause, isOK := scope.Statement.Clauses["LIMIT"]; isOK {
		limitClause := statementClause.Expression.(clause.Limit)
		limitValue = limitClause.Limit
	}

	return limitValue
}
