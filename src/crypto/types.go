package crypto

import "crypto/tls"

type HostCertAndOrPrivateKey struct {
	HostTLSKey *tls.Certificate
}

type CryptoStore struct {
	localhostIdentPairs []*HostCertAndOrPrivateKey
	localhostTLSCert    *tls.Certificate
}

type VmCryptoStore struct {
	*CryptoStore
}
