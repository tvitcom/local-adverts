package util 
import (
	"time"
	"encoding/json"
    "github.com/dvsekhvalnov/jose2go"
)

// get real jwt by next params: seckey, tokid, fqdn, uid, appsid, roles
func MakeJwtString(secretString, tokid, fqdn, uid, appids, roles string) (string, error) {
	c := struct {
	    Role       string `json:"rol"`
	    IssuedAt   int64  `json:"iat"`
		ExpiresAt  int64  `json:"exp"`
		Issuer     string `json:"iss"`
	    Subject    string `json:"sub"`
	    Audience   string `json:"aud"`
	}{
		Role: 	   roles,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Minute*120).Unix(),
		Issuer:    fqdn,
	    Subject:   uid,//тоесть user_id
	    Audience:  appids,
	}
	payload, err := json.Marshal(&c)
	if err != nil {
	    return "", err
	}
	return jose.SignBytes(
			payload, 
			jose.HS256, 
			[]byte(secretString), 
			jose.Header("kid", tokid),
	        jose.Header("typ", "jwt"),
        )
}

func GetJwtClaimHMAC(tokenString, secret string) ([]byte, map[string]interface{}, error) {
	payload, headers, err := jose.DecodeBytes(tokenString, []byte(secret))
        // println("JWTSTR:", tokenString)
        // println("Headers:", headers)
        // println("ERRTOK:",err.Error())
	return payload, headers, err
}

