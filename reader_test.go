package elecreader_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/pat-rohn/elecreader"
)

func TestMain(t *testing.T) {
	res, err := elecreader.Extract(GetExample())
	if err != nil {
		log.Fatalf("Failed to extract result: %v", err)
	}
	fmt.Printf("Res %v", res)
}

func GetExample() string {
	answer := "F.F(00)" + "\r\n"
	answer += "0.0(            1844)" + "\r\n"
	answer += "C.1.0(14593907)" + "\r\n"
	answer += "C.1.1(        )" + "\r\n"
	answer += "1.8.1(060740.931*kWh)" + "\r\n"
	answer += "1.8.2(093109.192*kWh)" + "\r\n"
	answer += "2.8.1(000000.000*kWh" + "\r\n"
	answer += "2.8.2(000000.000*kWh)" + "\r\n"
	answer += "1.8.0(153850.123*kWh)" + "\r\n"
	answer += "2.8.0(000000.000*kWh)" + "\r\n"
	answer += "15.8.0(153850.123*kWh)" + "\r\n"
	answer += "C.7.0(0035)" + "\r\n"
	answer += "32.7(233*V)" + "\r\n"
	answer += "52.7(233*V)" + "\r\n"
	answer += "72.7(234*V)" + "\r\n"
	answer += "31.7(000.62*A)" + "\r\n"
	answer += "51.7(000.42*A)" + "\r\n"
	answer += "71.7(001.29*A)" + "\r\n"
	answer += "82.8.1(0000)" + "\r\n"
	answer += "82.8.2(0000)" + "\r\n"
	answer += "0.2.0(M26)" + "\r\n"
	answer += "C.5.0(0400)" + "\r\n"
	answer += "!" + "\x03" // etx

	return answer
}
