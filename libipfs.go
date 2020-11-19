package main

import "C"

import (
	"fmt"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"unsafe"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

/*
	Internal helper function to convert a *loader.Plugin into a handle
*/
func handle_from_pluginLoader( loader *loader.PluginLoader ) ( int64 ) {
	handle := api_context.next_plugin_loader

	api_context.plugin_loaders[handle] = loader
	api_context.next_plugin_loader = handle + 1

	return handle
}
/*
	Internal helper function to convert a handle to a *loader.PluginLoader
*/
func pluginLoader_from_handle( handle int64 ) (*loader.PluginLoader, int64) {
	loader, ok := api_context.plugin_loaders[handle]
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
//export ipfs_Loader_Release
func ipfs_Loader_Release( handle int64 ) int64 {
	_, ok := api_context.plugin_loaders[handle]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid PluginLoader handle: %d", handle ) ) )
	}

	delete( api_context.plugin_loaders, handle )
	return 1
}


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

//export ipfs_CoreAPI_Unixfs_Get
func ipfs_CoreAPI_Unixfs_Get( api_handle int64, cid_str *C.char ) int64 {
	api, ec := coreAPI_from_handle( api_handle )
	if ec <= 0 {
		return ec
	}

	cid := icorepath.New( C.GoString( cid_str ) )

	node, err := api.Unixfs().Get( api_context.ctx, cid )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	return handle_from_node( node )
}
//export ipfs_CoreAPI_Unixfs_Add
func ipfs_CoreAPI_Unixfs_Add( api_handle int64, node_handle int64 ) int64 {
	api, ec := coreAPI_from_handle( api_handle )
	if ec <= 0 {
		return ec
	}

	node, ec := node_from_handle( node_handle )
	if ec <= 0 {
		return ec
	}

	cid, err := api.Unixfs().Add( api_context.ctx, node )
	if err != nil {
		return ipfs_SubmitError( err )
	}

	err = cid.IsValid()
	if err != nil {
		return ipfs_SubmitError( err )
	}

	return ipfs_SubmitString( cid.String() )
}


func handle_from_node( node files.Node ) ( int64 ) {
	api_context.nodes = append( api_context.nodes, node )
	return int64( len( api_context.nodes ) )
}
func node_from_handle( handle int64 ) (files.Node, int64) {
	// Get the PluginLoader Object
	if handle < 1 || int(handle) > len( api_context.nodes ) {
		return nil,ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid Node handle: %d", handle ) ) )
	}
	return api_context.nodes[handle-1], 1
}

func ipfs_Node_GetType( handle int64 ) int64 {
	node, ec := node_from_handle( handle )
	if ec <= 0 {
		return ec
	}

	switch node.(type) {
		case files.File:
			return 1
		case files.Directory:
			return 2
		default:
			return ipfs_SubmitError( errors.New( "Unknown Node type" ) )
	}
}
//export ipfs_Node_Read
func ipfs_Node_Read( handle int64, unsafe_bytes unsafe.Pointer, bytes_limit int32, offset int64 ) int64 {
	node, ec := node_from_handle( handle )
	if ec <= 0 {
		return ec
	}

	var file files.File
	switch f := node.(type) {
		case files.File:
			file = f
		case files.Directory:
			return 2
			return ipfs_SubmitError( errors.New( "Directory does not support file read" ) )
		default:
			return ipfs_SubmitError( errors.New( "Unknown Node type" ) )
	}

	bytes := make( []byte, bytes_limit )

	_, err := file.Seek( offset, io.SeekStart )
	if err != nil {
		return ipfs_SubmitError( err )
	}
	offset = 0

	read_count, err := file.Read( bytes )
	if err != nil && err != io.EOF {
		return ipfs_SubmitError( err )
	}

	// Copy the data read from the node to the source buffer
	// I couldn't get this to work any other way than the
	// pointer arithmetic below.
	uint_dst := uintptr(unsafe_bytes)
	for i := 0; i < read_count; i++ {
		dst := (*byte)( unsafe.Pointer(uint_dst) )
		*dst = bytes[i]
		uint_dst += 1
	}

	return int64( read_count + 1 )
}
//export ipfs_Node_NewFromPath
func ipfs_Node_NewFromPath( c_path *C.char ) int64 {
	path := C.GoString(c_path)

	st, err := os.Stat(path)
	if err != nil {
		return ipfs_SubmitError(err)
	}

	f, err := files.NewSerialFile( path, false, st )
	if err != nil {
		return ipfs_SubmitError(err)
	}

	return handle_from_node( f )
}


//export ipfs_RunGoroutines
func ipfs_RunGoroutines() int64 {
	runtime.Gosched()
	return 1
}

func main() {}
