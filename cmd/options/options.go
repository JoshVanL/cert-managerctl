package options

type Flags struct {
	Kubeconfig string

	Request Request
}

type Request struct {
	Cert Cert
	Sign Sign
}

type Cert struct {
	CommonName   string
	Organization []string
	DNSNames     []string
	IPs          []string
	URIs         []string

	Key string

	Issuer Issuer
	CRSpec CRSpec
}

type Sign struct {
	CSR string

	Issuer Issuer
	CRSpec CRSpec
}

type Object struct {
	Name string
}

type Issuer struct {
	IssuerName  string
	IssuerKind  string
	IssuerGroup string
}

type CRSpec struct {
	Duration string
	IsCA     bool
	Out      string
}