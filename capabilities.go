package twilio

import (
	"fmt"
	"strings"
	"time"

	"github.com/coreos/go-oidc/jose"
	"github.com/jonboulle/clockwork"
)

// Clock provides time for calculating expiry times in tokesn.
var Clock = clockwork.NewRealClock()

// Capabilities describes the possible capabilities for a Twilio token.
type Capabilities struct {
	AccountSid          string // AccountSid is the Twilio Account Sid
	AuthToken           string // AuthToken is the secret Twilio Account Auth Token
	AllowClientIncoming string // AllowClientIncoming specifies the possible client name fthat can be adopted
	AllowClientOutgoing string // AllowClientOutgoing specifies that this token can make outgoing calls using the given Application Sid
}

// Generate creates a Capabilities Token given some configuration values.
// See https://www.twilio.com/docs/api/client/capability-tokens for details.
func Generate(c Capabilities, expires time.Duration) (string, error) {
	signer := jose.NewSignerHMAC("", []byte(c.AuthToken))
	claims := jose.Claims{}

	claims.Add("iss", c.AccountSid)
	claims.Add("exp", Clock.Now().Add(expires).Unix())
	scopes := []string{}
	if c.AllowClientOutgoing != "" {
		scope := fmt.Sprintf("scope:client:outgoing?appSid=%s", c.AllowClientOutgoing)
		if c.AllowClientIncoming != "" {
			scope += fmt.Sprintf("&clientName=%s", c.AllowClientIncoming)
		}
		scopes = append(scopes, scope)
	}
	if c.AllowClientIncoming != "" {
		scopes = append(scopes, fmt.Sprintf("scope:client:incoming?clientName=%s", c.AllowClientIncoming))
	}
	claims.Add("scope", strings.Join(scopes, " "))

	jwt, err := jose.NewSignedJWT(claims, signer)
	if err != nil {
		return "", err
	}
	return jwt.Encode(), nil
}
