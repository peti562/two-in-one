package mocket

func (mh *Helper) Select(data *Data) {

	mh.queryType = "select"

	// Do we have a many to many join?
	if data.ManyToMany != nil {

		// Join statements need to be explicitly wrapped in quotes
		data.WrapQuotes = true

		// Handle the M2M2 query
		mh.catchNextQuery(data.ManyToMany.Response, false)
		mh.setupManyToManyQuery(data.ManyToMany)
		mh.unregister()
	}

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the SELECT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextQuery(data.Response, false)
		mh.setupSelectQuery(data)
		mh.unregister()
	}
}

func (mh *Helper) SelectWithException(data *Data) {

	mh.queryType = "select"

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the SELECT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextQuery(data.Response, true)
		mh.setupSelectQuery(data)
		mh.unregister()
	}
}
