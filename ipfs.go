package kerala

import (

    "github.com/jbenet/go-ipfs/core"
    "github.com/jbenet/go-ipfs/repo/fsrepo"
    "code.google.com/p/go.net/context"
	"fmt"
	u "github.com/jbenet/go-ipfs/util"
	merkledag "github.com/jbenet/go-ipfs/merkledag"
    "net/http"
	"bytes"
	"strings"
"io/ioutil"
"encoding/hex"
"encoding/json"
	
	

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

func GetStrings(node *core.IpfsNode, userID string) ([]string, error) {
	
	f,err := ioutil.ReadFile("output.html")
	if err != nil {
		return nil, err
	}
	if string(f) == "" {
		return nil, nil
	} else {
		if(userID == "") {
			var Key = u.B58KeyDecode(string(f))
		    var tweetArray = resolveAllInOrder(node,Key)
			return tweetArray, nil
			
			
		} else {
			var Key = u.B58KeyDecode(userID)
		    var tweetArray = resolveAllInOrder(node,Key)
			return tweetArray, nil
			
			
			
		}
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


func Pay(fee string, from_address string, to_address string, amount string, asset_id string, private_key string) (string) {
	unsignedResponse := sendasset(fee,from_address,to_address,amount,asset_id)
    signedResponse := signtransaction(unsignedResponse, private_key)
	transactionHash := pushtransaction(signedResponse)
	return transactionHash
}

func sendasset(fee_placeholder string, from_address string, to_address string, amount_placeholder string, asset_id_placeholder string) (string) {
	   
	//fee_placeholder 1000
	//from_address 1zLkEoZF7Zdoso57h9si5fKxrKopnGSDn
	//to_address akSjSW57xhGp86K6JFXXroACfRCw7SPv637
	//amount 10
	//asset_id AHthB6AQHaSS9VffkfMqTKTxVV43Dgst36
	
	
	client := &http.Client{}
	str := "{\n  \"fees\": fee_placeholder,\n  \"from\": \"from_address\",\n  \"to\": [\n    {\n      \"address\": \"to_address\",\n      \"amount\": \"amount_placeholder\",\n      \"asset_id\": \"asset_id_placeholder\"\n    }\n  ]\n}"
    strings.Replace(str, "fee_placeholder", fee_placeholder, -1)
    strings.Replace(str, "from_address", from_address, -1)
    strings.Replace(str, "to_address", to_address, -1)
    strings.Replace(str, "amount_placeholder", amount_placeholder, -1)
    strings.Replace(str, "asset_id_placeholder", asset_id_placeholder, -1)
	
	
	
	
		
		body := []byte(str)

		req, _ := http.NewRequest("POST", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/sendasset?format=json", bytes.NewBuffer(body))

		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Errored when sending request to the server")
		}

		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)

		fmt.Println(resp.Status)
		fmt.Println(string(resp_body))
		

		strhex := hex.EncodeToString(resp_body)
		
		return string(strhex)
	}


func signtransaction(hex_placeholder string, priv_key_placeholder string) (string) {
		client := &http.Client{}
		
		//hex representation of unsigned transaction placeholder 0100000001d238c42ec059b8c7747cd51debb4310108f6279d14957472822cf061a660828b000000001976a914760fdb3483204406ddb73a45b20b7c9be61d0a7e88acffffffff0288130000000000001976a91430a5d35558ade668b8829a2a0f60a3f10358327e88ac306f0100000000001976a914760fdb3483204406ddb73a45b20b7c9be61d0a7e88ac00000000
		
		//private keys placeholder D8414E7062013DD24D3A3E71EFA8C72142A63F45E3B1AFA4653AFDFD9BC8D67E
		
		str := "{\n  \"transaction\": \"hex_placeholder\",\n  \"keys\": [\n    \"priv_key_placeholder\"\n  ]\n}"
	    strings.Replace(str, "hex_placeholder", hex_placeholder, -1)
	    strings.Replace(str, "priv_key_placeholder", priv_key_placeholder, -1)
		body := []byte(str)

		
		
		
		req, _ := http.NewRequest("POST", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/signtransaction", bytes.NewBuffer(body))

		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("Errored when sending request to the server")
		}

		defer resp.Body.Close()
		resp_body, _ := ioutil.ReadAll(resp.Body)

		//fmt.Println(resp.Status)
		//fmt.Println(string(resp_body))
		return string(resp_body)
	}

func pushtransaction(input_placeholder string) (string) {
	client := &http.Client{}
	//placeholder string 0100000001d238c42ec059b8c7747cd51debb4310108f6279d14957472822cf061a660828b000000006b483045022100d326257244e8cb86889509cf5b4717edf273d9e6e643f571c434753059eb01a902204aa761f44e89b55af0e2fa0caef580401a4ba61eebf8bc29020ce23f6fab1ee2012102661ac805eef8015c7c8d617c65ef327c4f2272fd5d9e97456a0d32d3bcf6f563ffffffff0288130000000000001976a91430a5d35558ade668b8829a2a0f60a3f10358327e88ac306f0100000000001976a914760fdb3483204406ddb73a45b20b7c9be61d0a7e88ac00000000
    str := "\"input_placeholder\""
    strings.Replace(str, "input_placeholder", input_placeholder, -1)
	
	
	body := []byte(str)

	req, _ := http.NewRequest("POST", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/sendrawtransaction", bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	return string(resp_body)
	
}


func GetBalance(my_address string) (float64) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/addresses/address", nil)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)
	
    var dat map[string]interface{}
	
	if err := json.Unmarshal([]byte(string(resp_body)), &dat); err != nil {
	        panic(err)
	    }
		
	    num := dat["balance"]
	fmt.Println("NUM" , num.(float64))

	return num.(float64)
}

