package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
}

func New(endpoint string) API {
	impl := &APIImpl{
		APIBase: APIBase{endpoint: endpoint},
	}

	impl.GeneralImp.Base = impl.APIBase
	impl.BlocksImpl.Base = impl.APIBase
	impl.AccountsImpl.Base = impl.APIBase
	impl.StateImpl.Base = impl.APIBase
	impl.EventsImpl.Base = impl.APIBase
	impl.TransactionsImpl.Base = impl.APIBase
	return impl
}

type APIImpl struct {
	APIBase

	GeneralImp
	BlocksImpl
	TransactionsImpl
	AccountsImpl
	EventsImpl
	StateImpl
}

type APIBase struct {
	endpoint string
}

func (impl APIBase) Endpoint() string {
	return impl.endpoint
}

type Base interface {
	Endpoint() string
}

type API interface {
	General
	Blocks
	Transactions
	Accounts
	Events
	State
}

type ResponseHeader struct {
	AptosBlockHeight         uint64
	AptosChainID             uint16
	AptosEpoch               uint64
	AptosLedgerOldestVersion uint64
	AptosLedgerTimestampusec uint64
	AptosLedgerVersion       uint64
	AptosOldestBlockHeight   uint64
}

func Request(method, endpoint string, reqBody, resp interface{},
	query map[string]interface{}, respHeader *ResponseHeader) error {
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}

	if req.URL != nil && query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Add("Content-Type", "application/json")
	req = req.WithContext(context.Background())
	rsp, err := client.Do(req)
	if err != nil {
		return err
	}

	rspBody, err := ioutil.ReadAll(rsp.Body)
	defer rsp.Body.Close()
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK && rsp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("response(%d): %s", rsp.StatusCode, string(rspBody))
	}

	err = json.Unmarshal(rspBody, resp)
	if err != nil {
		return err
	}

	if respHeader != nil {
		if len(rsp.Header["X-Aptos-Block-Height"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Block-Height"][0], 10)
			respHeader.AptosBlockHeight = v.Uint64()
		}

		if len(rsp.Header["X-Aptos-Chain-Id"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Chain-Id"][0], 10)
			respHeader.AptosChainID = uint16(v.Uint64())
		}

		if len(rsp.Header["X-Aptos-Epoch"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Epoch"][0], 10)
			respHeader.AptosEpoch = v.Uint64()
		}

		if len(rsp.Header["X-Aptos-Ledger-Oldest-Version"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Ledger-Oldest-Version"][0], 10)
			respHeader.AptosLedgerOldestVersion = v.Uint64()
		}

		if len(rsp.Header["X-Aptos-Ledger-Timestampusec"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Ledger-Timestampusec"][0], 10)
			respHeader.AptosLedgerTimestampusec = v.Uint64()
		}

		if len(rsp.Header["X-Aptos-Ledger-Version"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Ledger-Version"][0], 10)
			respHeader.AptosLedgerVersion = v.Uint64()
		}

		if len(rsp.Header["X-Aptos-Oldest-Block-Height"]) > 0 {
			v, _ := new(big.Int).SetString(rsp.Header["X-Aptos-Oldest-Block-Height"][0], 10)
			respHeader.AptosOldestBlockHeight = v.Uint64()
		}
	}
	return nil
}

func requestOptions(opts ...interface{}) (rspHeader *ResponseHeader) {
	for _, opt := range opts {
		switch opt.(type) {
		case *ResponseHeader:
			rspHeader = opt.(*ResponseHeader)
		}
	}
	return
}
