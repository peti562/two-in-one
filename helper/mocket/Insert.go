package mocket

func (mh *Helper) Insert(data *Data) {

	mh.queryType = "insert"

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the INSERT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextInsert(data.Response, false)
		mh.setupInsertQuery(data)
		mh.unregister()
	}
}

func (mh *Helper) InsertWithException(data *Data) {

	mh.queryType = "insert"

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the INSERT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextInsert(data.Response, true)
		mh.setupInsertQuery(data)
		mh.unregister()
	}
}
