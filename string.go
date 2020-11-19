package main

import "C"

/*
	Get a string representation for an error code. This library treats any
	negative return code as an error condition. The code returned is specific
	to the function call producing it. Pass the error code to this function
	to get a string representation of the error.
*/
//export ipfs_GetError
func ipfs_GetError( handle int64 ) *C.char {
	err, ok := api_context.errors[handle]
	if ok {
		return C.CString( err.Error() )
	} else {
		return C.CString( "Invalid error handle" )
	}
}

/*
	Release the error code. This must be called for every error code received
	or a memory leak will result.
*/
//export ipfs_ReleaseError
func ipfs_ReleaseError( handle int64 ) int64 {
	_, ok := api_context.errors[handle]

	if !ok {
		return 0
	}
	delete( api_context.errors, handle )
	return 1
}

/*
	Internal helper function to take a Go error and turn it into an error
	code for returning thru the API.
*/
func ipfs_SubmitError( err error ) int64 {
	handle := api_context.next_error

	api_context.errors[handle] = err
	api_context.next_error = handle - 1

	return handle
}

/*
	Internal helper function to take a string and turn it into a return code.
	Because C can only return a single return value, this is used to turn a
	string into an integer handle that can be returned by a function with a int64_t
	return type or thru a inter64_t* parameter return buffer.
*/
func ipfs_SubmitString( str string ) int64 {
	handle := api_context.next_string

	api_context.strings[handle] = str
	api_context.next_error = handle + 1

	return handle
}

/*
	Get a C-string for a given string handle.
*/
//export ipfs_GetString
func ipfs_GetString( handle int64 ) *C.char {
	str, ok := api_context.strings[handle]
	if ok {
		return C.CString( str )
	} else {
		return nil
	}
}

/*
	Release a string handle. Failing to release a string returned from an API
	call will result in a memory leak.
*/
//export ipfs_ReleaseString
func ipfs_ReleaseString( handle int64 ) int64 {
	_, ok := api_context.strings[handle]

	if !ok {
		return 0
	}
	delete( api_context.strings, handle )
	return 1
}
