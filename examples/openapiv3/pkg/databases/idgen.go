package databases

import (
	"crypto/rand"
	"log"

	"github.com/oklog/ulid/v2"
)

func NewID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(Now()), rand.Reader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func MustNewID() string {
	id, err := NewID()
	if err != nil {
		log.Fatalf("failed to generate id. reason %v\n", err)
	}

	return id
}
