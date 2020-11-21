package main

import (
	"context"
	"fmt"
	"runtime"

	config "github.com/ipfs/go-ipfs-config"
	files "github.com/ipfs/go-ipfs-files"
	icore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/plugin/loader" // This package is needed so that all the preloaded plugins are loaded automatically
)

type ErrorsHolder struct {
	objects			map[int64]error
	next_handle		int64
}
type StringsHolder struct {
	objects			map[int64]string
	next_handle		int64
}
type PluginLoadersHolder struct {
	objects			map[int64]*loader.PluginLoader
	next_handle		int64
}
type UnixfsAddSettingsHolder struct {
	objects			map[int64]*options.UnixfsAddSettings
	next_handle		int64
}
type ConfigsHolder struct {
	objects			map[int64]*config.Config
	next_handle		int64
}
type ReposHolder struct {
	objects			map[int64]repo.Repo
	next_handle		int64
}

/*
 * Stuff all the values that need to be maintained between api calls here
 */
type libipfsAPIContext struct {
	ctx						context.Context
	ctx_cancel				context.CancelFunc

	errors					ErrorsHolder
	strings					StringsHolder
	plugin_loaders				PluginLoadersHolder
	configs					ConfigsHolder
	repos					ReposHolder

	//repos					[]repo.Repo
	build_cfgs				[]*core.BuildCfg
	core_apis					[]icore.CoreAPI
	nodes					[]files.Node

	unixfs_add_settings			UnixfsAddSettingsHolder
}
var api_context libipfsAPIContext

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

	api_context.errors = ErrorsHolder {
		objects: make(map[int64]error),
		next_handle: -1,
	}
	api_context.strings = StringsHolder {
		objects: make(map[int64]string),
		next_handle: 1,
	}
	api_context.plugin_loaders = PluginLoadersHolder {
		objects: make(map[int64]*loader.PluginLoader),
		next_handle: 1,
	}
	api_context.configs = ConfigsHolder {
		objects: make(map[int64]*config.Config),
		next_handle: 1,
	}
	api_context.repos = ReposHolder {
		objects: make(map[int64]repo.Repo),
		next_handle: 1,
	}


	api_context.unixfs_add_settings = UnixfsAddSettingsHolder {
		objects: make(map[int64]*options.UnixfsAddSettings),
		next_handle: 1,
	}
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

//export ipfs_RunGoroutines
func ipfs_RunGoroutines() int64 {
	runtime.Gosched()
	return 1
}

func main() {}

