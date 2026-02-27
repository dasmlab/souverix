package store

import (
	"testing"

	"github.com/dasmlab/ims/pkg/ims"
	"github.com/sirupsen/logrus"
)

func TestNewMemHSSStore(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	store, err := NewMemHSSStore(log)
	if err != nil {
		t.Fatalf("NewMemHSSStore() error = %v", err)
	}

	if store == nil {
		t.Error("NewMemHSSStore() returned nil")
	}
}

func TestMemHSSStore_SubscriberOperations(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	store, _ := NewMemHSSStore(log)

	sub := &ims.Subscriber{
		IMPU: "sip:alice@ims.local",
		IMPI: "alice@ims.local",
		Registered: false,
		ServiceProfile: ims.ServiceProfile{
			PublicIdentities: []string{"sip:alice@ims.local"},
		},
	}

	// Upsert
	err := store.UpsertSubscriber(sub)
	if err != nil {
		t.Fatalf("UpsertSubscriber() error = %v", err)
	}

	// Get
	retrieved, err := store.GetSubscriber("alice@ims.local")
	if err != nil {
		t.Fatalf("GetSubscriber() error = %v", err)
	}

	if retrieved.IMPU != sub.IMPU {
		t.Errorf("GetSubscriber() IMPU = %v, want %v", retrieved.IMPU, sub.IMPU)
	}

	// Get by IMPU
	retrievedByIMPU, err := store.GetSubscriberByIMPU("sip:alice@ims.local")
	if err != nil {
		t.Fatalf("GetSubscriberByIMPU() error = %v", err)
	}

	if retrievedByIMPU.IMPI != sub.IMPI {
		t.Errorf("GetSubscriberByIMPU() IMPI = %v, want %v", retrievedByIMPU.IMPI, sub.IMPI)
	}

	// List
	subs, err := store.ListSubscribers()
	if err != nil {
		t.Fatalf("ListSubscribers() error = %v", err)
	}

	if len(subs) == 0 {
		t.Error("ListSubscribers() returned empty list")
	}
}

func TestMemHSSStore_RegistrationOperations(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	store, _ := NewMemHSSStore(log)

	reg := &ims.Registration{
		IMPI:    "alice@ims.local",
		IMPU:    "sip:alice@ims.local",
		Contact: "sip:alice@192.168.1.1:5060",
		Expires: 3600,
		State:   ims.RegistrationStateRegistered,
	}

	// Upsert
	err := store.UpsertRegistration(reg)
	if err != nil {
		t.Fatalf("UpsertRegistration() error = %v", err)
	}

	// Get
	retrieved, err := store.GetRegistration("alice@ims.local")
	if err != nil {
		t.Fatalf("GetRegistration() error = %v", err)
	}

	if retrieved.State != reg.State {
		t.Errorf("GetRegistration() State = %v, want %v", retrieved.State, reg.State)
	}
}

func TestMemHSSStore_SCSCFAssignment(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	store, _ := NewMemHSSStore(log)

	// Assign
	scscf, err := store.AssignSCSCF("alice@ims.local")
	if err != nil {
		t.Fatalf("AssignSCSCF() error = %v", err)
	}

	if scscf == "" {
		t.Error("AssignSCSCF() returned empty string")
	}

	// Get assignment
	retrieved, err := store.GetSCSCFForSubscriber("alice@ims.local")
	if err != nil {
		t.Fatalf("GetSCSCFForSubscriber() error = %v", err)
	}

	if retrieved != scscf {
		t.Errorf("GetSCSCFForSubscriber() = %v, want %v", retrieved, scscf)
	}
}
