package main

import "C"

import (
	"fmt"
	"errors"

	icore "github.com/ipfs/interface-go-ipfs-core"

	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
)

func handle_from_CoreAPI( api icore.CoreAPI ) (int64) {
	handle := api_context.core_apis.next_handle

	api_context.core_apis.objects[handle] = api
	api_context.core_apis.next_handle = handle + 1

	return handle
}
func coreAPI_from_handle( handle int64 ) (icore.CoreAPI,int64) {
	api, ok := api_context.core_apis.objects[handle]
	if ok {
		return api, 1
	} else {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid CoreAPI handle: %d", handle ) ) )
	}
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

	return handle_from_CoreAPI( api )
}

//export ipfs_CoreAPI_Release
func ipfs_CoreAPI_Release( handle int64 ) int64 {
	_, ok := api_context.core_apis.objects[handle]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid CoreAPI handle: %d", handle ) ) )
	}

	delete( api_context.core_apis.objects, handle )

	return 1
}
