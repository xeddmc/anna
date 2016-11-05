package id

import (
	"github.com/xh3b4sd/anna/service/random"
	"github.com/xh3b4sd/anna/service/spec"
)

const (
	// Hex128 creates a new hexa decimal encoded, pseudo random, 128 bit hash.
	Hex128 spec.IDType = 16

	// Hex512 creates a new hexa decimal encoded, pseudo random, 512 bit hash.
	Hex512 spec.IDType = 64

	// Hex1024 creates a new hexa decimal encoded, pseudo random, 1024 bit hash.
	Hex1024 spec.IDType = 128

	// Hex2048 creates a new hexa decimal encoded, pseudo random, 2048 bit hash.
	Hex2048 spec.IDType = 256

	// Hex4096 creates a new hexa decimal encoded, pseudo random, 4096 bit hash.
	Hex4096 spec.IDType = 512
)

// ServiceConfig represents the configuration used to create a new ID service
// object.
type ServiceConfig struct {
	// Settings.

	// HashChars represents the characters used to create hashes.
	HashChars string

	// Random represents a service returning random numbers.
	RandomService spec.Random

	Type spec.IDType
}

// DefaultServiceConfig provides a default configuration to create a new ID service
// object by best effort.
func DefaultServiceConfig() ServiceConfig {
	newConfig := ServiceConfig{
		// Settings.
		HashChars:     "abcdef0123456789", // hex character set
		RandomService: random.MustNewService(),
		Type:          Hex128,
	}

	return newConfig
}

// NewService creates a new configured ID service object.
func NewService(config ServiceConfig) (spec.ID, error) {
	newService := &service{
		ServiceConfig: config,
	}

	if newService.HashChars == "" {
		return nil, maskAnyf(invalidConfigError, "hash characters must not be empty")
	}
	if newService.RandomService == nil {
		return nil, maskAnyf(invalidConfigError, "random service must not be empty")
	}
	if newService.Type == 0 {
		return nil, maskAnyf(invalidConfigError, "ID type must not be empty")
	}

	return newService, nil
}

// MustNewService creates either a new default configured id service object, or
// panics.
func MustNewService() spec.ID {
	newIDService, err := NewService(DefaultServiceConfig())
	if err != nil {
		panic(err)
	}

	return newIDService
}

type service struct {
	ServiceConfig
}

func (s *service) New() (string, error) {
	ID, err := s.WithType(s.Type)
	if err != nil {
		return "", maskAny(err)
	}

	return ID, nil
}

func (s *service) WithType(idType spec.IDType) (string, error) {
	n := int(idType)
	max := len(s.HashChars)

	newRandomNumbers, err := s.RandomService.CreateNMax(n, max)
	if err != nil {
		return "", maskAny(err)
	}

	b := make([]byte, n)

	for i, r := range newRandomNumbers {
		b[i] = s.HashChars[r]
	}

	return string(b), nil
}