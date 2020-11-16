
@[Link("ipfs")]
lib LibIPFS
	enum ErrorCode : Int64
		Error = 0
		NoError = 1
	end
	alias CString = UInt8*

	fun ipfs_Init() : Void
	fun ipfs_Cleanup() : Void

	fun ipfs_GetError( handle : Int64 ) : CString
	fun ipfs_ReleaseError( handle : Int64 ) : Void

	type PluginLoaderHandle = Int64
	fun ipfs_Loader_PluginLoader_Create( path : CString ) : PluginLoaderHandle
	fun ipfs_Loader_PluginLoader_Initialize( handle : PluginLoaderHandle ) : ErrorCode
	fun ipfs_Loader_PluginLoader_Inject( handle : PluginLoaderHandle ) : ErrorCode

	type ConfigHandle = Int64
	type IoHandle = Int64
	fun ipfs_Config_Init_unsafe = "ipfs_Config_Init"( io : Int64, size : Int32 ) : ConfigHandle

	type RepoHandle = Int64
	fun ipfs_FSRepo_Init( repo_path : CString, cfg_handle : ConfigHandle ) : ErrorCode
	fun ipfs_FSRepo_Open( repo_path : CString ) : RepoHandle

	type BuildCfgHandle = Int64
	fun ipfs_BuildCfg_Create() : BuildCfgHandle
	fun ipfs_BuildCfg_SetOnline( handle : BuildCfgHandle, state : Int32 ) : ErrorCode
	fun ipfs_BuildCfg_SetRouting( handle : BuildCfgHandle, option : Int32 ) : ErrorCode
	fun ipfs_BuildCfg_SetRepo( handle : BuildCfgHandle, repo : RepoHandle ) : ErrorCode
	fun ipfs_BuildCfg_Release( handle : BuildCfgHandle ) : ErrorCode

	type CoreAPIHandle = Int64
	type NodeHandle = Int64
	fun ipfs_CoreAPI_Create( cfg : BuildCfgHandle ) : CoreAPIHandle
	fun ipfs_CoreAPI_Swarm_Connect_async( api : CoreAPIHandle, peerAddr : CString, complete : Int64* ) : ErrorCode
	fun ipfs_CoreAPI_Unixfs_Get( api : CoreAPIHandle, cid : CString ) : NodeHandle

	fun ipfs_Node_GetType( node : NodeHandle ) : Int64
	fun ipfs_Node_Read( node : NodeHandle, unsafe_bytes : UInt8*, limit : Int32, offset : Int64 ) : Int64
	#fun ipfs_Node_Read_async( node : NodeHandle, unsafe_bytes : UInt8*, limit : Int32, offset : Int64, complete : Int64*) ) : Int64

	fun ipfs_RunGoroutines() : ErrorCode
end

