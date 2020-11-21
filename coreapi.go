package main

import "C"

import (
	"fmt"
	"errors"

	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
)

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
