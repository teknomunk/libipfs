This is a Proof of Concept C API for go-ipfs. Only a handful of API functions have
been written, barely enough to do anything besides demonstrate that such an API
is possible.

```
extern void ipfs_Init();
extern void ipfs_Cleanup();
extern char* ipfs_GetError(GoInt64 handle);
extern GoInt64 ipfs_ReleaseError(GoInt64 handle);
extern GoInt64 ipfs_Loader_PluginLoader_Create(char* plugin_path);
extern GoInt64 ipfs_Loader_PluginLoader_Initialize(GoInt64 handle);
extern GoInt64 ipfs_Loader_PluginLoader_Inject(GoInt64 handle);
extern GoInt64 ipfs_Config_Init(GoInt64 io_handle, GoInt32 size);
extern GoInt64 ipfs_FSRepo_Init(char* repo_path, GoInt64 cfg_handle);
extern GoInt64 ipfs_FSRepo_Open(char* repo_path);
extern GoInt64 ipfs_BuildCfg_Create();
extern GoInt64 ipfs_BuildCfg_SetOnline(GoInt64 handle, GoInt32 state);
extern GoInt64 ipfs_BuildCfg_SetRouting(GoInt64 handle, GoInt32 option);
extern GoInt64 ipfs_BuildCfg_SetRepo(GoInt64 cfg_handle, GoInt64 repo_handle);
extern GoInt64 ipfs_BuildCfg_Release(GoInt64 handle);
extern GoInt64 ipfs_CoreAPI_Create(GoInt64 cfg_handle);
extern GoInt64 ipfs_CoreAPI_Swarm_Connect_async(GoInt64 api_handle, char* addr, GoInt64* result);
extern GoInt64 ipfs_CoreAPI_Unixfs_Get(GoInt64 api_handle, char* cid_str);
extern GoInt64 ipfs_Node_Read(GoInt64 handle, void* unsafe_bytes, GoInt32 bytes_limit, GoInt64 offset);
extern GoInt64 ipfs_RunGoroutines();
```

# Library Methodology

The interface was built to follow the following rules:

* Handles for all types that are not:
  * integer or floating point numbers
  * strings
  * booleans
  * enumerations
  * callbacks of the type  ``int(*)( void* /* callback_data */, /* extra arguments */ )``
* Separate function calls for setting/getting fields in complex types
* function prefix of ipfs_

# Example Code

There is an example program written in crystal included to show that the
API works.  The code is based off the ipfs-as-a-library example in go-ipfs.
To build and run the example code, you will need the crystal compiler.

# Run build.sh to compile libipfs.so
# Copy to /usr/local/lib/
# ldconfig - /usr/local/lib/libipfs.so
# export LD\_LIBRARY\_PATH=/usr/local/lib/
# crystal build crystal/example.cr
# ./example

