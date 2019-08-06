package options

type Flags struct {
	Kubeconfig string
	LogLevel   string

	Request Request
	Get     Get
}

type Request struct {
	Cert ReqCert
	Sign ReqSign
}

type ReqCert struct {
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

type ReqSign struct {
	CSR string

	Issuer Issuer
	CRSpec CRSpec
	Object Object
}

type Get struct {
	Cert GetCert
}

type GetCert struct {
	Object
	Out  string
	Wait bool
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

type CROptions struct {
	Issuer Issuer
	CRSpec CRSpec
	Object Object
}
