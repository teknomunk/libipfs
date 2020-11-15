@[Link("ipfs")]
lib LibIPFS
	fun ipfs_Init() : Void
	fun ipfs_Cleanup() : Void
end

module IPFS
	LibIPFS.ipfs_Init()
	at_exit {
		LibIPFS.ipfs_Cleanup()
	}
end

