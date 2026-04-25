package testutil

import "github.com/google/uuid"

var (
	fixedTenantID = uuid.MustParse("a1b2c3d4-e5f6-7890-abcd-ef1234567890")
	fixedUserID   = uuid.MustParse("f0e1d2c3-b4a5-6789-0fed-cba987654321")
)

func TestTenantID() uuid.UUID {
	return fixedTenantID
}

func TestUserID() uuid.UUID {
	return fixedUserID
}
