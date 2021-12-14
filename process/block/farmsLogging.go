package block

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var farmAddresses = []string{
	"erd1qqqqqqqqqqqqqpgqv4ks4nzn2cw96mm06lt7s2l3xfrsznmp2jpsszdry5",
	"erd1qqqqqqqqqqqqqpgqye633y7k0zd7nedfnp3m48h24qygm5jl2jpslxallh",
	"erd1qqqqqqqqqqqqqpgqsw9pssy8rchjeyfh8jfafvl3ynum0p9k2jps6lwewp",
}

var farmEndpoints = []string{
	"getFarmingTokenReserve",
	"getFarmTokenSupply",
	"getRewardPerShare",
	"getRewardReserve",
}

var poolAddresses = []string{
	"erd1qqqqqqqqqqqqqpgqa0fsfshnff4n76jhcye6k7uvd7qacsq42jpsp6shh2",
	"erd1qqqqqqqqqqqqqpgqeel2kumf0r8ffyhth7pqdujjat9nx0862jpsg2pqaq",
}

var poolEndpoints = []string{
	"getReservesAndTotalSupply",
}

func (sp *shardProcessor) displayDEXInfo(headerHash []byte, timestamp uint64) {
	dexResponse := &DEXResponse{}

	client := http.Client{Timeout: 10 * time.Second}
	sp.displayFarmInfo(headerHash, timestamp, client, dexResponse)
	sp.displayPoolInfo(headerHash, timestamp, client, dexResponse)
}

func (sp *shardProcessor) displayFarmInfo(headerHash []byte, timestamp uint64, client http.Client, dexResponse *DEXResponse) {
	for _, fa := range farmAddresses {
		str := fmt.Sprintf("%s,%v,%s", hex.EncodeToString(headerHash), timestamp, fa)
		for _, method := range farmEndpoints {
			err := sp.getInformation(headerHash, timestamp, client, dexResponse, fa, method)

			if err != nil || dexResponse == nil || len(dexResponse.Error) > 0 || len(dexResponse.Data.Data.ReturnData) != 1 {
				str += "NaN,"
				continue
			}

			decoded, err := base64.StdEncoding.DecodeString(dexResponse.Data.Data.ReturnData[0])
			if err != nil {
				log.Error(err.Error(), "hash", headerHash, "timestamp", timestamp)
			}
			bigInt := big.NewInt(0).SetBytes(decoded)
			dexString := fmt.Sprintf("%s", bigInt.String())

			str = fmt.Sprintf("%s,%s", str, dexString)
		}

		log.Info("dex", "farmInfo", str)
	}
}

func (sp *shardProcessor) displayPoolInfo(headerHash []byte, timestamp uint64, client http.Client, dexResponse *DEXResponse) {
	for _, fa := range poolAddresses {
		str := fmt.Sprintf("%s,%v,%s", hex.EncodeToString(headerHash), timestamp, fa)
		for _, method := range poolEndpoints {
			err := sp.getInformation(headerHash, timestamp, client, dexResponse, fa, method)

			if err != nil || dexResponse == nil || len(dexResponse.Error) > 0 || len(dexResponse.Data.Data.ReturnData) != 3 {
				str += "NaN,"
				continue
			}

			dexString := ""
			for i := 0; i < 3; i++ {
				decoded, err := base64.StdEncoding.DecodeString(dexResponse.Data.Data.ReturnData[i])
				if err != nil {
					log.Error(err.Error(), "hash", headerHash, "timestamp", timestamp)
				}
				bigInt := big.NewInt(0).SetBytes(decoded)
				dexString += fmt.Sprintf("%s,", bigInt.String())
			}
			str = fmt.Sprintf("%s,%s", str, dexString)

		}

		log.Info("dex", "poolInfo", str)
	}
}

func (sp *shardProcessor) getInformation(headerHash []byte, timestamp uint64, client http.Client, dexResponse *DEXResponse, fa string, method string) error {
	//Encode the data
	postBody, _ := json.Marshal(map[string]string{
		"scAddress": fa,
		"funcName":  method,
	})

	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	resp, err := client.Post("http://127.0.0.1:8081/vm-values/query", "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Info("we have error", "err", err.Error())
		return err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Info("we have body error", "err", err.Error())
		return err
	}

	err = json.Unmarshal(body, dexResponse)
	if err != nil {
		log.Info("we have unmarshal error", "err", err.Error())
		return err
	}

	return nil
}

type DEXResponse struct {
	Data struct {
		Data struct {
			ReturnData     []string `json:"returnData"`
			ReturnCode     string   `json:"returnCode"`
			ReturnMessage  string   `json:"returnMessage"`
			GasRemaining   float64  `json:"gasRemaining"`
			GasRefund      int      `json:"gasRefund"`
			OutputAccounts struct {
				D0Acc53561C5D6F6Fd7D7E82Bf13247014F615483 struct {
					Address        string      `json:"address"`
					Nonce          int         `json:"nonce"`
					Balance        interface{} `json:"balance"`
					BalanceDelta   int         `json:"balanceDelta"`
					StorageUpdates struct {
					} `json:"storageUpdates"`
					Code            interface{}   `json:"code"`
					CodeMetaData    interface{}   `json:"codeMetaData"`
					OutputTransfers []interface{} `json:"outputTransfers"`
					CallType        int           `json:"callType"`
				} `json:"00000000000000000500656d0acc53561c5d6f6fd7d7e82bf13247014f615483"`
			} `json:"outputAccounts"`
			DeletedAccounts []interface{} `json:"deletedAccounts"`
			TouchedAccounts []interface{} `json:"touchedAccounts"`
			Logs            []interface{} `json:"logs"`
		} `json:"data"`
	} `json:"data"`
	Error string `json:"error"`
	Code  string `json:"code"`
}
