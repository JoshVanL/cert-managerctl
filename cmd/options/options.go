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
	Object Object
}

type Sign struct {
	CSR string

	Issuer Issuer
	CRSpec CRSpec
	Object Object
}

type Object struct {
	Name      string
	Namespace string
}

type Issuer struct {
	Name  string
	Kind  string
	Group string
}

type CRSpec struct {
	Duration string
	IsCA     bool
	Out      string
}
