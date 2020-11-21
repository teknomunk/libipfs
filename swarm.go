package main

import "C"

import (
	"fmt"
	"errors"

	icore "github.com/ipfs/interface-go-ipfs-core"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	ma "github.com/multiformats/go-multiaddr"
)

func perform_swarm_connect( api icore.CoreAPI, addr string, result *int64 ) error {
	maddr, err := ma.NewMultiaddr( addr )
	if err != nil {
		*result = ipfs_SubmitError( err )
		return err
	}

	pii, err := peerstore.InfoFromP2pAddr( maddr )
	if err != nil {
		*result = ipfs_SubmitError( err )
		return err
	}

	pi := peerstore.PeerInfo{ ID: pii.ID }
	pi.Addrs = append( pi.Addrs, pii.Addrs... )

	err = api.Swarm().Connect( api_context.ctx, pi )
	if err != nil {
		*result = ipfs_SubmitError( errors.New(fmt.Sprintf( "failed to connect to %s: %s", pi.ID, err )) )
		return err
	}

	*result = 1
	return nil
}

//export ipfs_CoreAPI_Swarm_Connect_async
func ipfs_CoreAPI_Swarm_Connect_async( api_handle int64, addr *C.char, result *int64 ) int64 {
	api, ec := coreAPI_from_handle( api_handle )
	if ec <= 0 {
		return ec
	}

	go perform_swarm_connect( api, C.GoString( addr ), result )
	return 1
}
