package gopaybtc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type NewAddressResponse struct {
	Address string `json:"address"`
}

type BlockonomicsClient struct {
	apiKey     string
	httpClient *http.Client
	addressURL string
	txURL      string
	newAddrURL string
}

type AddressInfo struct {
	Address           string `json:"address"`
	Balance           int64  `json:"balance"`
	TotalReceived     int64  `json:"total_received"`
	TotalTransactions int64  `json:"total_transactions"`
}

type TransactionInfo struct {
	Address string `json:"address"`
	Txid    string `json:"txid"`
	Value   int64  `json:"value"`
	Time    int64  `json:"time"`
}

type CreatePaymentRequest struct {
	Addr        string `json:"addr"`
	CallbackURL string `json:"callback_url"`
	Amount      int    `json:"amount"`
}

type CreatePaymentResponse struct {
	Address string `json:"address"`
	Status  string `json:"status"`
}

const (
	defaultAddressURL = "https://www.blockonomics.co/api/address?&addr=%s"
	defaultTxURL      = "https://www.blockonomics.co/api/searchhistory?txid=%s"
	defaultNewAddrURL = "https://www.blockonomics.co/api/new_address"
)

func NewBlockonomicsClient(config Config) *BlockonomicsClient {
	return &BlockonomicsClient{
		apiKey: config.ApiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		addressURL: defaultAddressURL,
		txURL:      defaultTxURL,
		newAddrURL: defaultNewAddrURL,
	}
}

func (client *BlockonomicsClient) CreatePaymentRequest(amount int, callbackURL string) (*CreatePaymentResponse, error) {
	reqBody := &CreatePaymentRequest{
		Addr:        client.newAddrURL,
		CallbackURL: callbackURL,
		Amount:      amount,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, client.addressURL+"/payment", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var createPaymentResp CreatePaymentResponse
	err = json.NewDecoder(resp.Body).Decode(&createPaymentResp)
	if err != nil {
		return nil, err
	}

	return &createPaymentResp, nil
}

func (client *BlockonomicsClient) GetAddressInfo(address string) (*AddressInfo, error) {
	url := fmt.Sprintf(client.addressURL, address)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data AddressInfo
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (client *BlockonomicsClient) GetTransactionInfo(txid string) (*TransactionInfo, error) {
	url := fmt.Sprintf(client.txURL, txid)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data []TransactionInfo
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("transaction not found: %s", txid)
	}

	return &data[0], nil
}

func (client *BlockonomicsClient) NewAddress() (string, error) {
	url := client.newAddrURL

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "Bearer "+client.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var data NewAddressResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}
	return data.Address, nil
}
