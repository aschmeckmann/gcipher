package api

type Response struct {
	Success bool        `json:"success"`
	Errors  []Error     `json:"errors,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Request struct {
	Data RequestData `json:"data"`
	Auth Auth        `json:"auth"`
}

type RequestData struct {
	Applicant    string `json:"applicant,omitempty"`
	CSR          string `json:"csr,omitempty"`
	Lifetime     int    `json:"lifetime,omitempty"`
	Type         string `json:"type,omitempty"`
	State        string `json:"state,omitempty"`
	SerialNumber string `json:"serialnumber,omitempty"`
}

type CertificateResponseData struct {
	CertificatePEM string `json:"cert"`
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
