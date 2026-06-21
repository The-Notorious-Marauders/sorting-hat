package properties

import (
	"github.com/golibs-starter/golib/config"
)

// JwksProperties represents ...
type JwksProperties struct {
	PublicKey  string `validate:"required"`
	PrivateKey string `validate:"required"`
}

// NewJwksProperties return a new JwksProperties instance
func NewJwksProperties(loader config.Loader) (*JwksProperties, error) {
	props := JwksProperties{}
	err := loader.Bind(&props)
	return &props, err
}

// Prefix return config prefix
func (t *JwksProperties) Prefix() string {
	return "app.jwks"
}
