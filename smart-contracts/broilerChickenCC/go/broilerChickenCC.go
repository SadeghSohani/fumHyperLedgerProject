package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type GrowthInformation struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Instruction string `json:"instruction"`
}

type Chicken struct {
	Type      string            `json:"type"`
	Birthday  string            `json:"birthday"`
	Breed     string            `json:"breed"`
	Price     float64           `json:"price"`
	Owner     string            `json:"owner"`
	GrowthInf GrowthInformation `json:"growthInformation"`
	ForSale   bool              `json:"forSale"`
	TxType      string            `json:"txType"`
	Status      string            `json:"status"`
}

type QueryResult struct {
	Key    string `json:"id"`
	Record *Chicken `json:"asset"`
}

type BatchResult struct {
	Id    string `json:"id"`
	Asset *Batch `json:"asset"`
}

type HistoryModel struct {
	TxId      string `json:"TxId"`
	Value     *Chicken
	Timestamp string `json:"Timestamp"`
	IsDelete  bool   `json:"IsDelete"`
}

type Token struct {
	Type        string  `json:"type"`
	User        string  `json:"user"`
	Amount      float64 `json:"amount"`
	BlockAmount float64 `json:"blockAmount"`
}

type Batch struct {
	Name      string            `json:"name"`
	Tag       string            `json:"tag"`
	Type      string            `json:"type"`
	Count     int64             `json:"count"`
	Owner     string            `json:"owner"`
	ForSale   bool              `json:"forSale"`
	TxType      string            `json:"txType"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

	// gi := GrowthInformation{
	// 	Key:         "Vaccine",
	// 	Value:       "VC123",
	// 	Instruction: "This is the first vaccination.",
	// }
	// bids := make(map[string]BidModel)

	// chickens := []Chicken{
	// 	Chicken{Type: "Chicken", Birthday: "1401/5/18", Breed: "Ross", Price: 12.0, Owner: "Sadegh", GrowthInf: gi, PublicToSell: false, Bids: bids},
	// 	Chicken{Type: "Chicken", Birthday: "1401/5/18", Breed: "Cobb", Price: 10.0, Owner: "Sadegh", GrowthInf: gi, PublicToSell: false, Bids: bids},
	// 	Chicken{Type: "Chicken", Birthday: "1401/5/18", Breed: "Cobb", Price: 10.0, Owner: "Sadegh", GrowthInf: gi, PublicToSell: false, Bids: bids},
	// 	Chicken{Type: "Chicken", Birthday: "1401/5/18", Breed: "Ross", Price: 12.0, Owner: "Sadegh", GrowthInf: gi, PublicToSell: false, Bids: bids},
	// }

	// for i, chicken := range chickens {
	// 	chickenAsBytes, _ := json.Marshal(chicken)
	// 	err := ctx.GetStub().PutState("CHICKEN"+strconv.Itoa(i), chickenAsBytes)

	// 	if err != nil {
	// 		return fmt.Errorf("Failed to put to world state. %s", err.Error())
	// 	}
	// }

	// tokenAsBytes, _ := json.Marshal(Token{Type: "Token", User: "Sadegh", Amount: 0.0, BlockAmount: 0.0})
	// err := ctx.GetStub().PutState("Sadegh", tokenAsBytes)
	// if err != nil {
	// 	return fmt.Errorf("Failed to put to world state. %s", err.Error())
	// }

	return nil
}

func (s *SmartContract) CreateChicken(ctx contractapi.TransactionContextInterface, id string, birthday string, breed string, price float64, owner string) (*Chicken, error) {

	gi := GrowthInformation{
		Key:         "",
		Value:       "",
		Instruction: "",
	}
	chicken := Chicken{Type: "Chicken", Birthday: birthday, Breed: breed, Price: price, Owner: owner, GrowthInf: gi, ForSale: false, TxType: "CreateChicken", Status:"Alive"}

	chickenAsBytes, _ := json.Marshal(chicken)
	err := ctx.GetStub().PutState(id, chickenAsBytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)

	return _chicken, nil
}

func (s *SmartContract) CreateBulkChicken(ctx contractapi.TransactionContextInterface, assetsIds string, birthday string, breed string, price float64, owner string) ([]Chicken, error) {

	ids := strings.Split(assetsIds, "#")

	gi := GrowthInformation{
		Key:         "",
		Value:       "",
		Instruction: "",
	}

	chickens := []Chicken{}

	for i := range ids {
		chicken := Chicken{Type: "Chicken", Birthday: birthday, Breed: breed, Price: price, Owner: owner, GrowthInf: gi, ForSale: false, TxType: "CreateChicken", Status:"Alive"}
		chickenAsBytes, _ := json.Marshal(chicken)
		err := ctx.GetStub().PutState(ids[i], chickenAsBytes)
	
		if err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		chickens = append(chickens, chicken)

	}

	// _chicken := new(Chicken)
	// _ = json.Unmarshal(chickenAsBytes, _chicken)

	return chickens, nil
}

func (s *SmartContract) CreateBulkChickenInBatch(ctx contractapi.TransactionContextInterface,
	 assetsIds string, birthday string, breed string, price float64, owner string,
	 batchId string, batchName string, batchTag string, batchType string) (*Batch, error) {

	ids := strings.Split(assetsIds, "#")

	batch := Batch{Name: batchName, Tag: batchTag, Type: batchType, Count: int64(len(ids)), Owner: owner, ForSale:false, TxType: "CreateBatch"}

	batchAsBytes, _ := json.Marshal(batch)
	err := ctx.GetStub().PutState(batchId, batchAsBytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	gi := GrowthInformation{
		Key:         "",
		Value:       "",
		Instruction: "",
	}

	for i := range ids {
		chicken := Chicken{Type: "Chicken", Birthday: birthday, Breed: breed, Price: price, Owner: batchId, GrowthInf: gi, ForSale: false, TxType: "CreateChicken", Status:"Alive"}
		chickenAsBytes, _ := json.Marshal(chicken)
		err := ctx.GetStub().PutState(ids[i], chickenAsBytes)
	
		if err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

	}

	_batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, _batch)

	return _batch, nil
}

func (s *SmartContract) CreateBatch(ctx contractapi.TransactionContextInterface,
	 id string, name string, tag string, batchType string, owner string) (*Batch, error) {

	batch := Batch{Name: name, Tag: tag, Type: batchType, Count: 0, Owner: owner, ForSale:false, TxType: "CreateBatch"}

	batchAsBytes, _ := json.Marshal(batch)
	err := ctx.GetStub().PutState(id, batchAsBytes)

	if err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, _batch)

	return _batch, nil
}

func (s *SmartContract) QueryChicken(ctx contractapi.TransactionContextInterface, id string) (*Chicken, error) {
	chickenAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if chickenAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, chicken)

	return chicken, nil
}

func (s *SmartContract) QueryBatch(ctx contractapi.TransactionContextInterface, id string) (*Batch, error) {
	batchAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if batchAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, batch)

	return batch, nil
}

func (s *SmartContract) QueryChickenByOwner(ctx contractapi.TransactionContextInterface, id string, owner string) (*Chicken, error) {
	chickenAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if chickenAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, chicken)

	if chicken.Owner == owner {
		return chicken, nil
	}

	return nil, fmt.Errorf("Permission denied.")

}

func (s *SmartContract) QueryBatchByOwner(ctx contractapi.TransactionContextInterface, id string, owner string) (*Batch, error) {
	batchAsBytes, err := ctx.GetStub().GetState(id)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if batchAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", id)
	}

	batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, batch)

	if batch.Owner == owner {
		return batch, nil
	}

	return nil, fmt.Errorf("Permission denied.")

}

func (s *SmartContract) QueryToken(ctx contractapi.TransactionContextInterface, user string) (*Token, error) {
	tokenAsBytes, err := ctx.GetStub().GetState(user)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if tokenAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", user)
	}

	token := new(Token)
	_ = json.Unmarshal(tokenAsBytes, token)

	return token, nil
}

func (s *SmartContract) BuyToken(ctx contractapi.TransactionContextInterface, user string, amount float64) (*Token, error) {
	tokenAsBytes, err := ctx.GetStub().GetState(user)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if tokenAsBytes == nil {

		token := Token{Type: "Token", User: user, Amount: amount, BlockAmount: 0.0}
		tokenAsBytes, _ := json.Marshal(token)
		err := ctx.GetStub().PutState(user, tokenAsBytes)
		if err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		return &token, nil

	} else {
		token := new(Token)
		_ = json.Unmarshal(tokenAsBytes, token)

		token.Amount = token.Amount + amount

		_tokenAsBytes, _ := json.Marshal(token)

		_err := ctx.GetStub().PutState(user, _tokenAsBytes)

		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		return token, nil
	}

}

func (s *SmartContract) TransferToken(ctx contractapi.TransactionContextInterface, sender string, receiver string, amount float64) (*Token, error) {
	sender_token, err := s.QueryToken(ctx, sender)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	receiver_token, err := s.QueryToken(ctx, receiver)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if sender_token.Amount >= amount {
		sender_token.Amount = sender_token.Amount - amount
		receiver_token.Amount = receiver_token.Amount + amount
		rtAsBytes, _ := json.Marshal(receiver_token)
		_err := ctx.GetStub().PutState(receiver, rtAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		stAsBytes, _ := json.Marshal(sender_token)
		_err_ := ctx.GetStub().PutState(sender, stAsBytes)
		if _err_ != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return sender_token, nil
}

func (s *SmartContract) QueryAllChickens(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		chicken := new(Chicken)
		_ = json.Unmarshal(queryResponse.Value, chicken)

		if chicken.Type == "Chicken" {
			queryResult := QueryResult{Key: queryResponse.Key, Record: chicken}
			results = append(results, queryResult)
		}

	}

	return results, nil
}

func (s *SmartContract) QueryChickensByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]QueryResult, error) {

	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		chicken := new(Chicken)
		_ = json.Unmarshal(queryResponse.Value, chicken)

		if chicken.Owner == owner && chicken.Breed != "" {
			queryResult := QueryResult{Key: queryResponse.Key, Record: chicken}
			results = append(results, queryResult)
		}

	}

	return results, nil
}

func (s *SmartContract) GetAssetsOfBatch(ctx contractapi.TransactionContextInterface, batchId string, owner string) ([]QueryResult, error) {

	batchAsBytes, err := ctx.GetStub().GetState(batchId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if batchAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", batchId)
	}

	batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, batch)

	if batch.Owner == owner {
		startKey := ""
		endKey := ""
	
		resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)
	
		if err != nil {
			return nil, err
		}
		defer resultsIterator.Close()
	
		results := []QueryResult{}
	
		for resultsIterator.HasNext() {
			queryResponse, err := resultsIterator.Next()
	
			if err != nil {
				return nil, err
			}
	
			chicken := new(Chicken)
			_ = json.Unmarshal(queryResponse.Value, chicken)
	
			if chicken.Owner == batchId && chicken.Breed != "" {
				queryResult := QueryResult{Key: queryResponse.Key, Record: chicken}
				results = append(results, queryResult)
			}
	
		}
		return results, nil		
	} else {
		return nil, fmt.Errorf("Permission denied.")
	}

}

func (s *SmartContract) QueryAllBatchsByOwner(ctx contractapi.TransactionContextInterface, owner string) ([]BatchResult, error) {

	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []BatchResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		batch := new(Batch)
		_ = json.Unmarshal(queryResponse.Value, batch)

		if batch.Owner == owner && batch.Name != "" {
			queryResult := BatchResult{Id: queryResponse.Key, Asset: batch}
			results = append(results, queryResult)
		}

	}

	return results, nil
}

func (s *SmartContract) QueryPublicChickens(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {

	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		chicken := new(Chicken)
		_ = json.Unmarshal(queryResponse.Value, chicken)

		if chicken.ForSale == true {
			queryResult := QueryResult{Key: queryResponse.Key, Record: chicken}
			results = append(results, queryResult)
		}

	}

	return results, nil
}

func (s *SmartContract) PutGrowthInformation(ctx contractapi.TransactionContextInterface, id string, key string, value string, instruction string, owner string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner {
		chicken.GrowthInf = GrowthInformation{Key: key, Value: value, Instruction: instruction}
		chicken.TxType = "PutGrowthInf"
	}

	chickenAsBytes, _ := json.Marshal(chicken)

	_err := ctx.GetStub().PutState(id, chickenAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)
	return _chicken, nil

}

func (s *SmartContract) ChangeChickenOwner(ctx contractapi.TransactionContextInterface, id string, owner string, newOwner string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner {
		chicken.Owner = newOwner
		chicken.TxType = "ChangeOwner"
	}

	chickenAsBytes, _ := json.Marshal(chicken)

	_err := ctx.GetStub().PutState(id, chickenAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)
	return _chicken, nil

}

func (s *SmartContract) ChangeChickenOwnerPhone(ctx contractapi.TransactionContextInterface, id string, owner string, newOwner string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Status == "FinalProduct" {
		chicken.Owner = newOwner
		chicken.TxType = "ChangeOwner"
	}

	chickenAsBytes, _ := json.Marshal(chicken)

	_err := ctx.GetStub().PutState(id, chickenAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)
	return _chicken, nil

}

func (s *SmartContract) ChangeChickenStatus(ctx contractapi.TransactionContextInterface, id string, owner string, status string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner{
		chicken.Status = status
		chicken.TxType = "ChangeStatus"
	}

	chickenAsBytes, _ := json.Marshal(chicken)

	_err := ctx.GetStub().PutState(id, chickenAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)
	return _chicken, nil

}

func (s *SmartContract) SetChickenPrice(ctx contractapi.TransactionContextInterface, id string, price float64, owner string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner {
		chicken.Price = price
		chicken.TxType = "SetPrice"
		chickenAsBytes, _ := json.Marshal(chicken)

		_err := ctx.GetStub().PutState(id, chickenAsBytes)

		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		_chicken := new(Chicken)
		_ = json.Unmarshal(chickenAsBytes, _chicken)
		return _chicken, nil
	} else {
		return nil, fmt.Errorf("You are not owner of this asset")
	}

}

func (t *SmartContract) GetHistoryForAsset(ctx contractapi.TransactionContextInterface, chickenName string) ([]HistoryModel, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(chickenName)
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	results := []HistoryModel{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		chicken := new(Chicken)
		_ = json.Unmarshal(queryResponse.Value, chicken)

		historyItem := HistoryModel{
			TxId:      queryResponse.TxId,
			Value:     chicken,
			Timestamp: time.Unix(queryResponse.Timestamp.Seconds, int64(queryResponse.Timestamp.Nanos)).String(),
			IsDelete:  queryResponse.IsDelete}
		results = append(results, historyItem)
	}

	return results, nil

}

func (s *SmartContract) SetPublicToSell(ctx contractapi.TransactionContextInterface, id string, owner string) (*Chicken, error) {
	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner {
		chicken.ForSale = true
		chicken.TxType = "SetChickenPublic"
	}

	chickenAsBytes, _ := json.Marshal(chicken)

	_err := ctx.GetStub().PutState(id, chickenAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_chicken := new(Chicken)
	_ = json.Unmarshal(chickenAsBytes, _chicken)

	return _chicken, nil

}

func (s *SmartContract) SetBatchPublicToSell(ctx contractapi.TransactionContextInterface, id string, owner string) (*Batch, error) {
	batch, err := s.QueryBatch(ctx, id)

	if err != nil {
		return nil, err
	}

	if batch.Owner == owner {
		batch.ForSale = true
		batch.TxType = "SetBatchPublic"
	}

	batchAsBytes, _ := json.Marshal(batch)

	_err := ctx.GetStub().PutState(id, batchAsBytes)

	if _err != nil {
		return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
	}

	_batch := new(Batch)
	_ = json.Unmarshal(batchAsBytes, _batch)

	return _batch, nil

}

func (s *SmartContract) BlockingToken(ctx contractapi.TransactionContextInterface, customer string, price float64) (*Token, error) {
	token, err := s.QueryToken(ctx, customer)

	if err != nil {
		return nil, err
	}

	if token.Amount >= price {
		token.Amount = token.Amount - price
		token.BlockAmount = token.BlockAmount + price
		tokenAsBytes, _ := json.Marshal(token)
		_err := ctx.GetStub().PutState(customer, tokenAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		return token, nil
	} else {
		return nil, fmt.Errorf("Insufficient balance. %s", err.Error())
	}

}

func (s *SmartContract) SellChicken(ctx contractapi.TransactionContextInterface, id string, owner string, customer string, price float64, biders string, bids string) (*Chicken, error) {

	chicken, err := s.QueryChicken(ctx, id)

	if err != nil {
		return nil, err
	}

	customer_token, err := s.QueryToken(ctx, customer)

	if err != nil {
		return nil, err
	}

	owner_token, err := s.QueryToken(ctx, owner)

	if err != nil {
		return nil, err
	}

	if chicken.Owner == owner {

		customer_token.BlockAmount = customer_token.BlockAmount - price
		cTokenAsBytes, _ := json.Marshal(customer_token)
		_err_ := ctx.GetStub().PutState(customer, cTokenAsBytes)
		if _err_ != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		owner_token.Amount = owner_token.Amount + price
		oTokenAsBytes, _ := json.Marshal(owner_token)
		_err__ := ctx.GetStub().PutState(owner, oTokenAsBytes)
		if _err__ != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		chicken.Owner = customer
		chicken.TxType = "SellChicken"
		chicken.ForSale = false

		chickenAsBytes, _ := json.Marshal(chicken)
		_err := ctx.GetStub().PutState(id, chickenAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		if biders != "" {
			bidersArr := strings.Split(biders, "#")
			bidsArr := strings.Split(bids, "#")
	
			for i := 0; i < len(bidersArr); i++ {
	
				bid, _ := strconv.ParseFloat(bidsArr[i], 8)
	
				token, err := s.QueryToken(ctx, bidersArr[i])
				token.BlockAmount = token.BlockAmount - bid
				token.Amount = token.Amount + bid
				tokenAsBytes, _ := json.Marshal(token)
				_err := ctx.GetStub().PutState(bidersArr[i], tokenAsBytes)
				if _err != nil {
					return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
				}
	
			}
		}
		                                                                                 
		_chicken := new(Chicken)
		_ = json.Unmarshal(chickenAsBytes, _chicken)
		return _chicken, nil

	} else {
		return nil, fmt.Errorf("You are not owner of this asset. %s", err.Error())
	}

}

func (s *SmartContract) SellBatch(ctx contractapi.TransactionContextInterface, id string, owner string, customer string, price float64, biders string, bids string) (*Batch, error) {

	batch, err := s.QueryBatch(ctx, id)

	if err != nil {
		return nil, err
	}

	customer_token, err := s.QueryToken(ctx, customer)

	if err != nil {
		return nil, err
	}

	owner_token, err := s.QueryToken(ctx, owner)

	if err != nil {
		return nil, err
	}

	if batch.Owner == owner {

		customer_token.BlockAmount = customer_token.BlockAmount - price
		cTokenAsBytes, _ := json.Marshal(customer_token)
		_err_ := ctx.GetStub().PutState(customer, cTokenAsBytes)
		if _err_ != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		owner_token.Amount = owner_token.Amount + price
		oTokenAsBytes, _ := json.Marshal(owner_token)
		_err__ := ctx.GetStub().PutState(owner, oTokenAsBytes)
		if _err__ != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
		batch.Owner = customer
		batch.TxType = "SellBatch"
		batch.ForSale = false

		batchAsBytes, _ := json.Marshal(batch)
		_err := ctx.GetStub().PutState(id, batchAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		if biders != "" {
			bidersArr := strings.Split(biders, "#")
			bidsArr := strings.Split(bids, "#")
	
			for i := 0; i < len(bidersArr); i++ {
	
				bid, _ := strconv.ParseFloat(bidsArr[i], 8)
	
				token, err := s.QueryToken(ctx, bidersArr[i])
				token.BlockAmount = token.BlockAmount - bid
				token.Amount = token.Amount + bid
				tokenAsBytes, _ := json.Marshal(token)
				_err := ctx.GetStub().PutState(bidersArr[i], tokenAsBytes)
				if _err != nil {
					return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
				}
	
			}
		}
		                                                                                 
		_batch := new(Batch)
		_ = json.Unmarshal(batchAsBytes, _batch)
		return _batch, nil

	} else {
		return nil, fmt.Errorf("You are not owner of this asset. %s", err.Error())
	}

}

func (s *SmartContract) PutAssetsInBatch(ctx contractapi.TransactionContextInterface, assetsIds string, owner string, batchId string) (*Batch, error) {

	batch, err := s.QueryBatch(ctx, batchId)

	if err != nil {
		return nil, err
	}

	if batch.Owner == owner {

		assetsArr := strings.Split(assetsIds, "#")
		for i := 0; i < len(assetsArr); i++ {
			chicken, err := s.QueryChicken(ctx, assetsArr[i])
			if err != nil {
				return nil, err
			}
			if chicken.Owner == owner {
				chicken.Owner = batchId
				chicken.TxType = "PutInBatch"
				chickenAsBytes, _ := json.Marshal(chicken)
				_err := ctx.GetStub().PutState(assetsArr[i], chickenAsBytes)
				if _err != nil {
					return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
				}
				batch.Count = batch.Count + 1
			}
		}

		batch.TxType = "AddAssets"

		batchAsBytes, _ := json.Marshal(batch)
		_err := ctx.GetStub().PutState(batchId, batchAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		_batch := new(Batch)
		_ = json.Unmarshal(batchAsBytes, _batch)
		return _batch, nil		

	} else {
		return nil, fmt.Errorf("You are not owner of this asset. %s", err.Error())
	}

}

func (s *SmartContract) RemoveAssetsFromBatch(ctx contractapi.TransactionContextInterface, assetsIds string, owner string, batchId string) (*Batch, error) {

	batch, err := s.QueryBatch(ctx, batchId)

	if err != nil {
		return nil, err
	}

	if batch.Owner == owner {

		assetsArr := strings.Split(assetsIds, "#")
		for i := 0; i < len(assetsArr); i++ {
			chicken, err := s.QueryChicken(ctx, assetsArr[i])
			if err != nil {
				return nil, err
			}
			if chicken.Owner == batchId {
				chicken.Owner = owner
				chicken.TxType = "RemoveFromBatch"
				chickenAsBytes, _ := json.Marshal(chicken)
				_err := ctx.GetStub().PutState(assetsArr[i], chickenAsBytes)
				if _err != nil {
					return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
				}
				batch.Count = batch.Count - 1
			}
		}

		batch.TxType = "RemoveAssets"

		batchAsBytes, _ := json.Marshal(batch)
		_err := ctx.GetStub().PutState(batchId, batchAsBytes)
		if _err != nil {
			return nil, fmt.Errorf("Failed to put to world state. %s", err.Error())
		}

		_batch := new(Batch)
		_ = json.Unmarshal(batchAsBytes, _batch)
		return _batch, nil		

	} else {
		return nil, fmt.Errorf("You are not owner of this asset. %s", err.Error())
	}

}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting broilerChickenCC Smart Contract: %s", err.Error())
	}

}