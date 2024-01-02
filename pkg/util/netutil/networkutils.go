package netutil

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// GetLocalIP retrieves the first non-loopback IP address found on the local machine.
func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", errors.New("no active network interface found")
}

// CheckHostAvailability checks if a host is available over TCP within a timeout duration.
func CheckHostAvailability(host string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

// SimpleGetRequest performs a GET request to the specified URL and returns the status code and body.
func SimpleGetRequest(url string) (int, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return 0, "", err
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	body := string(bodyBytes)

	return resp.StatusCode, body, nil
}

// ListenOnPort starts listening on the specified port and returns the listener.
func ListenOnPort(port string) (net.Listener, error) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}
	return listener, nil
}

// GetMACAddress retrieves the MAC address of the first interface found.
func GetMACAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, interf := range interfaces {
		if mac := interf.HardwareAddr; mac != nil {
			return mac.String(), nil
		}
	}

	return "", errors.New("no MAC address found")
}

// GetPublicIP makes an external request to identify the public IP of the machine.
func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(ip)), nil
}

// IsPortOpen checks if a port is open on a specific host.
func IsPortOpen(host string, port string, timeout time.Duration) bool {
	address := net.JoinHostPort(host, port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return false
	}

	conn.Close()

	return true
}

// DownloadFile will download a url to a local file. It's efficient because it will write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)

	return err
}

// GetSSLExpiry returns the expiry date of the SSL certificate of the given domain.
func GetSSLExpiry(domain string) (time.Time, error) {
	conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return time.Time{}, err
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	return cert.NotAfter, nil
}

// NetworkInterfaces lists all network interfaces along with their details.
func NetworkInterfaces() ([]net.Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	return interfaces, nil
}

// GetHostname returns the hostname of the local machine.
func GetHostname() (string, error) {
	return os.Hostname()
}

// ResolveHostname resolves the hostname to an IP address.
func ResolveHostname(hostname string) (string, error) {
	ips, err := net.LookupIP(hostname)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			return ipv4.String(), nil
		}
	}

	return "", errors.New("no IPv4 address found")
}
