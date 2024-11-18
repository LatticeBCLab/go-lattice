package key

import (
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"github.com/btcsuite/btcutil/base58"
	"github.com/rs/zerolog/log"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
	"golang.org/x/crypto/scrypt"
	"hash"
)

const defaultRandomBytesSize = 32
const defaultKeySize = 16 // 28 or 32

type OptFunc func(*Opts)
type HashFunc func() hash.Hash
type SymEncryptFunc func(data, key []byte) ([]byte, error)
type FormatFunc func([]byte) string

// Opts is the options for the ECDH exchange.
type Opts struct {
	curve          ecdh.Curve
	hashFunc       HashFunc
	symEncryptFunc SymEncryptFunc
	formatFunc     FormatFunc
	keySize        int
}

// defaultOpts returns the default options for the ECDH exchange.
func defaultOpts() *Opts {
	return &Opts{
		curve:          ecdh.P256(),
		hashFunc:       sm3.New,
		symEncryptFunc: Sm4ECB,
		formatFunc:     base58.Encode,
		keySize:        defaultKeySize,
	}
}

// WithCurve returns an OptFunc that sets the curve for the ECDH exchange.
func WithCurve(curve ecdh.Curve) OptFunc {
	return func(o *Opts) {
		o.curve = curve
	}
}

// WithHashFunc returns an OptFunc that sets the hash function for the ECDH exchange.
func WithHashFunc(h HashFunc) OptFunc {
	return func(o *Opts) {
		o.hashFunc = h
	}
}

// WithSymEncryptFunc returns an OptFunc that sets the symmetric encryption function for the ECDH exchange.
func WithSymEncryptFunc(fn SymEncryptFunc) OptFunc {
	return func(o *Opts) {
		o.symEncryptFunc = fn
	}
}

// WithFormatFunc returns an OptFunc that sets the format function for the ECDH exchange.
func WithFormatFunc(fn FormatFunc) OptFunc {
	return func(o *Opts) {
		o.formatFunc = fn
	}
}

// WithKeySize returns an OptFunc that sets the key size for the ECDH exchange.
func WithKeySize(size int) OptFunc {
	if size < defaultKeySize {
		log.Warn().Msg("key size should be greater than 16")
		size = defaultKeySize
	} else if size > defaultKeySize*2 {
		log.Warn().Msg("key size should be less than 32")
		size = defaultKeySize
	} else if size%8 != 0 {
		log.Warn().Msg("key size should be a multiple of 8")
		size = defaultKeySize
	}
	return func(o *Opts) {
		o.keySize = size
	}
}

// ExchangeParams is the parameters for the ECDH exchange.
type ExchangeParams struct {
	PK     *ecdh.PublicKey
	Random []byte
}

// ExchangeResult is the result of the ECDH exchange.
type ExchangeResult struct {
	Random     []byte
	PK         *ecdh.PublicKey
	CipherText string
	AccessKey  *AccessKey
}

// AccessKey is the access key for the ECDH exchange.
type AccessKey struct {
	ID     string `json:"id"`
	Secret string `json:"secret"`
}

type Exchange interface {
	// GenerateSharedParams generates the shared parameters for the ECDH exchange.
	// The params are the public key and the random number.
	// Warning: the private key should be kept secret.
	GenerateSharedParams() (*ecdh.PrivateKey, []byte, error)
	// Exchange performs the ECDH exchange and returns the shared secret.
	// You should generate client exchange params before calling this method. Can use GenerateSharedParams to generate the params.
	Exchange(client *ExchangeParams) (*ExchangeResult, error)
	// ClientExchange simulates the client's ECDH exchange.
	ClientExchange(clientSK *ecdh.PrivateKey, clientRandom []byte, server *ExchangeParams) (*AccessKey, error)
	// ConfirmAccessKeyIdOrigin confirms the access key id whether it comes from the secret.
	ConfirmAccessKeyIdOrigin(id, secret string) error
}

type ECDHEExchange struct {
	opts *Opts
}

// NewECDHEExchange initializes a new ECDHEExchange.
func NewECDHEExchange(opts ...OptFunc) Exchange {
	o := defaultOpts()
	for _, fn := range opts {
		fn(o)
	}
	return &ECDHEExchange{
		opts: o,
	}
}

// GenerateSharedParams generates the shared parameters for the ECDH exchange.
func (e *ECDHEExchange) GenerateSharedParams() (*ecdh.PrivateKey, []byte, error) {
	sk, err := e.opts.curve.GenerateKey(rand.Reader)
	if err != nil {
		log.Err(err).Msg("failed to generate private key")
		return nil, nil, err
	}

	random, err := RandomBytes()
	if err != nil {
		log.Err(err).Msg("failed to generate random")
		return nil, nil, err
	}

	return sk, random, nil
}

func (e *ECDHEExchange) Exchange(client *ExchangeParams) (*ExchangeResult, error) {
	// 1. Generate server private key, generate server random number
	serverSK, serverRandom, err := e.GenerateSharedParams()
	if err != nil {
		return nil, err
	}

	// 2. Sign server's public key that use node certificate
	sharedKey, err := serverSK.ECDH(client.PK)
	if err != nil {
		log.Err(err).Msg("failed to calculate shared key")
		return nil, err
	}

	// 3. Calculate shared secret
	salt := append(client.Random, serverRandom...)
	sessionKey, err := scrypt.Key(sharedKey, salt, 1<<2, 1, 8, 32)
	if err != nil {
		log.Err(err).Msg("failed to calculate session key")
		return nil, err
	}
	sessionKey = sessionKey[:e.opts.keySize]

	// 4. Hash the session key to get the access_key_id
	h := e.opts.hashFunc()
	h.Write(sessionKey)
	accessKeyId := h.Sum(nil)
	accessKeyId = accessKeyId[:e.opts.keySize]

	// 5. Symmetric encrypt the hashed result
	h2 := e.opts.hashFunc()
	h2.Write(salt)
	data := h2.Sum(nil)
	cipherText, err := e.opts.symEncryptFunc(data, sessionKey[:16])
	if err != nil {
		log.Err(err).Msg("failed to encrypt data")
		return nil, err
	}

	// 6. Return the Exchange Result
	return &ExchangeResult{
		Random:     serverRandom,
		PK:         serverSK.PublicKey(),
		CipherText: e.opts.formatFunc(cipherText),
		AccessKey: &AccessKey{
			ID:     e.opts.formatFunc(accessKeyId),
			Secret: e.opts.formatFunc(sessionKey),
		},
	}, nil
}

func (e *ECDHEExchange) ClientExchange(clientSK *ecdh.PrivateKey, clientRandom []byte, server *ExchangeParams) (*AccessKey, error) {
	sharedKey, err := clientSK.ECDH(server.PK)
	if err != nil {
		log.Err(err).Msg("failed to calculate shared key")
		return nil, err
	}

	salt := append(clientRandom, server.Random...)
	sessionKey, err := scrypt.Key(sharedKey, salt, 1<<2, 1, 8, 32)
	if err != nil {
		log.Err(err).Msg("failed to calculate session key")
		return nil, err
	}
	sessionKey = sessionKey[:e.opts.keySize]

	h := e.opts.hashFunc()
	h.Write(sessionKey)
	accessKeyId := h.Sum(nil)
	accessKeyId = accessKeyId[:e.opts.keySize]

	return &AccessKey{
		ID:     e.opts.formatFunc(accessKeyId),
		Secret: e.opts.formatFunc(sessionKey),
	}, nil
}

func (e *ECDHEExchange) ConfirmAccessKeyIdOrigin(id, secret string) error {
	// decode the secret that encoded by the format function,
	// now we fixed using base58 to decode the secret
	sessionKey := base58.Decode(secret)
	sessionKey = sessionKey[:e.opts.keySize]

	h := e.opts.hashFunc()
	h.Write(sessionKey)
	accessKeyId := h.Sum(nil)
	accessKeyId = accessKeyId[:e.opts.keySize]
	if e.opts.formatFunc(accessKeyId) != id {
		return errors.New("invalid access key id")
	}
	return nil
}

// RandomBytes generates random bytes.
func RandomBytes() ([]byte, error) {
	b := make([]byte, defaultRandomBytesSize)
	if _, err := rand.Read(b); err != nil {
		log.Err(err).Msg("failed to generate random bytes")
		return nil, err
	}
	return b, nil
}

// Sm4ECB encrypts the data using SM4 ECB mode.
// key size can be 128 bits, or 192 bits, or 256 bits.
func Sm4ECB(data, key []byte) ([]byte, error) {
	cipher, err := sm4.Sm4Ecb(key, data, true)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}
