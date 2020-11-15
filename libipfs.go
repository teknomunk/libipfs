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
	"runtime"

	config "github.com/ipfs/go-ipfs-config"
//	files "github.com/ipfs/go-ipfs-files"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"
//	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
//	"github.com/libp2p/go-libp2p-core/peer"
)

/*
 * Stuff all the values that need to be maintained between api calls here
 */
type ipfs_api_context struct {
	ctx				context.Context
	ctx_cancel		context.CancelFunc

	// Error handling for asynchronous functions
	errors			map[int64]error
	next_error		int64

	// Handle types
	plugin_loaders		[]*loader.PluginLoader
	configs			[]*config.Config
	repos			[]repo.Repo
	build_cfgs		[]*core.BuildCfg
	core_apis			[]icore.CoreAPI
}
var api_context ipfs_api_context


//export ipfs_Init
func ipfs_Init() {
	// Setup basic enviroment to allow the API to function
	ctx, cancel := context.WithCancel( context.Background() )
	api_context.ctx = ctx
	api_context.ctx_cancel = cancel
	api_context.next_error = -1
}

//export ipfs_Cleanup
func ipfs_Cleanup() {
	// Tear down everything
	api_context.ctx_cancel()

	fmt.Println("ipfs closed")
}

//export ipfs_GetError
func ipfs_GetError( handle int64 ) *C.char {
	err, ok := api_context.errors[handle]
	if ok {
		return C.CString( err.Error() )
	} else {
		return C.CString( "Invalid error handle" )
	}
}
//export ipfs_ReleaseError
func ipfs_ReleaseError( handle int64 ) int64 {
	_, ok := api_context.errors[handle]

	if !ok {
		return 0
	}
	delete( api_context.errors, handle )
	return 1
}

func ipfs_SubmitError( err error ) int64 {
	handle := api_context.next_error

	api_context.errors[handle] = err

	api_context.next_error = handle - 1
	return handle
}




func pluginLoader_from_handle( handle int64 ) (*loader.PluginLoader, int64) {
	// Get the PluginLoader Object
	if handle < 1 || int(handle) > len( api_context.plugin_loaders ) {
		return nil,ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid PluginLoader handle: %d", handle ) ) )
	}
	return api_context.plugin_loaders[handle-1], 1
}

//export ipfs_Loader_PluginLoader_Create
func ipfs_Loader_PluginLoader_Create( plugin_path *C.char ) int64 {
	loader,err := loader.NewPluginLoader( C.GoString( plugin_path ) )

	if err != nil {
		return ipfs_SubmitError( err )
	}

	// Add the loader to the object array and return its handle
	api_context.plugin_loaders = append( api_context.plugin_loaders, loader )
	return int64( len( api_context.plugin_loaders ) )
}

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
	cfg, err := config.Init( out, int(size) )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	// Add the config to the object array and return its handle
	api_context.configs = append( api_context.configs, cfg )
	return int64( len( api_context.configs ) )
}



func repo_from_handle( handle int64 ) (repo.Repo, int64) {
	// Get the PluginLoader Object
	if handle < 1 || int(handle) > len( api_context.repos ) {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid FSREpo handle: %d", handle ) ) )
	}
	return api_context.repos[handle-1], 1
}

//export ipfs_FSRepo_Init
func ipfs_FSRepo_Init( repo_path *C.char, cfg_handle int64 ) int64 {
	cfg, ec := config_from_handle( cfg_handle )
	if ec <= 0 {
		return ec
	}

	if err := fsrepo.Init( C.GoString( repo_path ), cfg ); err != nil {
		return ipfs_SubmitError( err )
	}

	return 1
}

//export ipfs_FSRepo_Open
func ipfs_FSRepo_Open( repo_path *C.char ) int64 {
	repo, err := fsrepo.Open( C.GoString( repo_path ) )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	// Add the repo to the object array and return its handle
	api_context.repos = append( api_context.repos, repo )
	return int64( len( api_context.repos ) )
}



func buildCfg_from_handle( handle int64 ) (*core.BuildCfg,int64) {
	if handle < 1 || int(handle) > len( api_context.build_cfgs ) {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid BuildCfg handle: %d", handle ) ) )
	}

	return api_context.build_cfgs[handle-1], 1
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
	buildCfg,ec := buildCfg_from_handle( handle )
	if ec <= 0 {
		return ec
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
	buildCfg,ec := buildCfg_from_handle( handle )
	if ec <= 0 {
		return ec
	}

	switch option {
		case 1:
			buildCfg.Routing = libp2p.DHTOption
		case 2:
			buildCfg.Routing = libp2p.DHTClientOption
		default:
			return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid routing option: %d", option ) ) )
	}

	return 1
}
//export ipfs_BuildCfg_SetRepo
func ipfs_BuildCfg_SetRepo( cfg_handle int64, repo_handle int64 ) int64 {
	buildCfg,ec := buildCfg_from_handle( cfg_handle )
	if ec <= 0 {
		return ec
	}

	repo,ec := repo_from_handle( repo_handle )
	if ec <= 0 {
		return ec
	}

	buildCfg.Repo = repo

	return 1
}
//export ipfs_BuildCfg_Release
func ipfs_BuildCfg_Release( handle int64 ) int64 {
	if handle < 1 || int(handle) > len( api_context.build_cfgs ) {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid BuildCfg handle %d", handle ) ) )
	}

	// Release the reference and allog the garbage collector to recover
	api_context.build_cfgs[handle-1] = nil
	return 1
}




func coreAPI_from_handle( handle int64 ) (icore.CoreAPI,int64) {
	if handle < 1 || int(handle) > len( api_context.core_apis ) {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid BuildCfg handle: %d", handle ) ) )
	}

	return api_context.core_apis[handle-1], 1
}
//export ipfs_CoreAPI_Create
func ipfs_CoreAPI_Create( cfg_handle int64 ) int64 {
	buildCfg,ec := buildCfg_from_handle( cfg_handle )
	if ec <= 0 {
		return ec
	}

	node, err := core.NewNode( api_context.ctx, buildCfg )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	api, err := coreapi.NewCoreAPI( node )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	api_context.core_apis = append( api_context.core_apis, api )
	return int64( len( api_context.core_apis ) )
}


func perform_swarm_connect( api icore.CoreAPI, addr string, result *int64 ) error {
	maddr, err := ma.NewMultiaddr( addr )
	if err != nil {
		*result = ipfs_SubmitError( err )
		return err
	}

	pii, err := peerstore.InfoFromP2pAddr( maddr )
	if err != nil {
		*result = ipfs_SubmitError( err )
		return err
	}

	pi := peerstore.PeerInfo{ ID: pii.ID }
	pi.Addrs = append( pi.Addrs, pii.Addrs... )

	err = api.Swarm().Connect( api_context.ctx, pi )
	if err != nil {
		*result = ipfs_SubmitError( errors.New(fmt.Sprintf( "failed to connect to %s: %s", pi.ID, err )) )
		return err
	}

	*result = 1
	return nil
}

//export ipfs_CoreAPI_Swarm_Connect_async
func ipfs_CoreAPI_Swarm_Connect_async( api_handle int64, addr *C.char, result *int64 ) int64 {
	api, ec := coreAPI_from_handle( api_handle )
	if ec <= 0 {
		return ec
	}

	go perform_swarm_connect( api, C.GoString( addr ), result )
	return 1
}


//export ipfs_RunGoroutines
func ipfs_RunGoroutines() int64 {
	runtime.Gosched()
	return 1
}

func main() {}
