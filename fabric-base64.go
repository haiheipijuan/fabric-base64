package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var (
	base64Digits     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	coder            = base64.NewEncoding(base64Digits)
	ordererMspFloder string
	peerMspFloder    string
	outFile          string
)

var mainCmd = &cobra.Command{
	Use: "fabric-base",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		result, err := generateBase64(ordererMspFloder, peerMspFloder, outFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			cmd.HelpFunc()(cmd, args)
			return
		}
		fmt.Println(string(result))
	},
}

func main() {
	mainFlags := mainCmd.PersistentFlags()
	mainFlags.StringVarP(&ordererMspFloder, "orderermsp", "O", "", "The floder which store orderer msp files")
	mainFlags.StringVarP(&peerMspFloder, "peermsp", "P", "", "The floder which store peer msp files")
	mainFlags.StringVarP(&outFile, "outfile", "o", "", "The file which store result json")

	if mainCmd.Execute() != nil {
		os.Exit(1)
	}
}

func generateBase64(ordMsp, peerMsp, outFile string) ([]byte, error) {
	var result []byte

	if ordMsp == "" && peerMsp == "" {
		return result, fmt.Errorf("Both orderer msp and peer msp is nil\n")
	}

	jsonMap := make(map[string]interface{})
	if ordMsp != "" {
		ordMap, err := readDirPem(ordMsp)
		if err != nil {
			return result, fmt.Errorf("Read orderer cert error:%v\n", err)
		}

		jsonMap["orderer"] = ordMap
	}

	if peerMsp != "" {
		peerMap, err := readDirPem(peerMsp)
		if err != nil {
			return result, fmt.Errorf("Read peer cert error:%v\n", err)
		}

		jsonMap["peer"] = peerMap
	}

	data, err := json.Marshal(jsonMap)
	if err != nil {
		return result, fmt.Errorf("Json marshal error:%v\n", err)
	}

	var out bytes.Buffer
	err = json.Indent(&out, data, "", "\t")
	if err != nil {
		return result, fmt.Errorf("Json indent error:%v\n", err)
	}

	if outFile != "" {
		err := ioutil.WriteFile(outFile, out.Bytes(), 0666)
		if err != nil {
			return result, fmt.Errorf("Write to %v error:%v\n", outFile, err)
		}
		return []byte("Success"), nil
	}

	return out.Bytes(), nil
}

func readDirPem(dir string) (result map[string]interface{}, err error) {
	result = make(map[string]interface{})

	var data []byte
	// admin cert
	admincertDirList, _ := ioutil.ReadDir(dir + "/admincerts")
	data, err = ioutil.ReadFile(dir + "/admincerts/" + admincertDirList[0].Name())
	if err != nil {
		return result, err
	}
	result["admincert"] = string(base64Encode(data))

	// ca cert
	cacertDirList, _ := ioutil.ReadDir(dir + "/cacerts")
	data, err = ioutil.ReadFile(dir + "/cacerts/" + cacertDirList[0].Name())
	if err != nil {
		return result, err
	}
	result["cacert"] = string(base64Encode(data))

	// tls cert
	tlscertDirList, _ := ioutil.ReadDir(dir + "/tlscacerts")
	data, err = ioutil.ReadFile(dir + "/tlscacerts/" + tlscertDirList[0].Name())
	if err != nil {
		return result, err
	}
	result["tlscert"] = string(base64Encode(data))

	return result, nil
}

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func base64Decode(src []byte) ([]byte, error) {
	return coder.DecodeString(string(src))
}
