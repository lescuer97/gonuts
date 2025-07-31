package nut18

import (
	"fmt"
	"testing"
)

func TestDecodingPaymentReq(t *testing.T) {

	encodedPayReq := "creqApmF0gaNhdGVub3N0cmFheKlucHJvZmlsZTFxeTI4d3VtbjhnaGo3dW45ZDNzaGp0bnl2OWtoMnVld2Q5aHN6OW1od2RlbjV0ZTB3ZmprY2N0ZTljdXJ4dmVuOWVlaHFjdHJ2NWhzenJ0aHdkZW41dGUwZGVoaHh0bnZkYWtxcWd5dnB6NXlzajBkcXgzZHpwdjg1eHdscmFwZncwOTR3c3EwdDdkeHd6cHl6eXAwem0zMGd1dWV6Zng1YWeBgmExZk5JUC0wNGFpanBheW1lbnRfaWRhYQ1hdWNzYXRhbYF4Imh0dHBzOi8vbm9mZWVzLnRlc3RudXQuY2FzaHUuc3BhY2VhZHB0aGlzIGlzIHRoZSBtZW1v"

	payReq, err := DecodePaymentRequest(encodedPayReq)

	if err != nil {
		t.Fatalf("DecodePaymentRequest(encodedPayReq) %+v", err)
	}

	fmt.Printf("payment req: %+v", payReq)
}
