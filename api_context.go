package main

import (
	"context"
	"fmt"
//	"errors"
//	"io"
//	"io/ioutil"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
//	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"
	icore "github.com/ipfs/interface-go-ipfs-core"
//	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
//	peerstore "github.com/libp2p/go-libp2p-peerstore"
//	ma "github.com/multiformats/go-multiaddr"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/core"
//	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
//	"github.com/ipfs/go-ipfs/repo/fsrepo"
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

	// String passing
	strings			map[int64]string
	next_string		int64

	// PluginLoader map
	//plugin_loaders	[]*loader.PluginLoader
	plugin_loaders		map[int64]*loader.PluginLoader
	next_plugin_loader	int64

	configs			[]*config.Config
	repos			[]repo.Repo
	build_cfgs		[]*core.BuildCfg
	core_apis			[]icore.CoreAPI
	nodes			[]files.Node
}
var api_context ipfs_api_context

/*
	Initialize the IPFS library. This must be called before calling any other
	ipfs_* functions.
*/
//export ipfs_Init
func ipfs_Init() {
	// Setup basic enviroment to allow the API to function
	ctx, cancel := context.WithCancel( context.Background() )
	api_context.ctx = ctx
	api_context.ctx_cancel = cancel

	api_context.errors = make(map[int64]error)
	api_context.next_error = -1

	api_context.strings = make(map[int64]string)
	api_context.next_string = 1

	api_context.plugin_loaders = make(map[int64]*loader.PluginLoader)
	api_context.next_plugin_loader = 1
}
/*
	Cleanup/teardown the IPFS library. Calling any ipfs_* functions after
	calling this function is undefined.
*/
//export ipfs_Cleanup
func ipfs_Cleanup() {
	// Tear down everything
	api_context.ctx_cancel()

	fmt.Println("ipfs closed")
}

