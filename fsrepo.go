package main

import "C"

import (
	"fmt"
	"errors"

	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
)

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
