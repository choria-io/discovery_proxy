package choria

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/go-openapi/strfmt"

	log "github.com/Sirupsen/logrus"
	apiclient "github.com/choria-io/discovery_proxy/client"
	httptransport "github.com/go-openapi/runtime/client"
)

// Choria is a utilty encompasing mcollective and choria config and various utilities
type Choria struct {
	Config *MCollectiveConfig
	Sets   *Sets
}

// Server is a representation of a network server host and port
type Server struct {
	Host string
	Port int
}

// New sets up a Choria with all its config loaded and so forth
func New(path string) (*Choria, error) {
	// TODO check SSL sanity

	c := Choria{}

	config, err := NewConfig(path)
	if err != nil {
		return &c, err
	}

	c.Config = config

	return &c, nil
}

// TrySrvLookup will attempt to lookup a series of names returning the first found
// if SRV lookups are disabled or nothing is found the default will be returned
func (c *Choria) TrySrvLookup(names []string, defaultSrv Server) (Server, error) {
	if !c.Config.Choria.UseSRVRecords {
		return defaultSrv, nil
	}

	for _, q := range names {
		a, err := c.QuerySrvRecords([]string{q})
		if err == nil {
			log.Infof("Found %s:%d from %s SRV lookups", a[0].Host, a[0].Port, strings.Join(names, ", "))

			return a[0], nil
		}
	}

	log.Debugf("Could not find SRV records for %s, returning defaults %s:%d", strings.Join(names, ", "), defaultSrv.Host, defaultSrv.Port)

	return defaultSrv, nil
}

// QuerySrvRecords looks for SRV records within the right domain either
// thanks to facter domain or the configured domain.
//
// If the config disables SRV then a error is returned.
func (c *Choria) QuerySrvRecords(records []string) ([]Server, error) {
	servers := []Server{}

	if !c.Config.Choria.UseSRVRecords {
		return servers, errors.New("SRV lookups are disabled in the configuration file")
	}

	domain, err := c.FacterDomain()
	if err != nil {
		return servers, err
	}

	for _, q := range records {
		record := q + "." + domain
		log.Debugf("Attempting SRV lookup for %s", record)

		cname, addrs, err := net.LookupSRV(record, "", "")
		if err != nil {
			return servers, err
		}

		log.Debugf("Found %d SRV records for %s", len(addrs), cname)

		for _, a := range addrs {
			servers = append(servers, Server{Host: a.Target, Port: int(a.Port)})
		}
	}

	return servers, nil
}

// DiscoveryServer is the server configured as a discovery proxy
func (c *Choria) DiscoveryServer() (Server, error) {
	s := Server{
		Host: c.Config.Choria.DiscoveryHost,
		Port: c.Config.Choria.DiscoveryPort,
	}

	if !c.ProxiedDiscovery() {
		return s, errors.New("Proxy discovery is not enabled")
	}

	result, err := c.TrySrvLookup([]string{"_mcollective-discovery._tcp"}, s)

	return result, err
}

// ProxiedDiscovery determines if a client is configured for proxied discover
func (c *Choria) ProxiedDiscovery() bool {
	if c.Config.HasOption("plugin.choria.discovery_host") || c.Config.HasOption("plugin.choria.discovery_port") {
		return true
	}

	return c.Config.Choria.DiscoveryProxy
}

// Certname determines the choria certname
func (c *Choria) Certname() string {
	certname := c.Config.Identity

	currentUser, _ := user.Current()

	if currentUser.Uid != "0" {
		if u, ok := os.LookupEnv("USER"); ok {
			certname = fmt.Sprintf("%s.mcollective", u)
		}
	}

	if u, ok := os.LookupEnv("MCOLLECTIVE_CERTNAME"); ok {
		certname = u
	}

	return certname
}

// CAPath determines the path to the CA file
func (c *Choria) CAPath() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "certs", "ca.pem"), nil
}

// ClientPrivateKey determines the location to the client cert
func (c *Choria) ClientPrivateKey() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "private_keys", fmt.Sprintf("%s.pem", c.Certname())), nil
}

// ClientPublicCert determines the location to the client cert
func (c *Choria) ClientPublicCert() (string, error) {
	ssl, err := c.SSLDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(ssl, "certs", fmt.Sprintf("%s.pem", c.Certname())), nil
}

// SSLDir determines the AIO SSL directory
func (c *Choria) SSLDir() (string, error) {
	if c.Config.Choria.SSLDir != "" {
		return c.Config.Choria.SSLDir, nil
	}

	u, _ := user.Current()
	if u.Uid == "0" {
		path, err := c.PuppetSetting("ssldir")
		if err != nil {
			return "", err
		}

		return path, nil
	}

	return filepath.Join(u.HomeDir, ".puppetlabs", "etc", "puppet", "ssl"), nil
}

// PuppetSetting retrieves a config setting by shelling out to puppet apply --configprint
func (c *Choria) PuppetSetting(setting string) (string, error) {
	args := []string{"apply", "--configprint", setting}

	out, err := exec.Command("puppet", args...).Output()
	if err != nil {
		return "", err
	}

	return strings.Replace(string(out), "\n", "", -1), nil
}

// FacterDomain determines the machines domain by querying facter. Returns "" when unknown
func (c *Choria) FacterDomain() (string, error) {
	cmd := c.FacterCmd()

	if cmd == "" {
		return "", errors.New("Could ont find your facter command")
	}

	out, err := exec.Command(cmd, "networking.domain").Output()
	if err != nil {
		return "", errors.New("Could not resolve the server domain via facter: " + err.Error())
	}

	return strings.Replace(string(out), "\n", "", -1), nil
}

// FacterCmd finds the path to facter using first AIO path then a `which` like command
// TODO: windows support
func (c *Choria) FacterCmd() string {
	if _, err := os.Stat("/opt/puppetlabs/bin/facter"); err == nil {
		return "/opt/puppetlabs/bin/facter"
	}

	path, err := exec.LookPath("facter")
	if err != nil {
		return ""
	}

	return path
}

// DiscoveryProxyClient is a client for the discovery REST service
func (c *Choria) DiscoveryProxyClient() (*apiclient.DiscoveryProxy, error) {
	server, err := c.DiscoveryServer()
	if err != nil {
		return nil, err
	}

	host := fmt.Sprintf("%s:%d", server.Host, server.Port)

	context, err := c.SSLContext()
	if err != nil {
		return nil, err
	}

	http := &http.Client{Transport: context}
	transport := httptransport.NewWithClient(host, apiclient.DefaultBasePath, apiclient.DefaultSchemes, http)

	return apiclient.New(transport, strfmt.NewFormats()), nil
}

// SSLContext creates a SSL context loaded with our certs and ca
func (c *Choria) SSLContext() (*http.Transport, error) {
	pub, _ := c.ClientPublicCert()
	pri, _ := c.ClientPrivateKey()
	ca, _ := c.CAPath()

	cert, err := tls.LoadX509KeyPair(pub, pri)

	if err != nil {
		return &http.Transport{}, errors.New("Could not load certificate " + pub + " and key " + pri + ": " + err.Error())
	}

	caCert, err := ioutil.ReadFile(ca)

	if err != nil {
		return &http.Transport{}, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}

	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}

	return transport, nil
}
