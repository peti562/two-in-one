package mocket

type Data struct {
	Model    ModelInterface
	Select   []string
	Update   []string
	Where    []Where
	Response []map[string]interface{}

	// Preloads require table name wrapping
	WrapQuotes bool

	// Used to create the M2M join table
	ManyToMany *ManyToMany

	// Support cache key tests
	CacheKey string

	// How many times should we bind?
	Times int

	// Used to determine if we back tick on .First(x) calls
	Limit int
}

type ManyToMany struct {
	TableName     string
	ForeignKey    string
	ForeignValues []interface{}
	Response      []map[string]interface{}
}

type Where struct {
	Field string
	Value interface{}
}

type ModelInterface interface {
	TableName() string
}
