package record

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
)

var namespacePattern = regexp.MustCompile(`^[a-z][a-z0-9._-]*$`)

// NewID creates an opaque, namespaced identifier. Semantic entities should use
// content-derived IDs where possible; this helper is for historical instances,
// commands, leases, and other facts that require a nonce.
func NewID(namespace string) (string, error) {
	if !namespacePattern.MatchString(namespace) {
		return "", fmt.Errorf("invalid ID namespace %q", namespace)
	}
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("generate %s ID: %w", namespace, err)
	}
	return namespace + ":" + hex.EncodeToString(buf[:]), nil
}

func ContentID(namespace string, digest Digest) (string, error) {
	if !namespacePattern.MatchString(namespace) {
		return "", fmt.Errorf("invalid ID namespace %q", namespace)
	}
	hexPart, err := digest.Hex()
	if err != nil {
		return "", err
	}
	return namespace + ":" + hexPart, nil
}

func ValidateID(namespace, raw string) error {
	if !namespacePattern.MatchString(namespace) {
		return fmt.Errorf("invalid ID namespace %q", namespace)
	}
	prefix := namespace + ":"
	if !strings.HasPrefix(raw, prefix) || len(raw) == len(prefix) {
		return fmt.Errorf("ID %q must have namespace %q", raw, namespace)
	}
	return nil
}

type SchemaID string

type SystemID string

type SnapshotID string

type PatchID string

type CandidateID string

type CampaignID string

type EventID string

type CommandID string

type EpisodeID string

type TrialID string

type ObservationID string

type EpochID string

type WorkID string

type LeaseID string

type ReservationID string

type SpanID string

type ActorRef string

func (id SchemaID) Validate() error { return ValidateID("schema", string(id)) }
