package interserviceclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gopkg.in/yaml.v2"
)

// Inter service token expire minutes. Specify after how long a token will expire
const (
	ISCExpireEnvVarName = "INTER_SERVICE_TOKEN_EXPIRE_MINUTES"
)

// ISCService defines the blueprint of a dependency service. This struct is here to maintain
// uniform structure definitions
type ISCService struct {
	// The name of the service that is been depended upon e.g mailgun, mpesa
	Name string

	// The endpoint where the service serves requests. The dependant should know forehand where to
	// this services lives
	RootDomain string
}

// GetJWTKey returns a byte slice of the JWT secret key
func GetJWTKey() []byte {
	key := serverutils.MustGetEnvVar(JWTSecretKey)
	return []byte(key)
}

// Claims a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide public claims
// Provides way for adding private claims
type Claims struct {
	jwt.StandardClaims
}

// InterServiceClient defines a client for use in interservice communication
type InterServiceClient struct {
	Name              string
	RequestRootDomain string
	httpClient        http.Client
}

// NewInterserviceClient initializes a new interservice client
func NewInterserviceClient(s ISCService) (*InterServiceClient, error) {
	return &InterServiceClient{
		Name:              s.Name,
		RequestRootDomain: s.RootDomain,
		httpClient: http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
			Timeout:   time.Duration(1 * time.Minute),
		},
	}, nil
}

// CreateAuthToken returns a signed JWT for use in authentication.
func (c InterServiceClient) CreateAuthToken(ctx context.Context) (string, error) {
	var expireMinutes int
	expireMinutesStr, err := serverutils.GetEnvVar(ISCExpireEnvVarName)
	if err != nil {
		// Fallback for when the env var is not set
		expireMinutesStr = "60"
	}
	expireMinutes, err = strconv.Atoi(expireMinutesStr)
	if err != nil {
		return "", fmt.Errorf("misconfigured ENV: %w", err)
	}
	claims := &Claims{
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Duration(expireMinutes) * time.Minute).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJWTKey())
	if err != nil {
		return "", fmt.Errorf("failed to create token with err: %v", err)
	}

	return tokenString, nil
}

// GenerateRequestURL generate a url with path for requested resource.
func (c InterServiceClient) generateRequestURL(path string) string {
	return fmt.Sprintf("%v/%v", c.RequestRootDomain, path)
}

// MakeRequest performs an inter service http request and returns a response
func (c InterServiceClient) MakeRequest(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {

	url := c.generateRequestURL(path)

	token, tknErr := c.CreateAuthToken(ctx)
	if tknErr != nil {
		return nil, tknErr
	}

	// A GET request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if method == http.MethodGet {
		req, reqErr := http.NewRequestWithContext(ctx, method, url, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		return c.httpClient.Do(req)
	}

	encoded, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBuffer(encoded)
	req, reqErr := http.NewRequestWithContext(ctx, method, url, payload)
	if reqErr != nil {
		return nil, reqErr
	}

	if serverutils.IsDebug() {
		r, _ := httputil.DumpRequest(req, true)
		log.Println(string(r))
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return c.httpClient.Do(req)
}

// jwtCheckFn is a function type for authorization and authentication checks
// there can be several e.g an authentication check runs first then an authorization
// check runs next if the authentication passes etc
type jwtCheckFn = func(r *http.Request) (bool, map[string]string, *jwt.Token)

// InterServiceAuthenticationMiddleware handles jwt authentication
func InterServiceAuthenticationMiddleware() func(http.Handler) http.Handler {
	// multiple checks can be run in sequence
	jwtCheckFuncs := []jwtCheckFn{HasValidJWTBearerToken}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				errs := []map[string]string{}

				for _, checkFunc := range jwtCheckFuncs {
					shouldContinue, errMap, _ := checkFunc(r)
					if shouldContinue {

						next.ServeHTTP(w, r)
						return
					}
					errs = append(errs, errMap)
				}

				serverutils.WriteJSONResponse(w, errs, http.StatusUnauthorized)
			})
	}
}

// HasValidJWTBearerToken returns true with no errors if the request has a valid bearer token in the authorization header.
// Otherwise, it returns false and the error in a map with the key "error"
func HasValidJWTBearerToken(r *http.Request) (bool, map[string]string, *jwt.Token) {
	bearerToken, err := firebasetools.ExtractBearerToken(r)
	if err != nil {

		return false, serverutils.ErrorMap(err), nil
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(bearerToken, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJWTKey(), nil
	})

	if err != nil {
		return false, serverutils.ErrorMap(err), nil
	}

	return true, nil, token
}

// Dep is the dependency definition
type Dep struct {
	DepName       string `yaml:"depName"`
	DepRootDomain string `yaml:"depRootDomain"`
}

//DepsConfig is the config for dependencies of a particular service
type DepsConfig struct {
	Staging    []Dep `yaml:"staging"`
	Testing    []Dep `yaml:"testing"`
	Demo       []Dep `yaml:"demo"`
	Production []Dep `yaml:"production"`
	E2E        []Dep `yaml:"e2e"`
}

// PathToDepsFile return the path to deps.yaml file
func PathToDepsFile() string {
	cwd, _ := os.Getwd()
	return getDepsPath(filepath.Join(cwd, DepsFileName))
}

// recursively get the path to the deps.yaml file
func getDepsPath(path string) string {
	_, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		n := strings.Split(filepath.Dir(path), "/")
		m := n[:len(n)-1]
		p := filepath.Join(strings.Join(m, "/"), DepsFileName)
		return getDepsPath(p)
	}
	return path
}

// GetDepFromConfig retrives a specific config from config slice
func GetDepFromConfig(name string, config []Dep) *Dep {
	var d Dep
	for _, dep := range config {
		if dep.DepName == name {
			d = dep
		}
	}
	return &d
}

// SetupISCclient returns an InterServiceClient
func SetupISCclient(config DepsConfig, serviceName string) (*InterServiceClient, error) {
	if serverutils.GetRunningEnvironment() == serverutils.StagingEnv {
		dep := GetDepFromConfig(serviceName, config.Staging)
		client, err := NewInterserviceClient(ISCService{Name: dep.DepName, RootDomain: dep.DepRootDomain})
		return client, err
	}

	if serverutils.GetRunningEnvironment() == serverutils.TestingEnv {
		dep := GetDepFromConfig(serviceName, config.Testing)
		client, err := NewInterserviceClient(ISCService{Name: dep.DepName, RootDomain: dep.DepRootDomain})
		return client, err
	}

	if serverutils.GetRunningEnvironment() == serverutils.DemoEnv {
		dep := GetDepFromConfig(serviceName, config.Demo)
		client, err := NewInterserviceClient(ISCService{Name: dep.DepName, RootDomain: dep.DepRootDomain})
		return client, err
	}

	if serverutils.GetRunningEnvironment() == serverutils.ProdEnv {
		dep := GetDepFromConfig(serviceName, config.Production)
		client, err := NewInterserviceClient(ISCService{Name: dep.DepName, RootDomain: dep.DepRootDomain})
		return client, err
	}

	if serverutils.GetRunningEnvironment() == E2eEnv {
		dep := GetDepFromConfig(serviceName, config.E2E)
		client, err := NewInterserviceClient(ISCService{Name: dep.DepName, RootDomain: dep.DepRootDomain})
		return client, err
	}

	return nil, fmt.Errorf("failed to setup isc client")
}

// LoadDepsFromYAML loads the interservice dependency config from a deps.yaml
// file that is at the default location
func LoadDepsFromYAML() (*DepsConfig, error) {
	var config DepsConfig

	file, err := ioutil.ReadFile(filepath.Clean(PathToDepsFile()))
	if err != nil {
		return nil, fmt.Errorf("can't read deps file: %w", err)
	}

	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, fmt.Errorf("can't unmarshal deps YAML: %w", err)
	}

	return &config, nil
}
