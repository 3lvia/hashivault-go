# hashivault-go
Provides a Vault client for the Hashicorp Vault secrets management solution. See the golang documentation for more 
information (specifically ./pkg/hashivault/doc.go).

## Generated documentation
The source code in this repository is documented using godoc standards (https://tip.golang.org/doc/comment). The
documentation can be generated and viewed on the local development machine as follows:
0. Install godoc: `go install -v golang.org/x/tools/cmd/godoc@latest` (if not already installed)
1. Run `godoc -http=:6060` in the root of the project
2. Open a browser and navigate to `http://localhost:6060/pkg/github.com/3lvia/hashivault-go/`