package twilio_test

import (
	"fmt"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/tmc/twilio"
)

func init() {
	twilio.Clock = clockwork.NewFakeClockAt(time.Unix(1257894000, 0))
}

func ExampleGenerate() {
	caps := twilio.Capabilities{
		AccountSid:          "ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
		AuthToken:           "yyyyyyyyyyyyyyyyyyyyyyyyyyyyy",
		AllowClientIncoming: "tommy",
		AllowClientOutgoing: "APFOOOOOOOo",
	}

	token, err := twilio.Generate(caps, time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	// output:
	// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjEyNTc4OTQwMDEsImlzcyI6IkFDeHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHh4eHgiLCJzY29wZSI6InNjb3BlOmNsaWVudDpvdXRnb2luZz9hcHBTaWQ9QVBGT09PT09PT29cdTAwMjZjbGllbnROYW1lPXRvbW15IHNjb3BlOmNsaWVudDppbmNvbWluZz9jbGllbnROYW1lPXRvbW15In0.n-GADTApMTBanP_o69br2djf8GSmycaL3FpLYHHcrTA
}
