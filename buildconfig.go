package main

import "C"

import (
	"fmt"
	"errors"

	libp2p "github.com/ipfs/go-ipfs/core/node/libp2p"

	"github.com/ipfs/go-ipfs/core"
)

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
