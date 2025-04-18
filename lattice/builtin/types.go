package builtin

// Contract builtin contract struct
type Contract struct {
	AbiString   string `json:"abiString,omitempty"`
	Address     string `json:"address,omitempty"`
	Description string `json:"description,omitempty"`
}

type UploadKeyParam struct {
	PublicKey []byte `json:"publicKey,omitempty"`
	Address   string
}
