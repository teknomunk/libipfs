package main

import "C"

import (
	"context"
	"fmt"
	"errors"
	"io"
	"io/ioutil"
//	"log"
//	"os"
//	"path/filepath"
//	"strings"
//	"sync"

	config "github.com/ipfs/go-ipfs-config"
//	files "github.com/ipfs/go-ipfs-files"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
//	icore "github.com/ipfs/interface-go-ipfs-core"
//	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
//	peerstore "github.com/libp2p/go-libp2p-peerstore"
//	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/core"
//	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
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
	configs			[]*config.Config
	repos			[]repo.Repo
	build_cfgs		[]*core.BuildCfg
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





func pluginLoader_from_handle( handle int64 ) (*loader.PluginLoader, error) {
	// Get the PluginLoader Object
	if handle < 1 || int(handle) > len( api_context.plugin_loaders ) {
		api_context.last_error = errors.New( fmt.Sprintf( "Invalid PluginLoader handle %d", handle ) )
		return nil,api_context.last_error
	}
	return api_context.plugin_loaders[handle-1], nil
}

//export ipfs_Loader_PluginLoader_Create
func ipfs_Loader_PluginLoader_Create( plugin_path *C.char ) int64 {
	loader,err := loader.NewPluginLoader( C.GoString( plugin_path ) )

	if err != nil {
		api_context.last_error = err
		return 0
	}

	// Add the loader to the object array and return its handle
	api_context.plugin_loaders = append( api_context.plugin_loaders, loader )
	return int64( len( api_context.plugin_loaders ) )
}

//export ipfs_Loader_PluginLoader_Initialize
func ipfs_Loader_PluginLoader_Initialize( handle int64 ) int64 {
	loader,err := pluginLoader_from_handle( handle )
	if err != nil {
		return 0
	}

	// Load preload and external plugins
	if err := loader.Initialize(); err != nil {
		api_context.last_error = err;
		return 0
	}

	return 1
}

//export ipfs_Loader_PluginLoader_Inject
func ipfs_Loader_PluginLoader_Inject( handle int64 ) int64 {
	loader,err := pluginLoader_from_handle( handle )
	if err != nil {
		return 0
	}

	// Load preload and external plugins
	if err := loader.Inject(); err != nil {
		api_context.last_error = err;
		return 0
	}

	return 1
}



func config_from_handle( handle int64 ) (*config.Config, error) {
	// Get the Config  Object
	if handle < 1 || int(handle) > len( api_context.configs ) {
		api_context.last_error = errors.New( fmt.Sprintf( "Invalid Config handle %d", handle ) )
		return nil, api_context.last_error
	}
	return api_context.configs[ handle - 1 ], nil
}

//export ipfs_Config_Init
func ipfs_Config_Init( io_handle int64, size int32 ) int64 {
	var out io.Writer
	if io_handle == 0 {
		out = ioutil.Discard
	} else {
		api_context.last_error = errors.New( "Invalid IO handle. IO objects not (yet) implemented" )
		return 0
	}

	// Create a config with default options and the specified keysize
	cfg, err := config.Init( out, int(size) )
	if err != nil {
		api_context.last_error = err
		return 0
	}

	// Add the config to the object array and return its handle
	api_context.configs = append( api_context.configs, cfg )
	return int64( len( api_context.configs ) )
}



func repo_from_handle( handle int64 ) (repo.Repo, error) {
	// Get the PluginLoader Object
	if handle < 1 || int(handle) > len( api_context.repos ) {
		api_context.last_error = errors.New( fmt.Sprintf( "Invalid FSREpo handle %d", handle ) )
		return nil,api_context.last_error
	}
	return api_context.repos[handle-1], nil
}

//export ipfs_FSRepo_Init
func ipfs_FSRepo_Init( repo_path *C.char, cfg_handle int64 ) int64 {
	cfg, err := config_from_handle( cfg_handle )
	if err != nil {
		return 0
	}

	if err := fsrepo.Init( C.GoString( repo_path ), cfg ); err != nil {
		api_context.last_error = err
		return 0
	}

	return 1
}

//export ipfs_FSRepo_Open
func ipfs_FSRepo_Open( repo_path *C.char ) int64 {
	repo, err := fsrepo.Open( C.GoString( repo_path ) )
	if err != nil {
		api_context.last_error = err
		return 0
	}

	// Add the repo to the object array and return its handle
	api_context.repos = append( api_context.repos, repo )
	return int64( len( api_context.repos ) )
}



func buildCfg_from_handle( handle int64 ) (*core.BuildCfg,error) {
	if handle < 1 || int(handle) > len( api_context.build_cfgs ) {
		api_context.last_error = errors.New( fmt.Sprintf( "Invalid BuildCfg handle: %d", handle ) )
		return nil,api_context.last_error
	}

	return api_context.build_cfgs[handle-1], nil
}

//export ipfs_BuildCfg_Create
func ipfs_BuildCfg_Create( ) int64 {
	options := &core.BuildCfg {
		Online: true,
	}

	api_context.build_cfgs = append( api_context.build_cfgs, options )
	return int64( len( api_context.build_cfgs ) )
}
//export ipfs_BuildCfg_SetOnline
func ipfs_BuildCfg_SetOnline( handle int64, state int32 ) int64 {
	buildCfg,err := buildCfg_from_handle( handle )
	if err != nil {
		return 0
	}

	if state == 0 {
		buildCfg.Online = false
	} else {
		buildCfg.Online = true
	}
	return 1
}
//export ipfs_BuildCfg_SetRouting
func ipfs_BuildCfg_SetRouting( handle int64, option int32 ) int64 {
	buildCfg,err := buildCfg_from_handle( handle )
	if err != nil {
		return 0
	}

	switch option {
		case 1:
			buildCfg.Routing = libp2p.DHTOption
		case 2:
			buildCfg.Routing = libp2p.DHTClientOption
		default:
			api_context.last_error = errors.New( fmt.Sprintf( "Invalid routing option: %d", option ) )
			return 0
	}

	return 1
}
//export ipfs_BuildCfg_SetRepo
func ipfs_BuildCfg_SetRepo( cfg_handle int64, repo_handle int64 ) int64 {
	buildCfg,err := buildCfg_from_handle( cfg_handle )
	if err != nil {
		return 0
	}

	repo,err := repo_from_handle( repo_handle )
	if err != nil {
		return 0
	}

	buildCfg.Repo = repo

	return 1
}
//export ipfs_BuildCfg_Release
func ipfs_BuildCfg_Release( handle int64 ) int64 {
	if handle < 1 || int(handle) > len( api_context.build_cfgs ) {
		api_context.last_error = errors.New( fmt.Sprintf( "Invalid BuildCfg handle %d", handle ) )
		return 0
	}

	// Release the reference and allog the garbage collector to recover
	api_context.build_cfgs[handle-1] = nil
	return 1
}

func main() {}
