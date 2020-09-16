package utils

import (
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	jwtSignedSecret = `Tes9tinas2kmskajirn`
	Issuer          = `bayen`
	UPLOAD          = `UPLOAD`
	DOWNLOAD        = `DOWNLOAD`
	DELETE          = `DELETE`
	ANY             = `ANY`
)

var ACTIONTYPES = []string{UPLOAD, DOWNLOAD, DELETE, ANY}

// required in uploading if you want to calculate the volume.
type volumeStatisticClaims struct {
	// calculate the usrs's volume if set user subject.
	Subject string `json:"sub,omitempty"`

	// calculate the organization's volume if set organization id.
	Org int64 `json:"org,omitempty"`
}

// used by meerastorage to verify the AC if you want.
type acClaims struct {
	Subject  string `json:"sub,omitempty"`
	Resource string `json:"res,omitempty"`
	Action   string `json:"act,omitempty"`
	Context  string `json:"ctx,omitempty"`
}

type FileClaims struct {
	VolumeStatisticClaims *volumeStatisticClaims `json:"vs,omitempty"`

	AcClaims *acClaims `json:"ac,omitempty"`

	// check the user if has logined if true
	AuthRequired bool `json:"auth,omitempty"`

	// Action: UPLOAD, DOWNLOAD, DELETE, ANY
	Action string `json:"act,omitempty"`

	// max file size in bytes allowed to be uploaded.
	// 0 means no limit.
	MaxSize uint64 `json:"mx,omitempty"`

	jwt.StandardClaims
}

func GenClaims(subject string) FileClaims {
	return FileClaims{
		StandardClaims: jwt.StandardClaims{
			// required, used as the :service(also as the prefix of the key in S3).
			Issuer: Issuer,
			// required if AuthRequired is true
			Subject: subject,
		},
		// check if the user is logined in if true
		AuthRequired: true,
		// required
		// UPLOAD:   used when upload the file
		// DOWNLOAD: used when download the file
		// DELETE:   used when delete the file
		// ANY:      used no matter upload, download, or delete the file. NOT RECOMMEND!
		Action: UPLOAD,
	}
}

func GenFileToken(subject string) string {
	claims := GenClaims(subject)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(jwtSignedSecret))
	if err != nil {
		fmt.Println("generate file token error: ", err)
		return ""
	}
	return signedToken
}
