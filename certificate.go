package wenex

import (
	"crypto/tls"
	"errors"

	"golang.org/x/crypto/acme/autocert"
)

func stringCert(wnx *Wenex, path string) (*tls.Config, error) {
	if wnx.Config.Get(path) == nil {
		return nil, nil
	}

	certPem, err := wnx.Config.String(path + ".cert")
	if err != nil {
		return nil, err
	}

	keyPem, err := wnx.Config.String(path + ".key")
	if err != nil {
		return nil, err
	}

	cert, err := tls.X509KeyPair([]byte(certPem), []byte(keyPem))
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func loadCert(wnx *Wenex, path string) (*tls.Config, error) {
	if wnx.Config.Get(path) == nil {
		return nil, nil
	}

	certFile, err := wnx.Config.String(path + ".cert")
	if err != nil {
		return nil, err
	}

	keyFile, err := wnx.Config.String(path + ".key")
	if err != nil {
		return nil, err
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}

func autoCert(wnx *Wenex, path string) (*autocert.Manager, error) {
	if wnx.Config.Get(path) == nil {
		return nil, nil
	}

	hostsInterface, err := wnx.Config.Slice(path + ".hosts")
	if err != nil {
		return nil, err
	}

	hostsString := make([]string, len(hostsInterface))
	for key, value := range hostsInterface {
		hostString, ok := value.(string)
		if !ok {
			return nil, errors.New("sefsef")
		}

		hostsString[key] = hostString
	}

	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(hostsString...),
	}

	dirCache, err := wnx.Config.String(path + ".dirCache")
	if err == nil {
		certManager.Cache = autocert.DirCache(dirCache)
	} else if err != ErrConfigValueNotFound {
		return nil, err
	}

	return certManager, nil
}
