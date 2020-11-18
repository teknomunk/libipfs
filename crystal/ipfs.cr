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
	class Node < IO
		getter handle : LibIPFS::NodeHandle
		@offset : Int64 = 0
		def initialize( @handle : LibIPFS::NodeHandle ); end

		def self.from_path( path : String )
			handle = LibIPFS.ipfs_Node_NewFromPath( path )
			IPFS.check_error(handle)

			return Node.new(handle)
		end
		def self.from_io( io : IO )
			raise "Not yet implemented"
		end

		def read( slice : Bytes )
			res = LibIPFS.ipfs_Node_Read( @handle, slice, slice.size, @offset )
			IPFS.check_error( res )
			@offset += ( res - 1 )
			return res - 1
		end
		def write( slice : Bytes ) : Nil
			raise "IPFS::Node is read only"
		end

		def to_s( io )
			io << "[Node handle=#{ handle.to_i64 }]"
		end
		def inspect( io )
			io << "[Node handle=#{ handle.to_i64 }]"
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
				completion : Int64 = 0
				ptr = pointerof(completion)
				LibIPFS.ipfs_CoreAPI_Swarm_Connect_async( @api.handle, peerAddr, ptr )
				while ptr[0] == 0
					Fiber.yield()
					LibIPFS.ipfs_RunGoroutines()
				end
				IPFS.check_error( ptr[0] )
				puts "connected to #{peerAddr}"
			end
		end
		def swarm(); Swarm.new(self); end

		struct UnixFS
			def initialize( @api : CoreAPI ); end

			def add( node : Node )
				res = LibIPFS.ipfs_CoreAPI_Unixfs_Add( @api.handle, node.handle )
				IPFS.check_error(res)

				cid = LibIPFS.ipfs_GetString(res)
				LibIPFS.ipfs_ReleaseString(res)

				return String.new(cid)
			end
			def get( cid : String )
				IPFS.check_error( handle = LibIPFS.ipfs_CoreAPI_Unixfs_Get( @api.handle, cid ) )
				return Node.new(handle)
			end
		end
		def unixfs(); UnixFS.new(self); end
	end
end

