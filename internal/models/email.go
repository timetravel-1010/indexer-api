package models

import "net/mail"

// An Email contains all the information of an e-mail.
type Email struct {
	MessageID               string          `json:"Message-Id"`
	Date                    string          `json:"Date"`
	From                    string          `json:"From"`
	To                      []*mail.Address `json:"To"`
	CC                      []*mail.Address `json:"Cc"`
	BCC                     []*mail.Address `json:"Bcc"`
	Subject                 string          `json:"Subject"`
	MimeVersion             string          `json:"Mime-Version"`
	ContentType             string          `json:"Content-Type"`
	ContentTransferEncoding string          `json:"Content-Transfer-Encoding"`
	XFrom                   string          `json:"X-From"`
	XTo                     []*mail.Address `json:"X-To"`
	Xcc                     []*mail.Address `json:"X-Cc"`
	Xbcc                    []*mail.Address `json:"X-Bcc"`
	XFolder                 string          `json:"X-Folder"`
	XOrigin                 string          `json:"X-Origin"`
	XFileName               string          `json:"X-Filename"`
	Body                    string          `json:"Body"`
}
