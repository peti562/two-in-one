package mocket

func (mh *Helper) Update(data *Data) {

	mh.queryType = "update"

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the INSERT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextUpdate(data.Response, false)
		mh.setupUpdateQuery(data)
		mh.unregister()
	}
}

func (mh *Helper) UpdateWithException(data *Data) {

	mh.queryType = "update"

	// Default the payload to 1
	if data.Times == 0 {
		data.Times = 1
	}

	// Handle the INSERT query
	for i := 0; i < data.Times; i++ {
		mh.catchNextUpdate(data.Response, true)
		mh.setupUpdateQuery(data)
		mh.unregister()
	}
}
