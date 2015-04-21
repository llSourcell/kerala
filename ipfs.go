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
commands "github.com/jbenet/go-ipfs/core/commands"



	
	

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
	return node, nil
}



func GetStrings(node *core.IpfsNode, userID string) ([]string, error) {

			var Key = u.B58KeyDecode(userID)
		    var tweetArray = resolveAllInOrder(node,Key)
			return tweetArray, nil
	
}


func GetDAG(node *core.IpfsNode, id string) (u.Key, error) {
	
	pointsTo, err := node.Namesys.Resolve(node.Context(), id)
	return pointsTo, err
	
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
			break;
		}

		node, err = node.Links[0].GetNode(nd.DAG)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s ", string(node.Data[:]))
		stringArr = append(stringArr, string(node.Data[:]))
	}

	fmt.Printf("\n");
	
	return stringArr

}

func AddString(node *core.IpfsNode, inputString string) (u.Key, error) {
	

		pointsTo, err := node.Namesys.Resolve(node.Context(), node.Identity.Pretty())
		
		//If there is an error, user is new and hasn't yet created a DAG.
		if err != nil {
			//[3] Initialize a MerkleDAG node and key
			var NewNode * merkledag.Node
			var Key u.Key
			//[4] Fill the node with user input
			NewNode = makeStringNode(inputString)
			//[5] Add the node to IPFS
			Key, _ = node.DAG.Add(NewNode)
			// //publish to IPNS
			output, err := commands.Publish(node, node.PrivateKey,Key.B58String())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("You published to IPNS. Your peer ID is ", output.Name)
			}
			
			return Key, nil
			
		} else {
			//[7] Initialize a new MerkleDAG node and key
			var NewNode * merkledag.Node
			var Key u.Key
			//[8] Fill the node with user input
			NewNode = makeStringNode(inputString)
			//[10] Convert it into a key
			Key = u.B58KeyDecode(pointsTo.B58String())
			//[11] Get the Old MerkleDAG node and key
			var OldNode * merkledag.Node
			var Key2 u.Key
			OldNode, _ = node.DAG.Get(Key)
			//[12]Add a link to the old node
	 		NewNode.AddNodeLink("next", OldNode)
			//[13] Add thew new node to IPFS
			Key2, _ = node.DAG.Add(NewNode)
			// //publish to IPNS
			output, err := commands.Publish(node, node.PrivateKey,Key2.B58String())
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("You published to IPNS. Your peer ID is ", output.Name)
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
			fmt.Println("Error when sending request to the server")
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

		return string(resp_body)
	}

func pushtransaction(input_placeholder string) (string) {
	client := &http.Client{}

    str := "\"input_placeholder\""
    strings.Replace(str, "input_placeholder", input_placeholder, -1)
	
	
	body := []byte(str)

	req, _ := http.NewRequest("POST", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/sendrawtransaction", bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error when sending request to the server")
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
		fmt.Println("Error when sending request to the server")
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

func GenerateAddress() (string) {
	client := &http.Client{}

	body := []byte("{\n  \"alias\": \"address label\"\n}")

	req, _ := http.NewRequest("POST", "https://private-anon-e4123b065-coinprism.apiary-mock.com/v1/account/createaddress", bytes.NewBuffer(body))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Coinprism-Username", "username")
	req.Header.Add("X-Coinprism-Password", "account-password")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	resp_body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	return string(resp_body)
}
}

