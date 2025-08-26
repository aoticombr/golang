package http

type TCert struct {
	PathCrt      string
	PathPriv     string
	CertPEMBlock []byte
	KeyPEMBlock  []byte
	PfxBlock     []byte
	PfxPass      string
}
