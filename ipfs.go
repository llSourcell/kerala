package kerala

import (

    "github.com/jbenet/go-ipfs/core"
    "github.com/jbenet/go-ipfs/repo/fsrepo"
    "code.google.com/p/go.net/context"
	"fmt"
	"io/ioutil"
	u "github.com/jbenet/go-ipfs/util"
	merkledag "github.com/jbenet/go-ipfs/merkledag"
	
	

)


func StartNode() (*core.IpfsNode, error) {
	//[1] init IPFS node
    builder := core.NewNodeBuilder().Online()
    r := fsrepo.At("~/.go-ipfs")
    if err := r.Open(); err != nil {
       return nil, err
    }
    builder.SetRepo(r)
    // Make our 'master' context and defer cancelling it
    ctx, _ := context.WithCancel(context.Background())
    //defer cancel()

    node, err := builder.Build(ctx)
    if err != nil {
       return nil, err
    }
    fmt.Printf("I am peer: %s\n", node.Identity)
	return node, nil
}

func GetStrings(node *core.IpfsNode) ([]string, error) {
	
	f,err := ioutil.ReadFile("output.html")
	if err != nil {
		return nil, err
	}
	if string(f) == "" {
		return nil, nil
	} else {
		var Key = u.B58KeyDecode(string(f))
	    var tweetArray = resolveAllInOrder(node,Key)
		return tweetArray, nil
	}
	
}

func resolveAllInOrder(nd * core.IpfsNode, k u.Key) []string {
	var stringArr []string
	var node * merkledag.Node
	node, err := nd.DAG.Get(k)
	fmt.Printf("the node is", node)
	if err != nil {
		fmt.Println(err)
		//return
	}
	fmt.Printf("bout to crash")
	fmt.Printf("%s ", string(node.Data[:]))
	fmt.Println("not crashed ")
	
	for ;; {
		var err error

		if len(node.Links) == 0 {
			fmt.Println("Links are 0")
			break;
		}

		node, err = node.Links[0].GetNode(nd.DAG)
		fmt.Printf("i pulled a node link %s", node)
		if err != nil {
			fmt.Println(err)
			//return
		}

		fmt.Printf("%s ", string(node.Data[:]))
		stringArr = append(stringArr, string(node.Data[:]))
	}

	fmt.Printf("\n");
	
	return stringArr

}

func AddString(node *core.IpfsNode, inputString string) (u.Key, error) {
	
	//[1] Check if key is stored locally
	f,err := ioutil.ReadFile("output.html")
 	if err!=nil {
  			return "", err
  		}
	//[2] If key is not stored locally, the user is new
		if string(f) == "" {
			//[3] Initialize a MerkleDAG node and key
			var NewNode * merkledag.Node
			var Key u.Key
			//[4] Fill the node with user input
			NewNode = makeStringNode(inputString)
			//[5] Add the node to IPFS
			Key, _ = node.DAG.Add(NewNode)
			//[6] Add the user input key to local store
			err := ioutil.WriteFile("output.html", []byte(Key.B58String()) , 0777)
			if err != nil {
				return "", err 
			}
			return Key, nil
			} else {
					//[7] Initialize a new MerkleDAG node and key
					var NewNode * merkledag.Node
					var Key u.Key
					//[8] Fill the node with user input
					NewNode = makeStringNode(inputString)
					//[9] Read the hash from the file
			 		f,err := ioutil.ReadFile("output.html")
					//[10] Convert it into a key
					Key = u.B58KeyDecode(string(f))
					//[11] Get the Old MerkleDAG node and key
					var OldNode * merkledag.Node
					var Key2 u.Key
					OldNode, _ = node.DAG.Get(Key)
					//[12]Add a link to the old node
			 		NewNode.AddNodeLink("next", OldNode)
					//[13] Add thew new node to IPFS
					Key2, _ = node.DAG.Add(NewNode)
					//[14] Add the new node key to local store (overwrite)
					err2 := ioutil.WriteFile("output.html", []byte(Key2.B58String()) , 0777)
					if err2 != nil {
						return "", err
					}
					return Key2, nil
				}
}


func makeStringNode(s string) * merkledag.Node {
	n := new(merkledag.Node)
	n.Data = make([]byte, len(s))
	copy(n.Data, s)
	return n;
}