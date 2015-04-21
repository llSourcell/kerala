# Kerala



<img align="right" height="300" src="http://i180.photobucket.com/albums/x247/Mandee_Candee/Palm-Tree-cartoon.gif">

> An IPFS wrapper for string storage and retrieval 

Kerala offers a an easy-to-use wrapper to store and retrieve links on the Interplanetary File System. (IPFS) I created this library to aid me in creating a decentralized twitter. It can be used for any sort of decentralized application that stores user tweets/posts/microblogs.

Kerala


Add data to IPFS.<br />
Retrieve data from IPFS.<br />
Resolve a peerID to get their data address<br />
Generate your own dapp asset address<br />
Pay anyone with dapp assets<br />


TODO

-Implement IPFS Keystore for encrypting data and sharing it with trusted nodes.
-Namecoin registration for PeerIDs

All pull requests, issue creation, and advice are welcome. 

## Install

```sh
$ go get -u github.com/siraj/go-kerala/kerala
```

Kerala depends on [IPFS](https://github.com/jbenet/go-ipfs). 

## Usage

```go
//Start a node
node, err := kerala.StartNode()
	if err != nil {
		panic(err)
	}

//Add your text to IPFS (Creates MerkleDAG)
var userInput = r.Form["sometext"]
Key, err := kerala.AddString(node, userInput[0])

//Resolve PeerID to get MerkleDAG
pointsTo, err := kerala.GetDAG(node, node.Identity.Pretty())

//Get all your text from IPFS (Retrieves MerkleDAG)  
tweetArray, err := kerala.GetStrings(node, pointsTo.B58String())

//Pay another node (Arguments are - fee, your address, their address, amount, asset address, private keys)
hash := kerala.Pay("1000","1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh", "akSjSW57xhGp86K6JFXXroACfRCw7SPv637", "10", "AHthB6AQHaSS9VffkfMqTKTxVV43Dgst36", "L1jftH241t2rhQSTrru9Vd2QumX4VuGsPhVfSPvibc4TYU4aGdaa" )

//Generate an asset address
address := kerala.GenerateAddress()

//Get your current balance
balance := kerala.GetBalance("1HihKUXo6UEjJzm4DZ9oQFPu2uVc9YK9Wh")

