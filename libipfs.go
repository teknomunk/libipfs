package main

import "C"

import (
	"context"
	"fmt"
//	"io/ioutil"
//	"log"
//	"os"
//	"path/filepath"
//	"strings"
//	"sync"

//	config "github.com/ipfs/go-ipfs-config"
//	files "github.com/ipfs/go-ipfs-files"
//	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
//	icore "github.com/ipfs/interface-go-ipfs-core"
//	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
//	peerstore "github.com/libp2p/go-libp2p-peerstore"
//	ma "github.com/multiformats/go-multiaddr"

//	"github.com/ipfs/go-ipfs/core"
//	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
//	"github.com/ipfs/go-ipfs/repo/fsrepo"
//	"github.com/libp2p/go-libp2p-core/peer"
)

/*
 * Stuff all the values that need to be maintained between api calls here
 */
type ipfs_api_context struct {
	ctx			context.Context
	ctx_cancel		context.CancelFunc
	last_error		error

	// Handle types
	plugin_loaders		[]*loader.PluginLoader
}
var api_context ipfs_api_context

//export ipfs_Init
func ipfs_Init() {
	// Setup basic enviroment to allow the API to function
	ctx, cancel := context.WithCancel( context.Background() )
	api_context.ctx = ctx
	api_context.ctx_cancel = cancel
}

//export ipfs_Cleanup
func ipfs_Cleanup() {
	// Tear down everything
	api_context.ctx_cancel()

	fmt.Println("ipfs closed")
}

//export ipfs_LastError
func ipfs_LastError() *C.char {
	return C.CString( api_context.last_error.Error() )
}

//export ipfs_Loader_PluginLoader_Create
func ipfs_Loader_PluginLoader_Create( plugin_path *C.char ) int64 {
	loader,err := loader.NewPluginLoader( C.GoString( plugin_path ) )

	if err != nil {
		api_context.last_error = err
		return 0
	}

	// Add the loader to the object array
	api_context.plugin_loaders = append( api_context.plugin_loaders, loader )

	// Return the position in the object array + 1
	return int64( len( api_context.plugin_loaders ) )
}

//export ipfs_Loader_PluginLoader_Initialize
func ipfs_Loader_PluginLoader_Initialize( handle int64 ) int64 {
	fmt.Println( "ipfs_Loader_PluginLoader_Initialize not yet implemented" )
	return 0
}

//export ipfs_Loader_PluginLoader_Inject
func ipfs_Loader_PluginLoader_Inject( handle int64 ) int64 {
	fmt.Println( "ipfs_Loader_PluginLoader_Inject not yet implemented" )
	return 0
}

func main() {}
