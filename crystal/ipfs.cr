module IPFS
	LibIPFS.ipfs_Init()
	at_exit {
		LibIPFS.ipfs_Cleanup()
	}
	spawn {
		LibIPFS.ipfs_RunGoroutines()
	}

	def self.check_error( e )
		raise_error(e) if e <= 0
	end
	def self.check_error( e : LibIPFS::ErrorCode )
		case e
		when LibIPFS::ErrorCode::NoError
			return
		else
			raise_error(e.to_i64)
		end
	end
	def self.raise_error( error_code )
		str = String.new( LibIPFS.ipfs_GetError( error_code ) )
		LibIPFS.ipfs_ReleaseError( error_code )
		raise str
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

	class CoreAPI
		getter handle
		def initialize( cfg : BuildCfg )
			IPFS.check_error( @handle = LibIPFS.ipfs_CoreAPI_Create( cfg.handle ) )
		end

		struct Swarm
			def initialize( @api : CoreAPI ); end

			def connect( peerAddr )
				completion : Int32 = 0
				ptr = pointerof(completion)
				LibIPFS.ipfs_CoreAPI_Swarm_Connect_async( @api.handle, peerAddr, ptr )
				while ptr[0] == 0
					Fiber.yield()
					LibIPFS.ipfs_RunGoroutines()
				end
				if ptr[0] < 0
					raise LibIPFS.ipfs_AsyncError( ptr[0] )
				end
			end
		end
		def swarm(); Swarm.new(self); end

		struct UnixFS
			def initialize( @api : CoreAPI ); end
		end
		def unixfs(); UnixFS.new(self); end
	end
end
