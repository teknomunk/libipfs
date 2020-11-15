require "./libipfs"
require "./libp2p"
require "./ipfs"

# Below is a port of the ipfs-as-a-library example program in crystal

def setupPlugins( plugin_path )
	loader = IPFS::PluginLoader.new( File.join( plugin_path, "plugins" ) )
	loader.initialize_plugins()
	loader.inject()
end
def createTempRepo()
	repoPath = File.tempname("ipfs-shell")
	puts "repo at #{repoPath}"
	cfg = IPFS::Config.new( keysize: 2048 )
	IPFS::FSRepo.init( repoPath, cfg )

	return repoPath
end
def createNode( repoPath )
	repo = IPFS::FSRepo.new( repoPath )

	nodeOptions = IPFS::BuildCfg.new(
		online: true,
		routing: LibP2P::DHTOption,
		repo: repo
	)
end

def spawnEphemeral()
	setupPlugins("")

	repoPath = createTempRepo()

	return createNode(repoPath)
end

spawnEphemeral()
