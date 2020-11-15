@[Link("ipfs")]
lib LibIPFS
	#type Handle = Int64
	#type ErrorCode = Int64

	fun ipfs_Init() : Void
	fun ipfs_Cleanup() : Void
	fun ipfs_LastError() : UInt8*

	fun ipfs_Loader_PluginLoader_Create( path : UInt8* ) : Int64
	fun ipfs_Loader_PluginLoader_Initialize( handle : Int64 ) : Int64
	fun ipfs_Loader_PluginLoader_Inject( handle : Int64 ) : Int64
end

module IPFS
	LibIPFS.ipfs_Init()
	at_exit {
		LibIPFS.ipfs_Cleanup()
	}

	class PluginLoader
		@handle : Int64

		def initialize( path : String )
			@handle = LibIPFS.ipfs_Loader_PluginLoader_Create( path )
			puts "@handle=#{@handle}"
		end
		def initialize_plugins()
			err = LibIPFS.ipfs_Loader_PluginLoader_Initialize( @handle )
		end
		def inject()
			err = LibIPFS.ipfs_Loader_PluginLoader_Inject( @handle )
		end
	end
end

# Below is a port of the ipfs-as-a-library example program in crystal

def setupPlugins( plugin_path )
	loader = IPFS::PluginLoader.new( File.join( plugin_path, "plugins" ) )
	loader.initialize_plugins()
	loader.inject()
end
def createTempRepo()

end
def createNode( repoPath )

end

def spawnEphemeral()
	setupPlugins("")

	repoPath = createTempRepo()

	return createNode(repoPath)
end

spawnEphemeral()
