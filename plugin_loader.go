package main

import "C"

import (
	"errors"
	"fmt"

	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
)

/*
	Internal helper function to convert a *loader.Plugin into a handle
*/
func handle_from_pluginLoader( loader *loader.PluginLoader ) ( int64 ) {
	handle := api_context.plugin_loaders.next_handle

	api_context.plugin_loaders.objects[handle] = loader
	api_context.plugin_loaders.next_handle = handle + 1

	return handle
}
/*
	Internal helper function to convert a handle to a *loader.PluginLoader
*/
func pluginLoader_from_handle( handle int64 ) (*loader.PluginLoader, int64) {
	loader, ok := api_context.plugin_loaders.objects[handle]
	if ok {
		return loader, 1
	} else {
		return nil,ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid PluginLoader handle: %d", handle ) ) )
	}
}

/*
	Create a new *loader.PluginLoader.

	Parameters:
		plugin_path
			Path to plugin path.  May be "".

	Return:
		On success:
			A handle to the plugin loader.  num > 0
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_Loader_PluginLoader_Create
func ipfs_Loader_PluginLoader_Create( plugin_path *C.char ) int64 {
	loader,err := loader.NewPluginLoader( C.GoString( plugin_path ) )

	if err != nil {
		return ipfs_SubmitError( err )
	}

	// Add the loader to the object array and return its handle
	return handle_from_pluginLoader( loader )
}

/*
	Initialize a plugin loader

	Parameters:
		handle
			Handle to plugin loader; the result of a previous
			ipfs_Loader_PluginLoader call

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_Loader_PluginLoader_Initialize
func ipfs_Loader_PluginLoader_Initialize( handle int64 ) int64 {
	loader,ec := pluginLoader_from_handle( handle )
	if ec <= 0 {
		return ec
	}

	// Load preload and external plugins
	if err := loader.Initialize(); err != nil {
		return ipfs_SubmitError( err )
	}

	return 1
}

/*
	Inject the plugins loaded into the IPFS node

	Parameters:
		handle
			Handle to plugin loader; the result of a previous
			ipfs_Loader_PluginLoader call

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_Loader_PluginLoader_Inject
func ipfs_Loader_PluginLoader_Inject( handle int64 ) int64 {
	loader,ec := pluginLoader_from_handle( handle )
	if ec <= 0 {
		return ec
	}

	// Load preload and external plugins
	if err := loader.Inject(); err != nil {
		return ipfs_SubmitError( err )
	}

	return 1
}

/*
	Release the plugin loader object

	Parameters:
		handle
			Handle to plugin loader; the result of a previous
			ipfs_Loader_PluginLoader call

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_Loader_PluginLoader_Release
func ipfs_Loader_PluginLoader_Release( handle int64 ) int64 {
	_, ok := api_context.plugin_loaders.objects[handle]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid PluginLoader handle: %d", handle ) ) )
	}

	delete( api_context.plugin_loaders.objects, handle )
	return 1
}

