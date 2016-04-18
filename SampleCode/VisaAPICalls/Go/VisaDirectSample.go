package visasample

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const USER_ID = ""
const USER_PASSWORD = ""

const SSL_PUBLIC_KEY_PATH = "certs/application.crt"
const SSL_PRIVATE_KEY_PATH = "certs/application.pem"
const SSL_CAPRIVATE_KEY_PATH = "certs/VDPCA-SBX.pem"

const API_SANDBOX = "https://sandbox.api.visa.com"
const API_PORT = 433

var PULL_FUNDS_TRANSACTIONS_URL = API_SANDBOX + "/visadirect/fundstransfer/v1/pullfundstransactions/"

type PullFundsTransactionRequest struct {
	SystemsTraceAuditNumber       int                       `json:"systemsTraceAuditNumber"`                 // required, 6
	RetrievalReferenceNumber      string                    `json:"retrievalReferenceNumber"`                // ydddhhnnnnnn(numeric characters only), Length: 12
	LocalTransactionDateTime      string                    `json:"localTransactionDateTime"`                // RFC3339. dateTime | YYYY-MM-DDThh:mm:ss. The date and time you submit the transaction
	AcquiringBin                  int                       `json:"acquiringBin"`                            // integer | positive, Length: 6 - 11
	AcquirerCountryCode           int                       `json:"acquirerCountryCode"`                     // integer | Length: 3
	SenderPrimaryAccountNumber    string                    `json:"senderPrimaryAccountNumber"`              // string | Length: 13 - 19
	SenderCardExpiryDate          string                    `json:"senderCardExpiryDate"`                    // string | YYYY-MM
	SenderCurrencyCode            string                    `json:"senderCurrencyCode"`                      // string | Length: 3
	Amount                        float64                   `json:"amount,omitempty"`                        // Optional: decimal | Length: totalDigits 12, fractionDigits 3 (minimum value is 0)
	Surcharge                     float64                   `json:"surcharge,omitempty"`                     // Optional: decimal | Length: totalDigits 12, fractionDigits 3(minimum value is 0)
	Cavv                          string                    `json:"cavv"`                                    // string | Length:40
	ForeignExchangeFeeTransaction float64                   `json:"foreignExchangeFeeTransaction,omitempty"` // Optional: decimal | Length: totalDigits 12, fractionDigits 3 (minimum value is 0)
	BusinessApplicationId         string                    `json:"businessApplicationId"`                   // string | Length: 2
	MerchantCategoryCode          int                       `json:"merchantCategoryCode,omitempty"`          // Conditional: integer | Length: total 4 digits
	CardAcceptor                  CardAcceptor              `json:"cardAcceptor"`                            // Object
	MagneticStripeData            *MagneticStripeData       `json:"magneticStripeData,omitempty"`            // Optional: Object
	PointOfServiceData            *PointOfServiceData       `json:"pointOfServiceData,omitempty"`            // Conditional: Object
	PointOfServiceCapability      *PointOfServiceCapability `json:"pointOfServiceCapability,omitempty"`      // Conditional: Object
	PinData                       *PinData                  `json:"pinData,omitempty"`                       // Conditional: Object
	FeeProgramIndicator           string                    `json:"feeProgramIndicator,omitempty"`           // Optional: string | Length:3
}

// --- START Element Structs ---
type CardAcceptor struct {
	Name       string              `json:"name"`       // string | Length: 1 - 25
	TerminalId string              `json:"terminalId"` // string | Length: 1 - 8
	IdCode     string              `json:"idCode"`     // string | Length: 1 - 15
	Address    CardAcceptorAddress `json:"address"`    // Object
}

type CardAcceptorAddress struct {
	State   string `json:"state,omitempty"`   // Conditional: string | Length: 2
	County  string `json:"county,omitempty"`  // Conditional: string | Length: 3
	Country string `json:"country"`           // string | Length: 3
	ZipCode string `json:"zipCode,omitempty"` // Conditional: string | Length: 5 - 9
}

type MagneticStripeData struct {
	Track1Data string `json:"track1Data,omitempty"` // Conditional: string | Length: maximum 76
	Track2Data string `json:"track2Data,omitempty"` // Conditional: string | hex binary value is sent as String, Length: maximum 19
}

type PointOfServiceData struct {
	PanEntryMode     int    `json:"panEntryMode,omitempty"`     // Conditional: integer | positive, Length: totaldigits 2
	PosConditionCode int    `json:"posConditionCode,omitempty"` // Conditional: integer | positive,Length: totalDigits 2
	MotoECIIndicator string `json:"motoECIIndicator,omitempty"` // Conditional: string | Length: 1 , max: 1 characters
}

type PointOfServiceCapability struct {
	PosTerminalType            int `json:"posTerminalType,omitempty"`            // Conditional: integer | positive, totalDigits 0
	PosTerminalEntryCapability int `json:"posTerminalEntryCapability,omitempty"` // Conditional: integer | positive, Length: totalDigits 1
}

type PinData struct {
	PinDataBlock               string                      `json:"pinDataBlock,omitempty"`               // Conditional: string | Length: 16, hexbinary format
	SecurityRelatedControlInfo *SecurityRelatedControlInfo `json:"securityRelatedControlInfo,omitempty"` // Conditional: object
}

type SecurityRelatedControlInfo struct {
	PinBlockFormatCode int `json:"pinBlockFormatCode,omitempty"` // Conditional: integer |positive Length: totalDigits 2
	ZoneKeyIndex       int `json:"zoneKeyIndex,omitempty"`       // Conditional: integer |positive Length: totalDigits 2
}

// --- END Element Structs ---

// --- START Response Structs ---
type PullFundsTransactionResponse struct {
	StatusIdentifier      string `json:"statusIdentifier"`              // string | required when call times out
	TransactionIdentifier int    `json:"transactionIdentifier"`         // integer | positive and required when call does not timeout, Length: 15
	ActionCode            string `json:"actionCode"`                    // string | Length: 2
	ApprovalCode          string `json:"ApprovalCode,omitempty"`        // Optional: string | Length: 6
	TransmissionDateTime  string `json:"transmissionDateTime"`          // dateTime | YYYY-MM-DDThh:mm:ss
	CavvResultCode        string `json:"cavvResultCode,omitempty"`      // Optional: string | Length: 1
	ResponseCode          string `json:"responseCode"`                  // string | Length: 1
	FeeProgramIndicator   string `json:"feeProgramIndicator,omitempty"` // Optional: string | Length:3
	ErrorMessage          string `json:"errorMessage,omitempty"`        // Optional: string | Length:3
}

// PullFundsTransactions (POST) Resource debits (pulls) funds from a sender's Visa account (in preparation for pushing funds to a recipient's account)
// by initiating a financial message called an Account Funding Transaction (AFT)
func PullFundsTransactionsPost(request PullFundsTransactionRequest, uuid string) (response PullFundsTransactionResponse, err error) {
	/*
	   You should log or otherwise retain all the information returned in the PullFundsTransactions response.
	   Should it be necessary to initiate a ReverseFundsTransactions POST operation, you may need to provide
	   several of the PullFundsTransactions Response values in the Request.
	*/
	body, err := json.Marshal(request)
	if err != nil {
		return response, err
	}
	responseJson, err := Client(USER_ID, USER_PASSWORD, PULL_FUNDS_TRANSACTIONS_URL, "POST", false, body, uuid)
	if err != nil {
		return response, err
	}
	// Unmarshall response
	err = json.Unmarshal(responseJson, &response)
	if err != nil {
		return response, err
	}
	return
}

func Client(userId string, userPassword string, url string, reqType string, production bool, body []byte, transactionID string) (response []byte, err error) {
	authHeader := createAuthHeader()

	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(body))
	req.Header.Set("X-Client-Transaction-ID", transactionID)
	req.Header.Set("Authorization:Basic ", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Load client cert
	cert, err := tls.LoadX509KeyPair(SSL_PUBLIC_KEY_PATH, SSL_PRIVATE_KEY_PATH)
	if err != nil {
		return nil, fmt.Errorf("Could not load key pair: %v", err)
	}
	// Load CA cert
	caCert, err := ioutil.ReadFile(SSL_CAPRIVATE_KEY_PATH)
	if err != nil {
		return nil, fmt.Errorf("Could not load CA key: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: true, //@FIXME: This call *must* be secure
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not load HTTPS client: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		response, _ = ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP error: %s, %v", resp.Status, string(response))
	}

	response, _ = ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %v\n", string(response))
	return response, nil
}

func createAuthHeader() (authHeader string) {
	// Auth header = \ase64(userid:user_password)
	authHeader = base64.StdEncoding.EncodeToString([]byte(USER_ID + ":" + USER_PASSWORD))
	return
}
