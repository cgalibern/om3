package arraypure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/devans10/pugo/pure1"
	"github.com/golang-jwt/jwt"
	"github.com/spf13/cobra"

	"github.com/opensvc/om3/core/array"
	"github.com/opensvc/om3/core/driver"
	"github.com/opensvc/om3/core/object"
)

var (
	WWIDPrefix      = "624a9370"
	RenewStatus     = 403
	ItemsPerPage    = 100
	MaxPages        = 1000
	DayMilliseconds = 24 * 60 * 60 * 1000
	RequestTimeout  = 10
	Head            = "/api/2.8"
)

type (
	Array struct {
		*array.Array
		token *pure1Token
	}

	pure1Response struct {
		TotalItems        int   `json:"total_item_count,omitempty"`
		ContinuationToken any   `json:"continuation_token,omitempty"`
		Items             []any `json:"items,omitempty"`
	}

	pure1Token struct {
		AccessToken     string `json:"access_token,omitempty"`
		IssuedTokenType string `json:"issued_token_type,omitempty"`
		TokenType       string `json:"token_type,omitempty"`
		ExpiresIn       int    `json:"expires_in,omitempty"`
	}
)

func init() {
	driver.Register(driver.NewID(driver.GroupArray, "pure"), NewDriver)
}

func NewDriver() array.Driver {
	t := New()
	var i any = t
	return i.(array.Driver)
}

func New() *Array {
	t := &Array{
		Array: array.New(),
	}
	return t
}

func (t *Array) Run(args []string) error {
	newParent := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:           "array",
			Short:         "Manage a purestorage storage array",
			SilenceUsage:  true,
			SilenceErrors: true,
		}
		return cmd
	}

	newMapCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "map",
			Short: "map commands",
		}
		return cmd
	}
	newUnmapCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "unmap",
			Short: "unmap commands",
		}
		return cmd
	}
	newAddCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "add",
			Short: "add commands",
		}
		return cmd
	}
	newDelCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "del",
			Short: "del commands",
		}
		return cmd
	}

	newUnmapDiskCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "disk",
			Short: "unmap a volume",
			RunE: func(_ *cobra.Command, _ []string) error {
				return fmt.Errorf("TODO")
			},
		}
		useFlagID(cmd)
		useFlagName(cmd)
		useFlagMapping(cmd)
		useFlagHostGroup(cmd)
		useFlagSerial(cmd)
		return cmd
	}
	newMapDiskCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "disk",
			Short: "map a volume",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if data, err := t.mapDisk(name, mappings, lun); err != nil {
					return err
				} else {
					return dump(data)
				}
			},
		}
		useFlagID(cmd)
		useFlagName(cmd)
		useFlagMapping(cmd)
		useFlagLUN(cmd)
		useFlagHost(cmd)
		useFlagHostGroup(cmd)
		useFlagSerial(cmd)
		return cmd
	}
	newDelDiskCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "disk",
			Short: "unmap a volume and delete",
			RunE: func(_ *cobra.Command, _ []string) error {
				if data, err := t.delDisk(name); err != nil {
					return err
				} else {
					return dump(data)
				}
			},
		}
		useFlagName(cmd)
		useFlagNow(cmd)
		useFlagID(cmd)
		useFlagSerial(cmd)
		return cmd
	}
	newAddDiskCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "disk",
			Short: "add a zvol-type dataset and map",
			RunE: func(cmd *cobra.Command, _ []string) error {
				if data, err := t.addDisk(name, size, mappings); err != nil {
					return err
				} else {
					return dump(data)
				}
			},
		}
		useFlagName(cmd)
		useFlagSize(cmd)
		useFlagMapping(cmd)
		return cmd
	}
	newGetCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "get",
			Short: "get commands",
		}
		return cmd
	}
	newGetHostsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "hosts",
			Short: "get hosts",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getHosts(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetConnectionsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "connections",
			Short: "get connections",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getConnections(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetVolumesCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "volumes",
			Short: "get volumes",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getVolumes(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetControllersCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "controllers",
			Short: "get controllers",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getControllers(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetDrivesCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "drives",
			Short: "get drives",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getDrives(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetPodsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "pods",
			Short: "get pods",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getPods(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetPortsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "ports",
			Short: "get ports",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getPorts(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetNetworkInterfacesCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "interfaces",
			Short: "get network interfaces",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getNetworkInterfaces(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetVolumeGroupsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "volumegroups",
			Short: "get volumegroups",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getVolumeGroups(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetHostGroupsCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "hostgroups",
			Short: "get hostgroups",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getHostGroups(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetArraysCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "arrays",
			Short: "get arrays",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getArrays(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}
	newGetHardwareCmd := func() *cobra.Command {
		cmd := &cobra.Command{
			Use:   "hardware",
			Short: "get hardware",
			RunE: func(_ *cobra.Command, _ []string) error {
				data, err := t.getHardware(filter)
				if err != nil {
					return err
				}
				return dump(data)
			},
		}
		useFlagFilter(cmd)
		return cmd
	}

	parent := newParent()

	// skip past the --array <array> arguments
	parent.SetArgs(os.Args[4:])

	addCmd := newAddCmd()
	addCmd.AddCommand(newAddDiskCmd())
	parent.AddCommand(addCmd)

	delCmd := newDelCmd()
	delCmd.AddCommand(newDelDiskCmd())
	parent.AddCommand(delCmd)

	getCmd := newGetCmd()
	getCmd.AddCommand(newGetHostsCmd())
	getCmd.AddCommand(newGetConnectionsCmd())
	getCmd.AddCommand(newGetControllersCmd())
	getCmd.AddCommand(newGetDrivesCmd())
	getCmd.AddCommand(newGetHardwareCmd())
	getCmd.AddCommand(newGetArraysCmd())
	getCmd.AddCommand(newGetHostGroupsCmd())
	getCmd.AddCommand(newGetVolumeGroupsCmd())
	getCmd.AddCommand(newGetNetworkInterfacesCmd())
	getCmd.AddCommand(newGetPortsCmd())
	getCmd.AddCommand(newGetPodsCmd())
	getCmd.AddCommand(newGetVolumesCmd())
	parent.AddCommand(getCmd)

	mapCmd := newMapCmd()
	mapCmd.AddCommand(newMapDiskCmd())
	parent.AddCommand(mapCmd)

	unmapCmd := newUnmapCmd()
	unmapCmd.AddCommand(newUnmapDiskCmd())
	parent.AddCommand(unmapCmd)

	return parent.Execute()
}

func (t Array) api() string {
	return t.Config().GetString(t.Key("api"))
}

func (t Array) clientID() string {
	return t.Config().GetString(t.Key("client_id"))
}

func (t Array) keyID() string {
	return t.Config().GetString(t.Key("key_id"))
}

func (t Array) username() string {
	return t.Config().GetString(t.Key("username"))
}

func (t Array) issuer() string {
	return t.Config().GetString(t.Key("issuer"))
}

func (t Array) insecure() bool {
	return t.Config().GetBool(t.Key("insecure"))
}

func (t Array) secret() string {
	return t.Config().GetString(t.Key("secret"))
}

func (t *Array) sec() (object.Sec, error) {
	s, err := t.Config().GetStringStrict(t.Key("secret"))
	if err != nil {
		return nil, err
	}
	return object.NewSec(s, object.WithVolatile(true))
}

func (t *Array) privateKey() ([]byte, error) {
	sec, err := t.sec()
	if err != nil {
		return nil, err
	}
	return sec.DecodeKey("private_key")
}

func (t *Array) getToken() (*pure1Token, error) {
	if t.token != nil {
		return t.token, nil
	}
	if err := t.newToken(); err != nil {
		return nil, err
	}
	return t.token, nil
}

func (t *Array) newToken() error {
	username := t.username()
	issuer := t.issuer()
	now := time.Now().Unix()
	if issuer == "" {
		issuer = username
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"aud": t.clientID(),
		"sub": username,
		"iss": issuer,
		"iat": now,
		"exp": now + int64(DayMilliseconds),
	})
	token.Header["kid"] = t.keyID()

	privateKey, err := t.privateKey()
	if err != nil {
		return err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return err
	}

	// Generate encoded token and send it as response.
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return err
	}

	values := url.Values{}
	values.Add("content-type", "application/json")
	values.Add("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	values.Add("subject_token", tokenStr)
	values.Add("subject_token_type", "urn:ietf:params:oauth:token-type:jwt")

	pure1URL := fmt.Sprintf("%s/oauth2/%s/token", t.api(), "1.0")
	req, err := http.NewRequest(http.MethodPost, pure1URL, strings.NewReader(values.Encode()))
	if err != nil {
		return err
	}
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	pToken := &pure1Token{}
	if err := json.Unmarshal(responseBody, pToken); err != nil {
		return err
	}

	t.token = pToken
	return nil
}

func (t *Array) client() (*pure1.Client, error) {
	restVersion := ""
	clientID := t.clientID()
	privateKey, err := t.privateKey()
	if err != nil {
		return nil, err
	}
	pureCli, err := pure1.NewClient(clientID, privateKey, restVersion)
	if err != nil {
		return nil, err
	}
	return pureCli, nil
}

func (c *Array) Do(req *http.Request, v interface{}, reestablishSession bool) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %w", err)
	}
	defer resp.Body.Close()

	if err := validateResponse(resp); err != nil {
		return nil, fmt.Errorf("validate response: %w", err)
		//return resp, err
	}

	err = decodeResponse(resp, v)
	if err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return resp, nil
}

func dump(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	return enc.Encode(data)
}

func (t *Array) addDisk(name, size string, mappings []string) (any, error) {
	return nil, nil
}

func (t *Array) mapDisk(name string, mappings []string, lun int) (any, error) {
	return nil, nil
}

func (t *Array) delDisk(name string) (any, error) {
	return nil, nil
}

func (t *Array) getHosts(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/hosts", params, nil)
}

func (t *Array) getConnections(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/connections", params, nil)
}

func (t *Array) getVolumes(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/volumes", params, nil)
}

func (t *Array) getVolumeGroups(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/volume-groups", params, nil)
}

func (t *Array) getControllers(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/controllers", params, nil)
}

func (t *Array) getDrives(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/drives", params, nil)
}

func (t *Array) getPods(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/pods", params, nil)
}

func (t *Array) getPorts(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/ports", params, nil)
}

func (t *Array) getNetworkInterfaces(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/network-interfaces", params, nil)
}

func (t *Array) getHostGroups(filter string) (any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/host-groups", params, nil)
}

func getParams(filter string) map[string]string {
	params := map[string]string{"total_item_count": "true", "limit": fmt.Sprint(ItemsPerPage)}
	if filter != "" {
		params["filter"] = filter
	}
	return params
}

func (t *Array) getHardware(filter string) ([]any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/hardware", params, nil)
}

func (t *Array) getArrays(filter string) ([]any, error) {
	params := getParams(filter)
	return t.doGet("GET", "/arrays", params, nil)
}

func (t *Array) doGet(method string, path string, params map[string]string, data interface{}) ([]any, error) {
	req, err := t.newRequest(method, path, params, data)
	if err != nil {
		return nil, err
	}
	var r pure1Response
	items := make([]any, 0)
	_, err = t.Do(req, &r, true)
	if err != nil {
		return nil, err
	}
	for len(items) < r.TotalItems {
		for _, item := range r.Items {
			//i := PureArray{}
			//s, _ := json.Marshal(item)
			//json.Unmarshal([]byte(s), &i)
			items = append(items, item)
		}

		if len(items) < r.TotalItems {
			if r.ContinuationToken != nil {
				if params == nil {
					params = map[string]string{"continuation_token": r.ContinuationToken.(string)}
				} else {
					params["continuation_token"] = r.ContinuationToken.(string)
				}
				req, err := t.newRequest(method, path, params, data)
				if err != nil {
					return nil, err
				}

				_, err = t.Do(req, r, false)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return items, nil
}

func (t *Array) newRequest(method string, path string, params map[string]string, data interface{}) (*http.Request, error) {
	fpath := t.api() + Head + path
	baseURL, err := url.Parse(fpath)
	if err != nil {
		return nil, err
	}
	if params != nil {
		ps := url.Values{}
		for k, v := range params {
			ps.Set(k, v)
		}
		baseURL.RawQuery = ps.Encode()
	}
	req, err := http.NewRequest(method, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	if data != nil {
		jsonString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, baseURL.String(), bytes.NewBuffer(jsonString))
		if err != nil {
			return nil, err
		}
	}

	token, err := t.getToken()
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", `Bearer `+token.AccessToken)

	return req, err
}

// decodeResponse function reads the http response body into an interface.
func decodeResponse(r *http.Response, v interface{}) error {
	if v == nil {
		return fmt.Errorf("nil interface provided to decodeResponse")
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	err := json.Unmarshal([]byte(bodyString), &v)

	return err
}

// validateResponse checks that the http response is within the 200 range.
// Some functionality needs to be added here to check for some specific errors,
// and probably add the equivlents to PureError and PureHTTPError from the Python
// REST client.
func validateResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)
	bodyString := string(bodyBytes)
	return fmt.Errorf("Response code: %d, Response body: %s", r.StatusCode, bodyString)
}
