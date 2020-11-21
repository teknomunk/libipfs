package main

import "C"

import (
	"fmt"
	"errors"

	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"

	"github.com/ipfs/go-ipfs/core"
)

func handle_from_buildCfg( config *core.BuildCfg ) (int64) {
	handle := api_context.build_cfgs.next_handle

	api_context.build_cfgs.objects[handle] = config
	api_context.build_cfgs.next_handle = handle + 1

	return handle
}
func buildCfg_from_handle( handle int64 ) (*core.BuildCfg,int64) {
	cfg, ok := api_context.build_cfgs.objects[ handle ]
	if ok {
		return cfg, 1
	} else {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid BuildCfg handle: %d", handle ) ) )
	}
}

//export ipfs_BuildCfg_Create
func ipfs_BuildCfg_Create( ) int64 {
	options := &core.BuildCfg {
		Online: true,
	}

	return handle_from_buildCfg( options )
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
	_, ok := api_context.build_cfgs.objects[handle]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid BuildCfg handle %d", handle ) ) )
	}

	// Release the reference and allow the garbage collector to recover
	delete( api_context.configs.objects, handle )

	return 1
}
