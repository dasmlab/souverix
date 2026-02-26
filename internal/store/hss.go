package store

import (
	"fmt"
	"sync"

	"github.com/dasmlab/ims/pkg/ims"
	"github.com/sirupsen/logrus"
)

// HSSStore is the Home Subscriber Server data store
type HSSStore interface {
	// Subscriber operations
	GetSubscriber(impi string) (*ims.Subscriber, error)
	GetSubscriberByIMPU(impu string) (*ims.Subscriber, error)
	UpsertSubscriber(sub *ims.Subscriber) error
	DeleteSubscriber(impi string) error
	ListSubscribers() ([]*ims.Subscriber, error)

	// Registration operations
	GetRegistration(impi string) (*ims.Registration, error)
	UpsertRegistration(reg *ims.Registration) error
	DeleteRegistration(impi string) error

	// S-CSCF assignment
	AssignSCSCF(impi string) (string, error)
	GetSCSCFForSubscriber(impi string) (string, error)
}

// MemHSSStore is an in-memory implementation of HSSStore
type MemHSSStore struct {
	mu           sync.RWMutex
	subscribers  map[string]*ims.Subscriber // key: IMPI
	registrations map[string]*ims.Registration // key: IMPI
	scscfPool    []string // Available S-CSCF names
	log          *logrus.Logger
}

// NewMemHSSStore creates a new in-memory HSS store
func NewMemHSSStore(log *logrus.Logger) (*MemHSSStore, error) {
	store := &MemHSSStore{
		subscribers:   make(map[string]*ims.Subscriber),
		registrations: make(map[string]*ims.Registration),
		scscfPool: []string{
			"scscf1.ims.local",
			"scscf2.ims.local",
		},
		log: log,
	}

	// Seed a test subscriber
	store.seedTestSubscriber()

	return store, nil
}

func (s *MemHSSStore) seedTestSubscriber() {
	sub := &ims.Subscriber{
		IMPU: "sip:alice@ims.local",
		IMPI: "alice@ims.local",
		Registered: false,
		ServiceProfile: ims.ServiceProfile{
			PublicIdentities: []string{"sip:alice@ims.local"},
			CoreNetworkServices: []string{"originating", "terminating"},
		},
		AuthData: ims.AuthData{
			AuthScheme: "Digest",
			Username:   "alice",
			Realm:      "ims.local",
			Password:   "secret123", // In production, this should be hashed
		},
	}
	s.subscribers[sub.IMPI] = sub
	s.log.WithField("impi", sub.IMPI).Info("seeded test subscriber")
}

// GetSubscriber retrieves a subscriber by IMPI
func (s *MemHSSStore) GetSubscriber(impi string) (*ims.Subscriber, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sub, ok := s.subscribers[impi]
	if !ok {
		return nil, fmt.Errorf("subscriber not found: %s", impi)
	}

	// Return a copy to prevent external modification
	subCopy := *sub
	return &subCopy, nil
}

// GetSubscriberByIMPU retrieves a subscriber by IMPU
func (s *MemHSSStore) GetSubscriberByIMPU(impu string) (*ims.Subscriber, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, sub := range s.subscribers {
		for _, pubID := range sub.ServiceProfile.PublicIdentities {
			if pubID == impu {
				subCopy := *sub
				return &subCopy, nil
			}
		}
	}

	return nil, fmt.Errorf("subscriber not found for IMPU: %s", impu)
}

// UpsertSubscriber creates or updates a subscriber
func (s *MemHSSStore) UpsertSubscriber(sub *ims.Subscriber) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if sub.IMPI == "" {
		return fmt.Errorf("IMPI is required")
	}

	// Create a copy
	subCopy := *sub
	s.subscribers[sub.IMPI] = &subCopy

	s.log.WithField("impi", sub.IMPI).Info("subscriber upserted")
	return nil
}

// DeleteSubscriber deletes a subscriber
func (s *MemHSSStore) DeleteSubscriber(impi string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subscribers, impi)
	delete(s.registrations, impi)

	s.log.WithField("impi", impi).Info("subscriber deleted")
	return nil
}

// ListSubscribers lists all subscribers
func (s *MemHSSStore) ListSubscribers() ([]*ims.Subscriber, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	subs := make([]*ims.Subscriber, 0, len(s.subscribers))
	for _, sub := range s.subscribers {
		subCopy := *sub
		subs = append(subs, &subCopy)
	}

	return subs, nil
}

// GetRegistration retrieves a registration by IMPI
func (s *MemHSSStore) GetRegistration(impi string) (*ims.Registration, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	reg, ok := s.registrations[impi]
	if !ok {
		return nil, fmt.Errorf("registration not found: %s", impi)
	}

	regCopy := *reg
	return &regCopy, nil
}

// UpsertRegistration creates or updates a registration
func (s *MemHSSStore) UpsertRegistration(reg *ims.Registration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if reg.IMPI == "" {
		return fmt.Errorf("IMPI is required")
	}

	regCopy := *reg
	s.registrations[reg.IMPI] = &regCopy

	s.log.WithFields(logrus.Fields{
		"impi": reg.IMPI,
		"impu": reg.IMPU,
		"state": reg.State,
	}).Info("registration upserted")

	return nil
}

// DeleteRegistration deletes a registration
func (s *MemHSSStore) DeleteRegistration(impi string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.registrations, impi)
	s.log.WithField("impi", impi).Info("registration deleted")
	return nil
}

// AssignSCSCF assigns an S-CSCF to a subscriber (round-robin)
func (s *MemHSSStore) AssignSCSCF(impi string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.scscfPool) == 0 {
		return "", fmt.Errorf("no S-CSCF available in pool")
	}

	// Simple round-robin assignment
	// In production, this would use more sophisticated logic
	assigned := s.scscfPool[0]
	
	// Update subscriber's S-CSCF assignment
	if sub, ok := s.subscribers[impi]; ok {
		sub.SCSCFName = assigned
	}

	s.log.WithFields(logrus.Fields{
		"impi": impi,
		"scscf": assigned,
	}).Info("S-CSCF assigned")

	return assigned, nil
}

// GetSCSCFForSubscriber retrieves the assigned S-CSCF for a subscriber
func (s *MemHSSStore) GetSCSCFForSubscriber(impi string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sub, ok := s.subscribers[impi]
	if !ok {
		return "", fmt.Errorf("subscriber not found: %s", impi)
	}

	if sub.SCSCFName == "" {
		return "", fmt.Errorf("no S-CSCF assigned for subscriber: %s", impi)
	}

	return sub.SCSCFName, nil
}
