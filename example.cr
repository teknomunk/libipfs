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

	fun ipfs_Config_Init( io : Int64, size : Int32 ) : Int64
end

module IPFS
	LibIPFS.ipfs_Init()
	at_exit {
		LibIPFS.ipfs_Cleanup()
	}

	def self.raise_ipfs_error()
		raise String.new( LibIPFS.ipfs_LastError() )
	end

	class PluginLoader
		@handle : Int64

		def initialize( path : String )
			@handle = LibIPFS.ipfs_Loader_PluginLoader_Create( path )
			IPFS.raise_ipfs_error if @handle <= 0
		end
		def initialize_plugins()
			error = LibIPFS.ipfs_Loader_PluginLoader_Initialize( @handle )
			IPFS.raise_ipfs_error if error == 0
		end
		def inject()
			error = LibIPFS.ipfs_Loader_PluginLoader_Inject( @handle )
			IPFS.raise_ipfs_error if error == 0
		end
	end

	class Config
		def initialize( keysize : Int )
			@handle = LibIPFS.ipfs_Config_Init( 0, keysize )
			IPFS.raise_ipfs_error if @handle <= 0
		end
	end

	class FSRepo
		def initialize( path : String, cfg : Config )
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
	repoPath = File.tempname("ipfs-shell")

	cfg = IPFS::Config.new( keysize: 2048 )
end
def createNode( repoPath )

end

def spawnEphemeral()
	setupPlugins("")

	repoPath = createTempRepo()

	return createNode(repoPath)
end

spawnEphemeral()
