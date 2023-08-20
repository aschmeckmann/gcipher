package repositories

import (
	"sync"
)

var (
	repoOnce      sync.Once
	certRepo      *CertificateRepository
	userRepo      *UserRepository
	crlRepo       *CRLRepository
	repoInitError error
)

// InitializeRepositories initializes the repositories once
func InitializeRepositories() error {
	repoOnce.Do(func() {
		certRepo, repoInitError = NewCertificateRepository()
		if repoInitError != nil {
			return
		}

		userRepo, repoInitError = NewUserRepository()
		if repoInitError != nil {
			return
		}

		crlRepo, repoInitError = NewCRLRepository()
		if repoInitError != nil {
			return
		}
	})

	return repoInitError
}

// GetCertificateRepository returns the singleton-like instance of the CertificateRepository
func GetCertificateRepository() *CertificateRepository {
	return certRepo
}

// GetUserRepository returns the singleton-like instance of the UserRepository
func GetUserRepository() *UserRepository {
	return userRepo
}

// G returns the singleton-like instance of the UserRepository
func GetCRLRepository() *CRLRepository {
	return crlRepo
}
