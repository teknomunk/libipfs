package main

import (
	"C"
	"fmt"
	"errors"
	"io"
	"os"
	"unsafe"

	files "github.com/ipfs/go-ipfs-files"
)

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
