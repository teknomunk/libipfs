package main

import "C"

import (
	"fmt"
	"errors"
	"io"
	"io/ioutil"

	config "github.com/ipfs/go-ipfs-config"
)

/*
	Internal helper function to convert a *config.Config into a handle
*/
func handle_from_config( config *config.Config ) (int64) {
	// Add the config to the object array and return its handle
	handle := api_context.configs.next_handle

	api_context.configs.objects[handle] = config
	api_context.configs.next_handle = handle + 1

	return handle
}

/*
	Internal helper function to convert a handle into a *config.Config
*/
func config_from_handle( handle int64 ) (*config.Config, int64) {
	// Get the Config  Object
	config, ok := api_context.configs.objects[handle]
	if ok {
		return config, 1
	} else {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid Config handle: %d", handle ) ) )
	}
}

//export ipfs_Config_Init
func ipfs_Config_Init( io_handle int64, size int32 ) int64 {
	var out io.Writer
	if io_handle == 0 {
		out = ioutil.Discard
	} else {
		return ipfs_SubmitError( errors.New( "Invalid IO handle. IO objects not (yet) implemented" ) )
	}

	// Create a config with default options and the specified keysize
	config, err := config.Init( out, int(size) )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	return handle_from_config( config )
}
//export ipfs_Config_Release
func ipfs_Config_Release( handle int64 ) int64 {
	_, ok := api_context.configs.objects[handle]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid Config handle: %d", handle ) ) )
	}

	delete( api_context.configs.objects, handle )
	return 1
}

