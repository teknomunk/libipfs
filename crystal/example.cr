require "./libipfs"
require "./libp2p"
require "./ipfs"

# Below is a partial port of the ipfs-as-a-library example program in crystal

## Spawn Ephemeral node
# setupPlugins
loader = IPFS::PluginLoader.new( File.join( "", "plugins" ) )
loader.initialize_plugins()
loader.inject()

# Create temporary repo
repoPath = File.tempname("ipfs-shell")
puts "repo at #{repoPath}"

cfg = IPFS::Config.new( keysize: 2048 )
IPFS::FSRepo.init( repoPath, cfg )

# Create Node
repo = IPFS::FSRepo.new( repoPath )

nodeOptions = IPFS::BuildCfg.new(
	online: true,
	routing: LibP2P::DHTOption,
	repo: repo
)

node = IPFS::CoreAPI.new( nodeOptions )

