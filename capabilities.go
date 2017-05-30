package twilio

import (
	"fmt"
	"strings"
	"time"

	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/jonboulle/clockwork"
)

// Clock provides time for calculating expiry times in tokens.
var Clock = clockwork.NewRealClock()

// Capabilities describes the possible capabilities for a Twilio token.
type Capabilities struct {
	AccountSid          string // AccountSid is the Twilio Account Sid
	AuthToken           string // AuthToken is the secret Twilio Account Auth Token
	AllowClientIncoming string // AllowClientIncoming specifies the possible client name that can be adopted
	AllowClientOutgoing string // AllowClientOutgoing specifies that this token can make outgoing calls using the given Application Sid
}

// Generate creates a Capabilities Token given some configuration values.
// See https://www.twilio.com/docs/api/client/capability-tokens for details.
func Generate(c Capabilities, expires time.Duration) (string, error) {
	sig, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: []byte(c.AuthToken)}, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		return "", err
	}

	cl := jwt.Claims{
		Issuer: c.AccountSid,
		Expiry: jwt.NewNumericDate(Clock.Now().Add(expires)),
	}

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
	//return jwt.Signed(sig).Claims(cl).CompactSerialize()
	return jwt.Signed(sig).Claims(cl).Claims(map[string]interface{}{
		"scope": strings.Join(scopes, " "),
	}).CompactSerialize()
}
