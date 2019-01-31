package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const version = "0.0.2"

var (
	// Ref is the buildtime code ref
	Ref string
	// Sha is the buildtime code sha
	Sha string
)

func appVersion() string {
	ver := []string{version, Ref, Sha}
	return fmt.Sprintf(`%s (%s on %s/%s; %s)`, strings.Join(ver, "/"), runtime.Version(), runtime.GOOS, runtime.GOARCH, runtime.Compiler)
}

func usage() string {
	t := template.Must(template.New("usage").Parse(`
Usage: {{ .Self }} user-spec command [args]
   eg: {{ .Self }} myuser bash
       {{ .Self }} nobody:root bash -c 'whoami && id'
       {{ .Self }} 1000:1 id

  Environment Variables:

    VEST_USER=user[:group]
      The user [and group] to run the command under. Overrides commandline if set.

    VEST_PROVIDERS=provider1,...
      Comma separated list of enabled providers. By default only vault is enabled.

    SOPS_FILES=/path/to/file[;/path/to/output[;mode]]:...
      If SOPS_FILES is set, will iterate over each file (colon separated), attempting to decrypt with Sops.
      The decrypted cleartext file can be optionally written out to a separate location (with optional filemode)
      or will be parsed into a map[string]string and injected into Environ
    
    VAULT_KEYS=/path/to/key[@version]:...
      If VAULT_KEYS is set, will iterate over each key (colon separated), attempting to get the secret from Vault.
      Secrets are pulled at the optional version or latest, then injected into Environ. If running in Kubernetes,
      the Pod's ServiceAccount token will automatically be looked up and used for Vault authentication.
    
    VAULT_*
      All vault client configuration environment variables are respected.
      More information at https://www.vaultproject.io/docs/commands/#environment-variables
    
    EJSON_FILES=/path/to/file1:...
    EJSON_KEYS=pubkey;privkey:...
      If EJSON_FILES is set, will iterate over each file (colon separated), attempting to decrypt using keys
      from EJSON_KEYS. If EJSON_FILES is not set, will look for any .ejson files in CWD. Cleartext decrypted
      json will be parsed into a map[string]string and injected into Environ.
    
    DOTENV_FILES=/path/to/file1:...
      if DOTENV_FILES is set, will iterate over each file, parse and inject into Environ. If DOTENV_FILES is
      not set, will look for any .env files in CWD.

{{ .Self }} version: {{ .Version }}
{{ .Self }} license: GPL-3 (full text at https://github.com/lumoslabs/vestibule)
`))
	var b bytes.Buffer
	template.Must(t, t.Execute(&b, struct {
		Self    string
		Version string
	}{
		Self:    filepath.Base(os.Args[0]),
		Version: appVersion(),
	}))
	return strings.TrimSpace(b.String()) + "\n"
}
