package visasample

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

func TestPullFundsTransactionPost(t *testing.T) {
	cases := []struct {
		systemsTraceAuditNumber        int
		retrievalReferenceNumber       string
		localTransactionDateTime       string
		acquiringBin                   int
		acquirerCountryCode            int
		senderPrimaryAccountNumber     string
		senderCardExpiryDate           string
		senderCurrencyCode             string
		amount                         float64
		surcharge                      float64
		cavv                           string
		foreignExchangeFeeTransaction  float64
		businessApplicationId          string
		merchantCategoryCode           int
		CAname                         string
		CAterminalId                   string
		CAidCode                       string
		CAAstate                       string
		CAAcounty                      string
		CAAcountry                     string
		CAAzipCode                     string
		MSDtrack1Data                  string
		MSDtrack2Data                  string
		POSDpanEntryMode               int
		POSDposConditionCode           int
		POSDmotoECIIndicator           string
		POSCposTerminalType            int
		POSCposTerminalEntryCapability int
		PDpinDataBlock                 string
		PDSRCIpinBlockFormatCode       int
		PDSRCIzoneKeyIndex             int
		feeProgramIndicator            string
	}{
		//{},
		{
			123456,
			"407509300259",
			time.Now().Format("2006-01-02T03:04:05"),
			409999,
			101,
			"4957030100009952",
			"2020-03",
			"USD",
			110.,
			2.00,
			"0000010926000071934977253000000000000000",
			10.00,
			"AA",
			6012,
			"Saranya",
			"365539",
			"VMT200911026070",
			"CA",
			"081",
			"USA",
			"94404",
			"", //"1010101010101010101010101010",
			"",
			90,
			0,
			"0",
			4,
			2,
			"",
			0,
			0,
			"123",
		},
	}

	for _, c := range cases {
		cardAcceptorAddress := CardAcceptorAddress{
			State:   c.CAAstate,
			County:  c.CAAcounty,
			Country: c.CAAcountry,
			ZipCode: c.CAAzipCode,
		}
		cardAcceptor := CardAcceptor{
			Name:       c.CAname,
			TerminalId: c.CAterminalId,
			IdCode:     c.CAidCode,
			Address:    cardAcceptorAddress,
		}
		/*magneticStripeData := MagneticStripeData{
			Track1Data: c.MSDtrack1Data,
			Track2Data: c.MSDtrack2Data,
		}*/
		/*
			pointOfServiceData := PointOfServiceData{
				PanEntryMode:     c.POSDpanEntryMode,
				PosConditionCode: c.POSDposConditionCode,
				MotoECIIndicator: c.POSDmotoECIIndicator,
			}
		*/
		/*
			pointOfServiceCapability := PointOfServiceCapability{
				PosTerminalType:            c.POSCposTerminalType,
				PosTerminalEntryCapability: c.POSCposTerminalEntryCapability,
			}
		*/
		/*
			securityRelatedControlInfo := SecurityRelatedControlInfo{
				PinBlockFormatCode: c.PDSRCIpinBlockFormatCode,
				ZoneKeyIndex:       c.PDSRCIzoneKeyIndex,
			}
		*/
		/*
			pinData := PinData{
				PinDataBlock: c.PDpinDataBlock,
				//SecurityRelatedControlInfo: securityRelatedControlInfo,
			}
		*/
		request := PullFundsTransactionRequest{
			SystemsTraceAuditNumber:    c.systemsTraceAuditNumber,
			RetrievalReferenceNumber:   c.retrievalReferenceNumber,
			LocalTransactionDateTime:   c.localTransactionDateTime,
			AcquiringBin:               c.acquiringBin,
			AcquirerCountryCode:        c.acquirerCountryCode,
			SenderPrimaryAccountNumber: c.senderPrimaryAccountNumber,
			SenderCardExpiryDate:       c.senderCardExpiryDate,
			SenderCurrencyCode:         c.senderCurrencyCode,
			Amount:                     c.amount,
			Surcharge:                  c.surcharge,
			Cavv:                       c.cavv,
			ForeignExchangeFeeTransaction: c.foreignExchangeFeeTransaction,
			BusinessApplicationId:         c.businessApplicationId,
			MerchantCategoryCode:          c.merchantCategoryCode,
			CardAcceptor:                  cardAcceptor,
			//MagneticStripeData:            &magneticStripeData,
			//PointOfServiceData:       &pointOfServiceData,
			//PointOfServiceCapability: &pointOfServiceCapability,
			//PinData:                  &pinData,
			FeeProgramIndicator: c.feeProgramIndicator,
		}

		// Set UUID
		newUuid := uuid.NewV4()
		uuid := newUuid.String()

		// Set user/password
		SetUserPassword("user_id", "user_password")
		//SetCertPaths("pubKeyPath", "pvtKeyPath", "caPvtKeyPath")

		response, err := PullFundsTransactionsPost(request, uuid)
		if err != nil {
			t.Errorf("Error when getting response: %v", err)
		}
		fmt.Printf("%+v\n", response)
		// 1. Check type
		if reflect.TypeOf(response).String() != "visasample.PullFundsTransactionResponse" {
			t.Errorf("Return should be of type PullFundsTransactionResponse. Looking for %s, got %s", "visa.PullFundsTransactionResponse", reflect.TypeOf(response).String())
		}
	}
}
