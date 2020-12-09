package main

import "C"

import (
	"errors"
	"fmt"
	"strings"

	icorepath "github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/ipfs/interface-go-ipfs-core/options"
	mh "github.com/multiformats/go-multihash"
)

type UnixfsAddOptionArray struct {
	options			[]options.UnixfsAddOption
}

func unixfsAddOptions_from_handle( handle int64 ) (*UnixfsAddOptionArray,int64) {
	opt, ok := api_context.unixfs_add_options.objects[ handle ]
	if ok {
		return opt, 1
	} else {
		return nil, ipfs_SubmitError( errors.New( fmt.Sprintf( "Invalid UnixsfsAddOptions handle: %d", handle ) ) )
	}
}

/*
	Create a new UnixfsAddOptionArray

	Parameters:
		none

	Return:
		A handle to the options array.
*/
//export ipfs_UnixfsAddOptions_Create
func ipfs_UnixfsAddOptions_Create() int64 {
	opts := &UnixfsAddOptionArray{}

	handle := api_context.unixfs_add_options.next_handle
	api_context.unixfs_add_options.next_handle += 1
	api_context.unixfs_add_options.objects[ handle ] = opts

	return handle
}

/*
	Specify option values.

	Parameters:
		int64_t opts_handle
			Handle to the options array created by ipfs_UnixfsAddOptions_Create

		int32_t option
			Option to set. Available options are:

			enum {
				IPFS_UNIXFS_AO_CID_VERSION		= 1,
						// Int parameter
						//   0 - CIDv0	(default)
						//	1 - CIDv1

				IPFS_UNIXFS_AO_HASH_TYPE			= 2,
						// String parameter
						// Acceptable values:
						//   sha-256

				IPFS_UNIXFS_AO_INLINE			= 3,
						// Bool parameter (experimental)
						//  true  - Inline small blocks into CIDs.
						//  false - Normal operation (default)

				IPFS_UNIXFS_AO_INLINE_LIMIT		= 4,
						// Int parameter (experimental)
						// Maximum block size to inline. Default: 32

				IPFS_UNIXFS_AO_RAW_LEAVES		= 5,
						// Bool parameter (experimental)
						//   true  - Use raw blocks for leaf nodes.
						//   false - Normal operation (default)

				IPFS_UNIXFS_AO_RAW_LEAVES_SET		= 6,
						// Bool

				IPFS_UNIXFS_AO_CHUNKER			= 7,
						// String

				IPFS_UNIXFS_AO_PIN				= 9,
						// Bool - Pin added object if true.
			}

		char* c_string_value
			parameter when option is a string value and ignored otherwise

		int64_t int_value
			parameter when option is an integer or boolean value and ignored otherwise

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_UnixfsAddOptions_SetOption
func ipfs_UnixfsAddOption_SetOption( opts_handle int64, option int32, c_string_value *C.char, int_value int64 ) int64 {
	return ipfs_SubmitError( errors.New( fmt.Sprintf( "Unrecognized Unixfs Add option: %d", option)))
}

/*
	Specify the hash function to use. Default if not specified is sha2-256

	Parameters:
		int64_t opts_handle
			Handle to the options array created by ipfs_UnixfsAddOptions_Create

		char* hashStr
			Hash function to use when adding files

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_UnixfsAddOptions_Hash
func ipfs_UnixfsAddOptions_Hash( opts_handle int64, c_hashStr *C.char ) int64 {
	hashStr := C.GoString( c_hashStr )

	opts,ec := unixfsAddOptions_from_handle( opts_handle )
	if ec <= 0 {
		return ec
	}

	hashFunCode, ok := mh.Names[ strings.ToLower( hashStr ) ]
	if !ok {
		return ipfs_SubmitError( errors.New( fmt.Sprintf( "Unrecognized hash function: %s", strings.ToLower(hashStr))))
	}

	opts.options = append( opts.options, options.Unixfs.Hash(hashFunCode) )

	return 1
}

/*
	Inline small blocks into CID (experimental option)

	Parameters:
		int64_t opts_handle
			Handle to the options array created by ipfs_UnixfsAddOptions_Create

		bool inlinea
			Pass TRUE to enable option

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_UnixfsAddOption_Inline
func ipfs_UnixfsAddOption_Inline( opts_handle int64, nline bool ) int64 {
	opts,ec := unixfsAddOptions_from_handle( opts_handle )
	if ec <= 0 {
		return ec
	}
	opts.options = append( opts.options, options.Unixfs.Inline(nline) )

	return 1
}
/*
	Specify maximum block size to inline. Default is 32 bytes (experimental option)

	Parameters:
		int64_t opts_handle
			Handle to the options array created by ipfs_UnixfsAddOptions_Create

		bool inlinea
			Pass TRUE to enable option

	Return:
		On success:
			1
		Otherwise:
			Error code. See ipfs_GetError and ipfs_ReleaseError for details
*/
//export ipfs_UnixfsAddOption_InlineLimit
func ipfs_UnixfsAddOption_InlineLimit( opts_handle int64, inlineLimit int64 ) int64 {
	opts,ec := unixfsAddOptions_from_handle( opts_handle )
	if ec <= 0 {
		return ec
	}
	opts.options = append( opts.options, options.Unixfs.InlineLimit(int(inlineLimit)) )

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

