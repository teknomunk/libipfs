module IPFS
	LibIPFS.ipfs_Init()
	at_exit {
		LibIPFS.ipfs_Cleanup()
	}

	def self.check_error( e )
		raise_error if e <= 0
	end
	def self.check_error( e : LibIPFS::ErrorCode )
		case e
		when LibIPFS::ErrorCode::NoError
			return
		else
			raise_error
		end
	end
	def self.raise_error()
		raise String.new( LibIPFS.ipfs_LastError() )
	end

	class PluginLoader
		getter handle : LibIPFS::PluginLoaderHandle

		def initialize( path : String )
			IPFS.check_error( @handle = LibIPFS.ipfs_Loader_PluginLoader_Create( path ) )
		end
		def initialize_plugins()
			IPFS.check_error LibIPFS.ipfs_Loader_PluginLoader_Initialize( @handle )
		end
		def inject()
			IPFS.check_error LibIPFS.ipfs_Loader_PluginLoader_Inject( @handle )
		end
	end

	class Config
		getter handle : LibIPFS::ConfigHandle

		def initialize( keysize : Int )
			IPFS.check_error( @handle = LibIPFS.ipfs_Config_Init_unsafe( 0, keysize ) )
		end
	end

	class FSRepo
		getter handle : LibIPFS::RepoHandle

		def initialize( path : String )
			IPFS.check_error( @handle = LibIPFS.ipfs_FSRepo_Open( path ) )
		end
		def self.init( path : String, cfg : Config )
			IPFS.check_error LibIPFS.ipfs_FSRepo_Init( path, cfg.handle )
		end
	end

	class BuildCfg
		getter handle : LibIPFS::BuildCfgHandle
		def initialize( *, online = true, routing = LibP2P::DHTClientOption, repo = nil )
			IPFS.check_error( @handle = LibIPFS.ipfs_BuildCfg_Create() )

			IPFS.check_error LibIPFS.ipfs_BuildCfg_SetOnline( @handle, online )
		end
		def finalize()
			IPFS.check_error LibIPFS.ipfs_BuildCfg_Release(@handle)
		end
	end
end
