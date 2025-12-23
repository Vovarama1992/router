package wg

import (
	"bytes"
	"text/template"
)

type ClientConfigParams struct {
	ClientPrivateKey string
	ClientAddress    string
	DNS              string

	ServerPublicKey string
	ServerEndpoint  string
	AllowedIPs      string
}

func RenderClientConfig(tplPath string, params ClientConfigParams) (string, error) {
	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}
