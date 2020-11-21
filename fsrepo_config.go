package main

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
func handle_from_config( config *config.Config ) (int64 ) {
	// Add the config to the object array and return its handle
	api_context.configs = append( api_context.configs, config )
	return int64( len( api_context.configs ) )
}

/*
	Internal helper function to convert a handle into a *config.Config
*/
func config_from_handle( handle int64 ) (*config.Config, int64) {
	// Get the Config  Object
	if handle < 1 || int(handle) > len( api_context.configs ) {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid Config handle: %d", handle ) ) )
	}
	return api_context.configs[ handle - 1 ], 1
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

