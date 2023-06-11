// Package hashivault provides a Vault client for the Hashicorp Vault secrets management solution. It supports three
// modes of authentication against Vault:
// 1. GitHub authentication for people
// 2. Kubernetes authentication for pods
// 3. Azure AD SSO authentication for people (not yet implemented)
//
// The mode of authentication is determined by the environment variables:
//  1. GITHUB_TOKEN. If this variable is set, the client will authenticate against Vault using the GitHub auth method.
//     This takes precedence over the other methods.
//  2. MOUNT_PATH and ROLE. If these variables are set, the client will authenticate using the Kubernetes auth method.
//
// The environment variable VAULT_ADDR must be set to the address of the Vault server.
//
// The client will periodically renew the authentication token. The token is renewed when it has less than 30 seconds
// left to live. The token is renewed in a separate goroutine, so the client will not block while waiting for the token
// to be renewed.
//
// The main abstraction of this package is the SecretsManager interface. A new instance of SecretsManager is created
// with the New function. The returned SecretsManager is safe to use concurrently. Normally, only a single instance of
// SecretsManager should be created in an application, presumably in the main function as the application is being
// initialised.
//
// The secrets that are returned by the SecretsManager are represented as functions that return a map[string]any. The
// point is that the returned function will always return the latest version of the secret. Therefore, clients should
// save a reference to the function rather than saving the actual secrets, and invoke the func just-in-time as the
// secret is needed. The returned function is safe to use concurrently.
//
// The SecretsManager interface also provides a method for setting the default Google credentials for the current
// process.
//
// The token refresh functionality runs in a separate goroutine, and also a new goroutine will be started for each
// fetched secret that is renewable and has a lease duration. In order to communicate errors from these goroutines, the
// New function returns a channel of errors in addition to the SecretsManager. Clients should start a goroutine that
// reads from this channel and handles errors as appropriate.
//
// The following example shows how to use the SecretsManager:
// ```
// import (
//
//	"github.com/3lvia/hashivault-go/pkg/hashivault"
//	"log"
//
// )
//
//	func main() {
//	 v, errChan, err := hashivault.New()
//   if err != nil {
//	   log.Fatal(err)
//	 }

//	 go func(ec <-chan error) {
//	   for err := range ec {
//	     log.Println(err)
//	   }
//	 }(errChan)
//
//	 secret, err := v.GetSecret("kunde/kv/data/appinsights/kunde")
//	 if err != nil {
//	   log.Fatal(err)
//	 }
//
//	 mapOfSecrets := secret()
//	 _ = mapOfSecrets
//	}
//
// ```
package hashivault
