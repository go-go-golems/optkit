package budget

import (
	"testing"

	"github.com/go-go-golems/optkit/record"
)

func TestReservationIdentityNormalizesClaimOrder(t *testing.T) {
	first, err := ReservationID(ReserveRequest{
		Campaign: "campaign:test", Work: "work:test",
		Claims: []Quantity{{Resource: "tokens", Units: 20}, {Resource: "calls", Units: 1}},
	})
	if err != nil {
		t.Fatal(err)
	}
	second, err := ReservationID(ReserveRequest{
		Campaign: "campaign:test", Work: "work:test",
		Claims: []Quantity{{Resource: "calls", Units: 1}, {Resource: "tokens", Units: 20}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("claim order changed reservation identity: %s != %s", first, second)
	}
	if err := record.ValidateID("reservation", string(first)); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeQuantitiesRejectsDuplicateAndNegative(t *testing.T) {
	if _, err := NormalizeQuantities([]Quantity{{Resource: "calls", Units: 1}, {Resource: "calls", Units: 2}}, false); err == nil {
		t.Fatal("duplicate resource unexpectedly accepted")
	}
	if _, err := NormalizeQuantities([]Quantity{{Resource: "calls", Units: -1}}, false); err == nil {
		t.Fatal("negative quantity unexpectedly accepted")
	}
}
