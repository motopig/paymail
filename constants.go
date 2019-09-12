package paymail

const Version = "1.0"

type NotFound struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Capability Discovery is the process by which a paymail client learns the supported features of a paymail service and their respective endpoints and configurations
// http://bsvalias.org/02-02-capability-discovery.html#setup
type ServiceDiscoveryResponse struct {
	Version      string `json:"bsvalias"`
	Capabilities `json:"capabilities"`
}

type Capabilities struct {
	PkiUrl                  string `json:"pki"`
	PaymentDestinationUrl   string `json:"paymentDestination"`
	SenderValidationUrl     string `json:"6745385c3fc0"`
	VerifyPublicKeyOwnerUrl string `json:"a9f510c16bde"`
}

// The capabilities.pki path returns a URI template. Clients should replace the {alias} and {domain.tld} template parameters and then make an HTTP GET request against this URI.
// http://bsvalias.org/03-public-key-infrastructure.html#client-request
type PKIResponse struct {
	Version string `json:"bsvalias"`
	Handle  string `json:"handle"`
	PubKey  string `json:"pubkey"`
}

// This capability allows clients to verify if a given public key is a valid identity key for a given paymail handle.
// http://bsvalias.org/05-verify-public-key-owner.html#verify-public-key-owner
type VerifyResponse struct {
	Handle string `json:"handle"`
	Pubkey string `json:"pubkey"`
	Match  bool   `json:"match"`
}

// The capabilities.pki path returns a URI template. Senders should replace the {alias} and {domain.tld} template parameters with the values from the receiver's paymail handle and then make an HTTP POST request against this URI
// http://bsvalias.org/04-01-basic-address-resolution.html#sender-request
type BasicAddressRequest struct {
	SenderName   string  `json:"sendername"`
	SenderHandle string  `json:"senderhandle"`
	Dt           string  `json:"dt"`
	Amount       float64 `json:"amount"`
	Purpose      string  `json:"purpose"`
	Signature    string  `json:"signature"`
}
