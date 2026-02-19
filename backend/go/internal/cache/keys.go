package cache

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	KeyPrefixPerson  = "person:"
	KeyPrefixOrg     = "org:"
	KeyPrefixMeeting = "meeting:"
	KeyPrefixSession    = "session:"
	KeyPrefixProfile    = "profile:"
	KeyPrefixIncrement  = "increment:"
	KeyPrefixAuth       = "auth:"
	KeyPrefixPermission = "permission:"
	KeyPrefixRole       = "role:"
	KeyPrefixConsent    = "consent:"
)

func KeyPerson(id uuid.UUID) string {
	return KeyPrefixPerson + id.String()
}

func KeyMeeting(id uuid.UUID) string {
	return KeyPrefixMeeting + id.String()
}

func KeyMeetingByExternalID(externalType, externalID string) string {
	return fmt.Sprintf("%sexternal:%s:%s", KeyPrefixMeeting, externalType, externalID)
}

func KeyPersonByEmail(email string) string {
	return KeyPrefixPerson + "email:" + email
}

func KeyOrganization(id uuid.UUID) string {
	return KeyPrefixOrg + id.String()
}

func KeyOrganizationBySlug(slug string) string {
	return KeyPrefixOrg + "slug:" + slug
}

// KeyOrganizationMeetings returns a cache key for paginated meetings in an org.
func KeyOrganizationMeetings(orgID uuid.UUID, page, pageSize int) string {
	return fmt.Sprintf("org:%s:meetings:page:%d:size:%d", orgID.String(), page, pageSize)
}

func KeySession(tokenHash string) string {
	return KeyPrefixSession + tokenHash
}

func KeyProfile(id uuid.UUID) string {
	return KeyPrefixProfile + id.String()
}

func KeyProfileByPersonAndOrg(personID, orgID uuid.UUID) string {
	return fmt.Sprintf("%sperson:%s:org:%s", KeyPrefixProfile, personID.String(), orgID.String())
}

func KeyIncrement(id uuid.UUID) string {
	return KeyPrefixIncrement + id.String()
}

func KeyMeetingIncrements(meetingID uuid.UUID) string {
	return fmt.Sprintf("meeting:%s:increments", meetingID.String())
}

func KeyAuthMethod(id uuid.UUID) string {
	return KeyPrefixAuth + id.String()
}

func KeyAuthMethodByProvider(provider, providerID string) string {
	return fmt.Sprintf("%sprovider:%s:%s", KeyPrefixAuth, provider, providerID)
}

func KeyRole(id uuid.UUID) string {
	return KeyPrefixRole + id.String()
}

func KeyPermission(id uuid.UUID) string {
	return KeyPrefixPermission + id.String()
}

func KeyHasPermission(personID, orgID uuid.UUID, resourceName string, resourceID *uuid.UUID, activity string) string {
	resIDStr := "nil"
	if resourceID != nil {
		resIDStr = resourceID.String()
	}
	return fmt.Sprintf("has_perm:%s:%s:%s:%s:%s", personID.String(), orgID.String(), resourceName, resIDStr, activity)
}

func KeyConsentBySession(sessionID string) string {
	return KeyPrefixConsent + "session:" + sessionID
}

func KeyConsentByPerson(personID uuid.UUID) string {
	return KeyPrefixConsent + "person:" + personID.String()
}

func ChannelMeetingEvents(meetingID uuid.UUID) string {
	return fmt.Sprintf("events:meeting:%s", meetingID.String())
}

