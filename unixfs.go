package main

import "C"

import (
	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
)

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

