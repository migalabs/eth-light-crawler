# Ethereum Light Crawler
This tool uses Ethereum's peer-discovery protocol to measure the size of the Ethereum network (testnets included) 

### How it works?
The `light-crawler` creates an [`Ethereum node`](https://github.com/ethereum/go-ethereum/blob/f53ff0ff4a68ffc56004ab1d5cc244bcb64d3277/p2p/enode/localnode.go#L47) with a new generated ecdsa Private Key to launch a [`Discovery5`](https://github.com/ethereum/go-ethereum/blob/f53ff0ff4a68ffc56004ab1d5cc244bcb64d3277/p2p/discover/v5_udp.go#L61) service. 

Once the node and the service are running, the light-crawler will request random ENRs from the network, indexing in an SQL database all the different atributes from the discovered Nodes (we index all the available ENRs in the network, or in other words, all the Ethereum networks). 

### Requirements
- (recomended) Go1.17+
- (must) make 
- (must) local or remote postgreSQL database

### How to run it?
Once the repo is cloned, you can compile it and generate a binary running:
```
$ make build
```

Once the binary was generated, you can execute it by running:
```
USAGE:
   ./build/eth-light-crawler [subcommands] [arguments]

COMMANDS:
   discv5   crawl Ethereum's public DHT thought the Discovery 5.1 protocol
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

The current crawler only support the discovery5.1 protocol, so to crawl it your can run the following command:
```
USAGE:
   ./build/eth-light-crawler discv5 [arguments...]

OPTIONS:
   --log-level value    verbosity of the logs that will be displayed [debug,warn,info,error] (default: "info") [$IPFS_CID_HOARDER_LOGLEVEL]
   --db-endpoint value  login endpoint to the database (default: "postgres://test:password@localhost:5432/eth_light_crawler") [$IPFS_CID_HOARDER_DB_ENDPOINT]
   --port value         port number that we want to use/advertise in the Ethereum network (default: 9001)
   --reset-db           reset the content of the db tables (default: false)
   --help, -h           show help (default: false)
```
_NOTE: the `light-crawler` will require to have a postgreSQL database created before running it, it will only create the required tables to run._

### Maintainers
Miga Labs / @cortze

### Contributions 
This tool is open for any kind of feedback and contribution, so please feel free to approach us!

