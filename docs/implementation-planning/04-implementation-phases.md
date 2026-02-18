# Implementation Phases

## Overview

This document breaks down the implementation into 7 phases, with each phase containing detailed tickets that can be worked on step-by-step. Each ticket includes:
- Description
- Acceptance criteria
- Dependencies
- Implementation notes
- Estimated complexity

## Phase 1: Foundation & Data Models

**Goal**: Set up project structure, database schema, and core data models.

### Phase 1.1: Project Setup

#### Ticket 1.1.1: Initialize Go Backend Project
**Description**: Set up the Go backend project structure with proper module configuration.

**Acceptance Criteria**:
- [ ] Go module initialized in `backend/go/`
- [ ] `go.mod` file created with proper module path
- [ ] Basic directory structure created:
  - `backend/go/cmd/` - Application entry points
  - `backend/go/internal/` - Internal packages
  - `backend/go/internal/models/` - GORM models
  - `backend/go/internal/repository/` - Repository interfaces and implementations
  - `backend/go/internal/service/` - Service layer
  - `backend/go/internal/handler/` - HTTP handlers
  - `backend/go/internal/middleware/` - HTTP middleware
  - `backend/go/internal/config/` - Configuration
  - `backend/go/internal/cache/` - Cache abstractions
  - `backend/go/internal/errors/` - Error definitions
  - `backend/go/migrations/` - Database migrations
  - `backend/go/pkg/` - Public packages (if any)
- [ ] `.gitignore` entries for Go build artifacts
- [ ] Basic `README.md` in `backend/go/`

**Dependencies**: None

**Implementation Notes**:
- Use `go mod init github.com/yourorg/meeting-cost/backend/go`
- Follow Go project layout best practices
- Set up proper package naming conventions

**Files to Create**:
- `backend/go/go.mod`
- `backend/go/README.md`
- `backend/go/.gitignore`

---

#### Ticket 1.1.2: Initialize React Frontend Project
**Description**: Set up the React frontend project structure.

**Acceptance Criteria**:
- [ ] React project initialized in `frontend/react/`
- [ ] Basic directory structure created:
  - `frontend/react/src/` - Source code
  - `frontend/react/src/components/` - React components
  - `frontend/react/src/pages/` - Page components
  - `frontend/react/src/services/` - API client services
  - `frontend/react/src/hooks/` - Custom React hooks
  - `frontend/react/src/context/` - React context providers
  - `frontend/react/src/utils/` - Utility functions
  - `frontend/react/src/types/` - TypeScript types
  - `frontend/react/public/` - Static assets
- [ ] `package.json` configured with dependencies
- [ ] TypeScript configuration
- [ ] ESLint and Prettier configured
- [ ] Basic `README.md` in `frontend/react/`

**Dependencies**: None

**Implementation Notes**:
- Use `create-react-app` with TypeScript template or Vite
- Configure path aliases for cleaner imports
- Set up environment variable handling

**Files to Create**:
- `frontend/react/package.json`
- `frontend/react/tsconfig.json`
- `frontend/react/.eslintrc.json`
- `frontend/react/.prettierrc`
- `frontend/react/README.md`

---

#### Ticket 1.1.3: Set Up Development Infrastructure
**Description**: Create Docker Compose setup for local development.

**Acceptance Criteria**:
- [ ] `docker-compose.yml` in `infrastructure/docker/`
- [ ] PostgreSQL service configured
- [ ] Valkey/Redis service configured
- [ ] Backend service configured (for local development)
- [ ] Frontend service configured (for local development)
- [ ] Environment variable files (`.env.example`)
- [ ] Network configuration for service communication
- [ ] Volume mounts for development

**Dependencies**: None

**Implementation Notes**:
- Use official PostgreSQL and Redis images
- Configure health checks
- Set up proper port mappings
- Use environment variables for configuration

**Files to Create**:
- `infrastructure/docker/docker-compose.yml`
- `infrastructure/docker/.env.example`
- `infrastructure/docker/README.md`

---

### Phase 1.2: Database Schema & Models

#### Ticket 1.2.1: Create Core GORM Models
**Description**: Implement all GORM models as defined in [01-data-models.md](01-data-models.md).

**Acceptance Criteria**:
- [ ] `Person` model implemented
- [ ] `Organization` model implemented
- [ ] `PersonOrganizationProfile` model implemented
- [ ] `Role` model implemented
- [ ] `RoleAssignment` model implemented
- [ ] `Permission` model implemented
- [ ] All models include proper GORM tags
- [ ] All models include `CreatedAt`, `UpdatedAt`, `DeletedAt`
- [ ] All foreign key relationships defined
- [ ] Model validation tags added where appropriate

**Dependencies**: Ticket 1.1.1

**Implementation Notes**:
- Place models in `backend/go/internal/models/`
- Use `uuid.UUID` for primary keys
- Use `gorm.DeletedAt` for soft deletes
- Add JSON tags for API serialization
- Add validation tags using `github.com/go-playground/validator`

**Files to Create**:
- `backend/go/internal/models/person.go`
- `backend/go/internal/models/organization.go`
- `backend/go/internal/models/person_organization_profile.go`
- `backend/go/internal/models/role.go`
- `backend/go/internal/models/role_assignment.go`
- `backend/go/internal/models/permission.go`

---

#### Ticket 1.2.2: Create Meeting & Increment Models
**Description**: Implement meeting-related GORM models.

**Acceptance Criteria**:
- [ ] `Meeting` model implemented
- [ ] `Increment` model implemented
- [ ] `MeetingParticipant` model implemented
- [ ] All relationships defined correctly
- [ ] Deduplication fields included
- [ ] External ID fields for Zoom/Teams/Slack integration

**Dependencies**: Ticket 1.2.1

**Implementation Notes**:
- Include computed fields (TotalCost, TotalDuration, etc.)
- Add indexes via GORM tags
- Consider JSONB for flexible external ID storage

**Files to Create**:
- `backend/go/internal/models/meeting.go`
- `backend/go/internal/models/increment.go`
- `backend/go/internal/models/meeting_participant.go`

---

#### Ticket 1.2.3: Create Authentication Models
**Description**: Implement authentication-related GORM models.

**Acceptance Criteria**:
- [ ] `AuthMethod` model implemented
- [ ] `Session` model implemented
- [ ] OAuth provider support fields
- [ ] Token storage fields (encrypted)
- [ ] Email verification fields

**Dependencies**: Ticket 1.2.1

**Implementation Notes**:
- Use encryption for sensitive token fields
- Support multiple OAuth providers
- Include session expiration logic

**Files to Create**:
- `backend/go/internal/models/auth_method.go`
- `backend/go/internal/models/session.go`

---

#### Ticket 1.2.4: Create Subscription & Payment Models
**Description**: Implement subscription and payment models.

**Acceptance Criteria**:
- [ ] `Subscription` model implemented
- [ ] `Payment` model implemented
- [ ] Stripe integration fields
- [ ] Subscription status tracking
- [ ] Payment history tracking

**Dependencies**: Ticket 1.2.1

**Implementation Notes**:
- Include Stripe customer and subscription IDs
- Track payment status and receipts
- Support multiple subscription plans

**Files to Create**:
- `backend/go/internal/models/subscription.go`
- `backend/go/internal/models/payment.go`

---

#### Ticket 1.2.5: Create Audit Log Model
**Description**: Implement audit logging model.

**Acceptance Criteria**:
- [ ] `AuditLog` model implemented
- [ ] Actor tracking (PersonID, OrganizationID)
- [ ] Action and resource tracking
- [ ] Flexible details storage (JSONB)
- [ ] Timestamp indexing

**Dependencies**: Ticket 1.2.1

**Implementation Notes**:
- Use JSONB for flexible audit details
- Index on timestamp for efficient queries
- Consider partitioning for large-scale deployments

**Files to Create**:
- `backend/go/internal/models/audit_log.go`

---

#### Ticket 1.2.6: Create Cookie Consent Model
**Description**: Implement CookieConsent GORM model for GDPR/CCPA compliance.

**Acceptance Criteria**:
- [ ] `CookieConsent` model implemented
- [ ] All consent preference fields
- [ ] Session tracking
- [ ] Audit trail fields (PreviousConsentID, ConsentSource)
- [ ] Version tracking
- [ ] IP address and user agent tracking
- [ ] Proper indexes for audit queries

**Dependencies**: Ticket 1.2.1

**Implementation Notes**:
- Support both authenticated and anonymous users
- Maintain full audit trail via PreviousConsentID
- Track consent version for policy changes
- Never hard delete, only soft delete for auditability

**Files to Create**:
- `backend/go/internal/models/cookie_consent.go`

---

#### Ticket 1.2.7: Create Database Migration Files
**Description**: Create initial database migration SQL files.

**Acceptance Criteria**:
- [ ] Migration tool configured (`golang-migrate` or similar)
- [ ] Initial schema migration created
- [ ] All tables created with proper constraints
- [ ] All indexes created
- [ ] Foreign key constraints defined
- [ ] Check constraints added
- [ ] Rollback migration created

**Dependencies**: Tickets 1.2.1-1.2.6

**Implementation Notes**:
- Use versioned migrations
- Include both up and down migrations
- Test migrations on clean database
- Document migration dependencies

**Files to Create**:
- `backend/go/migrations/001_initial_schema.up.sql`
- `backend/go/migrations/001_initial_schema.down.sql`
- `backend/go/migrations/README.md`

---

#### Ticket 1.2.7: Database Connection & Configuration
**Description**: Set up database connection and configuration management.

**Acceptance Criteria**:
- [ ] Database configuration struct
- [ ] Connection pooling configured
- [ ] GORM database connection established
- [ ] Migration runner implemented
- [ ] Health check endpoint for database
- [ ] Environment variable configuration

**Dependencies**: Ticket 1.2.6

**Implementation Notes**:
- Use connection pooling (set max open, idle connections)
- Configure connection timeouts
- Support both local and production configurations
- Use `gorm.io/driver/postgres` driver

**Files to Create**:
- `backend/go/internal/config/database.go`
- `backend/go/internal/config/config.go`
- `backend/go/cmd/migrate/main.go`

---

## Phase 2: Backend Infrastructure & Contracts

**Goal**: Implement core infrastructure, repository layer, and service contracts.

### Phase 2.1: Core Infrastructure

#### Ticket 2.1.1: Error Handling Infrastructure
**Description**: Implement error handling system as defined in [03-backend-patterns.md](03-backend-patterns.md).

**Acceptance Criteria**:
- [ ] Custom error types defined
- [ ] Domain error struct implemented
- [ ] Error code constants
- [ ] Error wrapping utilities
- [ ] Error conversion to HTTP responses
- [ ] Error logging integration

**Dependencies**: Ticket 1.1.1

**Implementation Notes**:
- Follow error handling patterns from patterns document
- Use `errors.Is()` and `errors.As()` for error checking
- Include error details for debugging

**Files to Create**:
- `backend/go/internal/errors/domain.go`
- `backend/go/internal/errors/codes.go`
- `backend/go/internal/errors/helpers.go`

---

#### Ticket 2.1.2: Logging Infrastructure
**Description**: Set up structured logging system.

**Acceptance Criteria**:
- [ ] Logger interface defined
- [ ] Logger implementation (using `zap` or `logrus`)
- [ ] Log levels configured
- [ ] Structured logging format
- [ ] Request ID tracking
- [ ] Context-aware logging
- [ ] Log rotation configuration

**Dependencies**: Ticket 1.1.1

**Implementation Notes**:
- Use structured logging (JSON format)
- Include request IDs for traceability
- Support different log levels per environment

**Files to Create**:
- `backend/go/internal/logger/logger.go`
- `backend/go/internal/logger/zap.go` (or logrus implementation)
- `backend/go/internal/logger/middleware.go`

---

#### Ticket 2.1.3: Cache Infrastructure
**Description**: Implement cache abstraction layer.

**Acceptance Criteria**:
- [ ] Cache interface defined
- [ ] Valkey/Redis implementation
- [ ] Cache key generation utilities
- [ ] TTL configuration
- [ ] Cache connection management
- [ ] Health check for cache

**Dependencies**: Ticket 1.1.3

**Implementation Notes**:
- Abstract cache behind interface for testability
- Use `github.com/redis/go-redis` or similar
- Implement cache key patterns from patterns document
- Support both local and distributed caching

**Files to Create**:
- `backend/go/internal/cache/interface.go`
- `backend/go/internal/cache/redis.go`
- `backend/go/internal/cache/keys.go`

---

#### Ticket 2.1.4: Configuration Management
**Description**: Implement configuration management system.

**Acceptance Criteria**:
- [ ] Configuration struct defined
- [ ] Environment variable loading
- [ ] Configuration validation
- [ ] Default values
- [ ] Environment-specific configs
- [ ] Secret management support

**Dependencies**: Ticket 1.1.1

**Implementation Notes**:
- Use `viper` or similar for configuration
- Support `.env` files for local development
- Validate required configuration on startup
- Support AWS Secrets Manager for production

**Files to Create**:
- `backend/go/internal/config/config.go`
- `backend/go/internal/config/loader.go`
- `backend/go/internal/config/validation.go`

---

### Phase 2.2: Repository Layer

#### Ticket 2.2.1: Person Repository Implementation
**Description**: Implement PersonRepository interface with GORM.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Database queries use GORM
- [ ] Cache integration (cache-aside pattern)
- [ ] Soft delete support
- [ ] Pagination support
- [ ] Filter support
- [ ] Unit tests with mocks

**Dependencies**: Tickets 1.2.1, 2.1.3

**Implementation Notes**:
- Follow repository patterns from patterns document
- Use transactions where needed
- Implement proper error handling
- Cache frequently accessed data

**Files to Create**:
- `backend/go/internal/repository/person_repository.go`
- `backend/go/internal/repository/person_repository_test.go`

---

#### Ticket 2.2.2: Organization Repository Implementation
**Description**: Implement OrganizationRepository interface.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Member management methods
- [ ] Meeting queries
- [ ] Cache integration
- [ ] Soft delete support
- [ ] Unit tests

**Dependencies**: Tickets 1.2.1, 2.2.1

**Implementation Notes**:
- Handle organization-member relationships
- Support active/inactive member filtering
- Cache organization data

**Files to Create**:
- `backend/go/internal/repository/organization_repository.go`
- `backend/go/internal/repository/organization_repository_test.go`

---

#### Ticket 2.2.3: PersonOrganizationProfile Repository
**Description**: Implement PersonOrganizationProfileRepository.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Wage update methods
- [ ] Membership activation/deactivation
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.1, 2.2.1, 2.2.2

**Implementation Notes**:
- Handle wage privacy concerns
- Support membership history
- Cache profile data

**Files to Create**:
- `backend/go/internal/repository/person_org_profile_repository.go`
- `backend/go/internal/repository/person_org_profile_repository_test.go`

---

#### Ticket 2.2.4: Meeting Repository Implementation
**Description**: Implement MeetingRepository interface.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Increment management
- [ ] Participant management
- [ ] Deduplication queries
- [ ] External ID lookups
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.2, 2.2.2

**Implementation Notes**:
- Support meeting deduplication
- Handle increment ordering
- Cache active meetings

**Files to Create**:
- `backend/go/internal/repository/meeting_repository.go`
- `backend/go/internal/repository/meeting_repository_test.go`

---

#### Ticket 2.2.5: Increment Repository Implementation
**Description**: Implement IncrementRepository interface.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Batch operations
- [ ] Time-based queries
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.2, 2.2.4

**Implementation Notes**:
- Support efficient batch inserts
- Index on time ranges for queries

**Files to Create**:
- `backend/go/internal/repository/increment_repository.go`
- `backend/go/internal/repository/increment_repository_test.go`

---

#### Ticket 2.2.6: Auth Repository Implementation
**Description**: Implement AuthRepository interface.

**Acceptance Criteria**:
- [ ] AuthMethod operations
- [ ] Session operations
- [ ] Token hash lookups
- [ ] Expired session cleanup
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.3, 2.1.3

**Implementation Notes**:
- Cache active sessions
- Implement session expiration cleanup
- Support multiple auth methods per person

**Files to Create**:
- `backend/go/internal/repository/auth_repository.go`
- `backend/go/internal/repository/auth_repository_test.go`

---

#### Ticket 2.2.7: Permission Repository Implementation
**Description**: Implement PermissionRepository interface.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Role operations
- [ ] Permission operations
- [ ] Role assignment operations
- [ ] Permission checking logic
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.1, 2.2.1, 2.2.2

**Implementation Notes**:
- Implement efficient permission checking
- Cache permission results
- Support both role-based and direct permissions

**Files to Create**:
- `backend/go/internal/repository/permission_repository.go`
- `backend/go/internal/repository/permission_repository_test.go`

---

#### Ticket 2.2.8: Consent Repository Implementation
**Description**: Implement ConsentRepository interface for cookie consent management.

**Acceptance Criteria**:
- [ ] All interface methods implemented
- [ ] Current consent retrieval by session
- [ ] Current consent retrieval by person
- [ ] Consent history queries
- [ ] Audit trail maintenance
- [ ] Soft delete support
- [ ] Unit tests

**Dependencies**: Ticket 1.2.6, 2.1.3

**Implementation Notes**:
- Support both session-based and person-based queries
- Maintain full history for auditability
- Efficient queries for current consent

**Files to Create**:
- `backend/go/internal/repository/consent_repository.go`
- `backend/go/internal/repository/consent_repository_test.go`

**Acceptance Criteria**:
- [ ] Role operations
- [ ] Permission operations
- [ ] Role assignment operations
- [ ] Permission checking logic
- [ ] Cache integration
- [ ] Unit tests

**Dependencies**: Tickets 1.2.1, 2.2.1, 2.2.2

**Implementation Notes**:
- Implement efficient permission checking
- Cache permission results
- Support both role-based and direct permissions

**Files to Create**:
- `backend/go/internal/repository/permission_repository.go`
- `backend/go/internal/repository/permission_repository_test.go`

---

### Phase 2.3: Dependency Injection & Service Contracts

#### Ticket 2.3.1: Dependency Injection Container
**Description**: Implement DI container as defined in [03-backend-patterns.md](03-backend-patterns.md).

**Acceptance Criteria**:
- [ ] Container struct defined
- [ ] Repository initialization
- [ ] Service initialization
- [ ] Dependency wiring
- [ ] Context support
- [ ] Cleanup methods

**Dependencies**: All Phase 2.2 tickets

**Implementation Notes**:
- Follow container pattern from patterns document
- Initialize in dependency order
- Support graceful shutdown

**Files to Create**:
- `backend/go/internal/container/container.go`
- `backend/go/internal/container/factory.go`

---

#### Ticket 2.3.2: Service Interface Definitions
**Description**: Define all service interfaces as specified in [02-backend-contracts.md](02-backend-contracts.md).

**Acceptance Criteria**:
- [ ] AuthService interface
- [ ] PersonService interface
- [ ] OrganizationService interface
- [ ] MeetingService interface
- [ ] ConsentService interface
- [ ] All DTO types defined
- [ ] Request/Response types defined

**Dependencies**: Ticket 2.1.1

**Implementation Notes**:
- Match interfaces exactly from contracts document
- Include all method signatures
- Define all DTOs including consent DTOs

**Files to Create**:
- `backend/go/internal/service/auth_service.go` (interface)
- `backend/go/internal/service/person_service.go` (interface)
- `backend/go/internal/service/organization_service.go` (interface)
- `backend/go/internal/service/meeting_service.go` (interface)
- `backend/go/internal/service/consent_service.go` (interface)
- `backend/go/internal/service/dto.go` (all DTOs)

**Acceptance Criteria**:
- [ ] AuthService interface
- [ ] PersonService interface
- [ ] OrganizationService interface
- [ ] MeetingService interface
- [ ] All DTO types defined
- [ ] Request/Response types defined

**Dependencies**: Ticket 2.1.1

**Implementation Notes**:
- Match interfaces exactly from contracts document
- Include all method signatures
- Define all DTOs

**Files to Create**:
- `backend/go/internal/service/auth_service.go` (interface)
- `backend/go/internal/service/person_service.go` (interface)
- `backend/go/internal/service/organization_service.go` (interface)
- `backend/go/internal/service/meeting_service.go` (interface)
- `backend/go/internal/service/dto.go` (all DTOs)

---

## Phase 3: Authentication & Authorization

**Goal**: Implement authentication and authorization systems.

### Phase 3.1: Authentication Service

#### Ticket 3.1.1: JWT Token Management
**Description**: Implement JWT token generation and validation.

**Acceptance Criteria**:
- [ ] JWT token generation
- [ ] Token validation
- [ ] Token refresh logic
- [ ] Token expiration handling
- [ ] Secret key management
- [ ] Token payload structure
- [ ] Unit tests

**Dependencies**: Tickets 2.2.6, 2.3.2

**Implementation Notes**:
- Use `github.com/golang-jwt/jwt/v5`
- Store secrets securely
- Support token rotation
- Include session ID in token

**Files to Create**:
- `backend/go/internal/auth/jwt.go`
- `backend/go/internal/auth/jwt_test.go`

---

#### Ticket 3.1.2: Password Management
**Description**: Implement password hashing and validation.

**Acceptance Criteria**:
- [ ] Password hashing (bcrypt)
- [ ] Password validation
- [ ] Password strength requirements
- [ ] Password change functionality
- [ ] Unit tests

**Dependencies**: Ticket 2.2.6

**Implementation Notes**:
- Use `golang.org/x/crypto/bcrypt`
- Enforce password complexity rules
- Support password reset flow

**Files to Create**:
- `backend/go/internal/auth/password.go`
- `backend/go/internal/auth/password_test.go`

---

#### Ticket 3.1.3: Auth Service Implementation - Registration
**Description**: Implement user registration flow.

**Acceptance Criteria**:
- [ ] Register method implemented
- [ ] Email validation
- [ ] Password validation
- [ ] Duplicate email check
- [ ] Email verification token generation
- [ ] Person creation
- [ ] AuthMethod creation
- [ ] Unit tests

**Dependencies**: Tickets 3.1.2, 2.2.1, 2.2.6

**Implementation Notes**:
- Validate email format
- Check for existing email
- Create person and auth method
- Generate verification token

**Files to Create**:
- `backend/go/internal/service/auth_service_impl.go` (partial)

---

#### Ticket 3.1.4: Auth Service Implementation - Login
**Description**: Implement login flow.

**Acceptance Criteria**:
- [ ] Login method implemented
- [ ] Email/password authentication
- [ ] Session creation
- [ ] JWT token generation
- [ ] Refresh token generation
- [ ] Error handling
- [ ] Unit tests

**Dependencies**: Tickets 3.1.1, 3.1.2, 3.1.3

**Implementation Notes**:
- Verify credentials
- Create session record
- Generate access and refresh tokens
- Handle invalid credentials

**Files to Create**:
- `backend/go/internal/service/auth_service_impl.go` (partial)

---

#### Ticket 3.1.5: Auth Service Implementation - OAuth2
**Description**: Implement OAuth2 authentication flow.

**Acceptance Criteria**:
- [ ] OAuth provider configuration
- [ ] OAuth initiation endpoint
- [ ] OAuth callback handling
- [ ] Token exchange
- [ ] User info retrieval
- [ ] AuthMethod creation/update
- [ ] Support for Zoom, Google, Teams, Slack
- [ ] Unit tests

**Dependencies**: Tickets 3.1.1, 2.2.6

**Implementation Notes**:
- Use OAuth2 library
- Handle state parameter for CSRF
- Store OAuth tokens securely
- Support multiple providers

**Files to Create**:
- `backend/go/internal/auth/oauth.go`
- `backend/go/internal/auth/oauth_providers.go`
- `backend/go/internal/service/auth_service_impl.go` (partial)

---

#### Ticket 3.1.6: Auth Service Implementation - Password Reset
**Description**: Implement password reset flow.

**Acceptance Criteria**:
- [ ] Forgot password method
- [ ] Reset token generation
- [ ] Email sending (or token storage)
- [ ] Reset password method
- [ ] Token validation
- [ ] Token expiration
- [ ] Unit tests

**Dependencies**: Tickets 3.1.2, 2.2.6

**Implementation Notes**:
- Generate secure reset tokens
- Set token expiration (e.g., 1 hour)
- Validate token before allowing reset

**Files to Create**:
- `backend/go/internal/service/auth_service_impl.go` (partial)

---

#### Ticket 3.1.7: Session Management
**Description**: Implement session management functionality.

**Acceptance Criteria**:
- [ ] Session validation
- [ ] Session refresh
- [ ] Session revocation
- [ ] Multiple session support
- [ ] Session cleanup (expired)
- [ ] Unit tests

**Dependencies**: Tickets 3.1.1, 2.2.6

**Implementation Notes**:
- Validate session on each request
- Support revoking individual or all sessions
- Clean up expired sessions periodically

**Files to Create**:
- `backend/go/internal/service/auth_service_impl.go` (partial)

---

### Phase 3.2: Authorization & Permissions

#### Ticket 3.2.1: Permission Checking Logic
**Description**: Implement permission checking system.

**Acceptance Criteria**:
- [ ] Permission check method
- [ ] Role-based permission resolution
- [ ] Direct permission checking
- [ ] Resource-scoped permissions
- [ ] Organization-scoped permissions
- [ ] Caching of permission results
- [ ] Unit tests

**Dependencies**: Tickets 2.2.7, 2.1.3

**Implementation Notes**:
- Check both role and direct permissions
- Cache permission results
- Support resource-level permissions

**Files to Create**:
- `backend/go/internal/auth/permission_checker.go`
- `backend/go/internal/auth/permission_checker_test.go`

---

#### Ticket 3.2.2: Authorization Middleware
**Description**: Implement HTTP middleware for authorization.

**Acceptance Criteria**:
- [ ] JWT token extraction
- [ ] Token validation
- [ ] User context injection
- [ ] Permission checking middleware
- [ ] Organization context extraction
- [ ] Error responses for unauthorized/forbidden
- [ ] Unit tests

**Dependencies**: Tickets 3.1.1, 3.2.1

**Implementation Notes**:
- Extract token from Authorization header
- Validate token and extract user info
- Inject user into request context
- Support permission-based route protection

**Files to Create**:
- `backend/go/internal/middleware/auth.go`
- `backend/go/internal/middleware/permission.go`
- `backend/go/internal/middleware/auth_test.go`

---

#### Ticket 3.2.3: Role & Permission Management
**Description**: Implement role and permission management in OrganizationService.

**Acceptance Criteria**:
- [ ] Create role method
- [ ] Assign role method
- [ ] Remove role assignment
- [ ] Permission assignment to roles
- [ ] Permission assignment to persons
- [ ] Role listing
- [ ] Unit tests

**Dependencies**: Tickets 2.2.7, 2.3.2

**Implementation Notes**:
- Support creating custom roles
- Assign permissions to roles
- Support direct person permissions
- Validate permissions exist

**Files to Create**:
- `backend/go/internal/service/organization_service_impl.go` (partial)

---

### Phase 3.3: Cookie Consent Service

#### Ticket 3.3.1: Cookie Classification System
**Description**: Implement system to classify cookies by category for consent enforcement.

**Acceptance Criteria**:
- [ ] Cookie classification rules engine
- [ ] Category mapping (necessary, analytics, marketing, functional)
- [ ] Pattern matching for cookie names
- [ ] Configuration-based classification rules
- [ ] Default classification for unknown cookies
- [ ] Cookie name pattern matching (regex support)
- [ ] Unit tests with various cookie examples

**Dependencies**: Ticket 2.3.2

**Implementation Notes**:
- Define rules for each cookie category
- Support pattern matching (e.g., "_ga*" = analytics, "session*" = necessary)
- Allow configuration override via config file
- Default to "necessary" if classification unknown (fail-safe)
- Support wildcard patterns and regex
- Document all cookie classifications

**Files to Create**:
- `backend/go/internal/consent/cookie_classifier.go`
- `backend/go/internal/consent/cookie_classifier_test.go`
- `backend/go/internal/consent/cookie_rules.go`
- `backend/go/internal/consent/cookie_rules.yaml` (or .json)

---

#### Ticket 3.3.2: Consent Service Implementation
**Description**: Implement ConsentService for cookie consent management with auditability.

**Acceptance Criteria**:
- [ ] GetConsent method
- [ ] UpdateConsent method
- [ ] WithdrawConsent method
- [ ] CheckCookieAllowed method (runtime enforcement)
- [ ] ClassifyCookie method (uses classification system)
- [ ] GetConsentHistory method
- [ ] ExportConsentData method
- [ ] GetCurrentPolicyVersion method
- [ ] Audit trail creation on consent changes
- [ ] Session ID generation for anonymous users
- [ ] IP address and user agent tracking
- [ ] Integration with cookie classifier
- [ ] Unit tests

**Dependencies**: Tickets 2.2.8, 2.3.2, 2.1.2, 3.3.1

**Implementation Notes**:
- Create new consent record on each change (never update existing)
- Link to previous consent via PreviousConsentID
- Track consent source (initial, update, withdrawal)
- Support both authenticated and anonymous users
- Maintain full audit trail for compliance
- CheckCookieAllowed should check consent and return boolean
- Always allow "necessary" cookies regardless of consent
- Cache consent lookups for performance

**Files to Create**:
- `backend/go/internal/service/consent_service_impl.go`
- `backend/go/internal/service/consent_service_impl_test.go`

**Acceptance Criteria**:
- [ ] Create role method
- [ ] Assign role method
- [ ] Remove role assignment
- [ ] Permission assignment to roles
- [ ] Permission assignment to persons
- [ ] Role listing
- [ ] Unit tests

**Dependencies**: Tickets 2.2.7, 2.3.2

**Implementation Notes**:
- Support creating custom roles
- Assign permissions to roles
- Support direct person permissions
- Validate permissions exist

**Files to Create**:
- `backend/go/internal/service/organization_service_impl.go` (partial)

---

## Phase 4: Core Business Logic Services

**Goal**: Implement core business logic for meetings, organizations, and persons.

### Phase 4.1: Person Service

#### Ticket 4.1.1: Person Service Implementation - CRUD
**Description**: Implement basic CRUD operations for PersonService.

**Acceptance Criteria**:
- [ ] GetPerson method
- [ ] UpdatePerson method
- [ ] Authorization checks
- [ ] Input validation
- [ ] DTO conversion
- [ ] Unit tests

**Dependencies**: Tickets 2.2.1, 2.3.2, 3.2.1

**Implementation Notes**:
- Check permissions before operations
- Validate input data
- Convert models to DTOs
- Handle not found errors

**Files to Create**:
- `backend/go/internal/service/person_service_impl.go` (partial)

---

#### Ticket 4.1.2: Person Service - Organization Management
**Description**: Implement organization membership methods.

**Acceptance Criteria**:
- [ ] GetOrganizations method
- [ ] JoinOrganization method
- [ ] LeaveOrganization method
- [ ] Authorization checks
- [ ] Unit tests

**Dependencies**: Tickets 2.2.1, 2.2.2, 2.2.3, 4.1.1

**Implementation Notes**:
- Check if person can join organization
- Handle membership activation/deactivation
- Update cache on changes

**Files to Create**:
- `backend/go/internal/service/person_service_impl.go` (partial)

---

#### Ticket 4.1.3: Person Service - GDPR Compliance
**Description**: Implement GDPR compliance features.

**Acceptance Criteria**:
- [ ] RequestDataExport method
- [ ] RequestDeletion method
- [ ] Data anonymization
- [ ] Export all person data
- [ ] Anonymize without affecting meeting costs
- [ ] Unit tests

**Dependencies**: Tickets 4.1.1, 2.2.1

**Implementation Notes**:
- Export all person-related data
- Anonymize person data
- Preserve meeting cost calculations
- Soft delete with anonymization flag

**Files to Create**:
- `backend/go/internal/service/person_service_impl.go` (partial)

---

### Phase 4.2: Organization Service

#### Ticket 4.2.1: Organization Service - CRUD
**Description**: Implement basic CRUD operations for OrganizationService.

**Acceptance Criteria**:
- [ ] CreateOrganization method
- [ ] GetOrganization method
- [ ] UpdateOrganization method
- [ ] DeleteOrganization method
- [ ] Authorization checks
- [ ] Slug generation
- [ ] Unit tests

**Dependencies**: Tickets 2.2.2, 2.3.2, 3.2.1

**Implementation Notes**:
- Generate unique slug from name
- Check permissions for all operations
- Handle soft deletes
- Update cache

**Files to Create**:
- `backend/go/internal/service/organization_service_impl.go` (partial)

---

#### Ticket 4.2.2: Organization Service - Member Management
**Description**: Implement member management methods.

**Acceptance Criteria**:
- [ ] GetMembers method
- [ ] AddMember method
- [ ] RemoveMember method
- [ ] UpdateMemberWage method
- [ ] Authorization checks (wage privacy)
- [ ] Unit tests

**Dependencies**: Tickets 2.2.2, 2.2.3, 4.2.1

**Implementation Notes**:
- Check permissions before viewing/updating wages
- Handle wage privacy concerns
- Support organization default wage
- Support blended wage option

**Files to Create**:
- `backend/go/internal/service/organization_service_impl.go` (partial)

---

#### Ticket 4.2.3: Organization Service - Settings
**Description**: Implement organization settings management.

**Acceptance Criteria**:
- [ ] UpdateSettings method
- [ ] UpdateDefaultWage method
- [ ] SetBlendedWage method
- [ ] Authorization checks
- [ ] Unit tests

**Dependencies**: Tickets 4.2.1

**Implementation Notes**:
- Support flexible settings via JSONB
- Validate wage values
- Update cache on changes

**Files to Create**:
- `backend/go/internal/service/organization_service_impl.go` (partial)

---

### Phase 4.3: Meeting Service

#### Ticket 4.3.1: Meeting Service - Cost Calculation Logic
**Description**: Implement meeting cost calculation logic.

**Acceptance Criteria**:
- [ ] Increment cost calculation
- [ ] Total cost calculation
- [ ] Cost per second/minute/hour
- [ ] Wage aggregation logic
- [ ] Blended wage support
- [ ] Individual wage support
- [ ] Immutable value objects
- [ ] Unit tests

**Dependencies**: Tickets 1.2.2, 2.2.3

**Implementation Notes**:
- Follow immutability patterns
- Calculate costs based on increments
- Support both blended and individual wages
- Use value objects for money calculations

**Files to Create**:
- `backend/go/internal/service/meeting_cost_calculator.go`
- `backend/go/internal/service/meeting_cost_calculator_test.go`
- `backend/go/internal/valueobject/money.go`
- `backend/go/internal/valueobject/increment_cost.go`

---

#### Ticket 4.3.2: Meeting Service - CRUD Operations
**Description**: Implement basic CRUD operations for MeetingService.

**Acceptance Criteria**:
- [ ] CreateMeeting method
- [ ] GetMeeting method
- [ ] UpdateMeeting method
- [ ] DeleteMeeting method
- [ ] Authorization checks
- [ ] Deduplication logic
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 2.3.2, 3.2.1, 4.3.1

**Implementation Notes**:
- Check for duplicate meetings
- Generate deduplication hash
- Check permissions
- Handle external IDs

**Files to Create**:
- `backend/go/internal/service/meeting_service_impl.go` (partial)

---

#### Ticket 4.3.3: Meeting Service - Meeting Control
**Description**: Implement meeting start/stop/reset functionality.

**Acceptance Criteria**:
- [ ] StartMeeting method
- [ ] StopMeeting method
- [ ] ResetMeeting method
- [ ] State validation
- [ ] Increment management
- [ ] Transaction support
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 2.2.5, 4.3.1, 4.3.2

**Implementation Notes**:
- Validate meeting state before operations
- Use transactions for atomicity
- Create initial increment on start
- Close final increment on stop
- Recalculate totals

**Files to Create**:
- `backend/go/internal/service/meeting_service_impl.go` (partial)

---

#### Ticket 4.3.4: Meeting Service - Increment Management
**Description**: Implement increment update logic.

**Acceptance Criteria**:
- [ ] UpdateAttendeeCount method
- [ ] UpdateAverageWage method
- [ ] UpdatePurpose method
- [ ] Increment closing logic
- [ ] New increment creation
- [ ] Cost recalculation
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 2.2.5, 4.3.1, 4.3.3

**Implementation Notes**:
- Close current increment on change
- Create new increment with new values
- Recalculate running totals
- Use transactions

**Files to Create**:
- `backend/go/internal/service/meeting_service_impl.go` (partial)

---

#### Ticket 4.3.5: Meeting Service - Deduplication
**Description**: Implement meeting deduplication logic.

**Acceptance Criteria**:
- [ ] Deduplication hash generation
- [ ] Duplicate detection
- [ ] External ID matching
- [ ] Merge logic (optional)
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 4.3.2

**Implementation Notes**:
- Generate hash from external ID and organization
- Check for existing meetings
- Support multiple external providers
- Handle conflicts

**Files to Create**:
- `backend/go/internal/service/meeting_deduplicator.go`
- `backend/go/internal/service/meeting_deduplicator_test.go`

---

#### Ticket 4.3.6: Meeting Service - Participant Management
**Description**: Implement participant tracking.

**Acceptance Criteria**:
- [ ] AddParticipant method
- [ ] RemoveParticipant method
- [ ] GetParticipants method
- [ ] Join/leave time tracking
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 4.3.2

**Implementation Notes**:
- Track participant join/leave times
- Calculate participant duration
- Support bulk participant operations

**Files to Create**:
- `backend/go/internal/service/meeting_service_impl.go` (partial)

---

#### Ticket 4.3.7: Meeting Service - Queries
**Description**: Implement meeting query methods.

**Acceptance Criteria**:
- [ ] ListMeetings method
- [ ] GetMeetingCost method
- [ ] Filtering support
- [ ] Pagination support
- [ ] Sorting support
- [ ] Unit tests

**Dependencies**: Tickets 2.2.4, 4.3.1, 4.3.2

**Implementation Notes**:
- Support various filters
- Implement efficient pagination
- Cache query results
- Support sorting options

**Files to Create**:
- `backend/go/internal/service/meeting_service_impl.go` (partial)

---

## Phase 5: API Endpoints & HTTP Layer

**Goal**: Implement HTTP API layer with Fiber framework.

### Phase 5.1: API Infrastructure

#### Ticket 5.1.1: Fiber Application Setup
**Description**: Set up Fiber application with basic configuration.

**Acceptance Criteria**:
- [ ] Fiber app initialization
- [ ] Middleware setup (CORS, logging, recovery)
- [ ] Error handling middleware
- [ ] Request ID middleware
- [ ] Health check endpoint
- [ ] Graceful shutdown

**Dependencies**: Tickets 2.1.2, 2.1.4

**Implementation Notes**:
- Configure CORS for frontend
- Set up structured logging
- Implement graceful shutdown
- Add health checks
- Include cookie consent middleware in middleware chain (after Ticket 5.1.4)

**Files to Create**:
- `backend/go/cmd/api/main.go`
- `backend/go/internal/server/server.go`
- `backend/go/internal/middleware/cors.go`
- `backend/go/internal/middleware/recovery.go`
- `backend/go/internal/middleware/request_id.go`

---

#### Ticket 5.1.2: API Routing Structure
**Description**: Set up API routing structure with versioning.

**Acceptance Criteria**:
- [ ] Versioned routes (`/api/v1/`)
- [ ] Route groups (auth, persons, organizations, meetings)
- [ ] Route registration
- [ ] 404 handler
- [ ] Method not allowed handler

**Dependencies**: Ticket 5.1.1

**Implementation Notes**:
- Organize routes by resource
- Support API versioning
- Handle 404 and 405 errors

**Files to Create**:
- `backend/go/internal/handler/router.go`
- `backend/go/internal/handler/routes.go`

---

#### Ticket 5.1.3: Request Validation Middleware
**Description**: Implement request validation using validator.

**Acceptance Criteria**:
- [ ] Request validation middleware
- [ ] Validation error formatting
- [ ] Custom validators
- [ ] DTO validation tags
- [ ] Unit tests

**Dependencies**: Ticket 5.1.1

**Implementation Notes**:
- Use `github.com/go-playground/validator`
- Create custom validators as needed
- Format validation errors consistently

**Files to Create**:
- `backend/go/internal/middleware/validation.go`
- `backend/go/internal/validator/custom.go`

---

#### Ticket 5.1.4: Cookie Consent Enforcement Middleware
**Description**: Implement middleware to enforce cookie consent at runtime by filtering Set-Cookie headers.

**Acceptance Criteria**:
- [ ] Middleware that intercepts Set-Cookie headers
- [ ] Session ID extraction from request (cookie or header)
- [ ] Consent lookup from cache/database
- [ ] Cookie classification using consent service
- [ ] Automatic filtering of Set-Cookie headers based on consent
- [ ] Always allow necessary cookies (session, auth)
- [ ] Block analytics/marketing/functional cookies if not consented
- [ ] Fallback behavior when consent not found (default to necessary only)
- [ ] Logging of blocked cookie attempts for audit
- [ ] Performance optimization (cache consent lookups)
- [ ] Unit tests
- [ ] Integration tests

**Dependencies**: Tickets 3.3.1, 3.3.2, 5.1.1, 2.1.3

**Implementation Notes**:
- Intercept Set-Cookie headers in Fiber response
- Extract session ID from request (check cookie, header, or generate)
- Use ConsentService.CheckCookieAllowed() for each cookie
- Classify each cookie using ConsentService.ClassifyCookie()
- Remove Set-Cookie headers for non-allowed cookies
- Log blocked cookies for compliance auditing
- Cache consent lookups to minimize database hits
- Handle edge cases: no session ID, no consent record, expired consent
- Ensure necessary cookies (session, JWT) always allowed

**Files to Create**:
- `backend/go/internal/middleware/cookie_consent.go`
- `backend/go/internal/middleware/cookie_consent_test.go`
- `backend/go/internal/middleware/cookie_filter.go`

---

### Phase 5.2: Consent Endpoints

#### Ticket 5.2.1: Consent Handlers
**Description**: Implement cookie consent API endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/consent` endpoint
- [ ] POST `/api/v1/consent` endpoint
- [ ] DELETE `/api/v1/consent` endpoint (withdraw)
- [ ] GET `/api/v1/consent/history` endpoint
- [ ] GET `/api/v1/consent/policy-version` endpoint
- [ ] Request validation
- [ ] IP address and user agent extraction
- [ ] Session ID handling
- [ ] Integration with consent service
- [ ] Integration tests

**Dependencies**: Tickets 3.3.2, 5.1.2, 5.1.3

**Implementation Notes**:
- Extract IP and user agent from request
- Support both authenticated and anonymous users
- Validate consent preferences
- Return audit trail information
- Note: Cookie enforcement happens at middleware level (Ticket 5.1.4), not in handlers

**Files to Create**:
- `backend/go/internal/handler/consent_handler.go`
- `backend/go/internal/handler/consent_handler_test.go`

---

### Phase 5.3: Authentication Endpoints

#### Ticket 5.3.1: Auth Handlers - Registration & Login
**Description**: Implement registration and login endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/auth/register` endpoint
- [ ] POST `/api/v1/auth/login` endpoint
- [ ] Request validation
- [ ] Error handling
- [ ] Response formatting
- [ ] Integration tests

**Dependencies**: Tickets 3.1.3, 3.1.4, 5.1.2, 5.1.3

**Implementation Notes**:
- Validate request bodies
- Call service methods
- Format responses consistently
- Handle errors appropriately

**Files to Create**:
- `backend/go/internal/handler/auth_handler.go` (partial)

---

#### Ticket 5.3.2: Auth Handlers - OAuth
**Description**: Implement OAuth endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/auth/oauth/{provider}` endpoint
- [ ] GET `/api/v1/auth/oauth/{provider}/callback` endpoint
- [ ] State parameter handling
- [ ] Error handling
- [ ] Integration tests

**Dependencies**: Tickets 3.1.5, 5.1.2

**Implementation Notes**:
- Generate state for CSRF protection
- Handle OAuth callback
- Exchange code for tokens

**Files to Create**:
- `backend/go/internal/handler/auth_handler.go` (partial)

---

#### Ticket 5.3.3: Auth Handlers - Password Management
**Description**: Implement password reset endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/auth/forgot-password` endpoint
- [ ] POST `/api/v1/auth/reset-password` endpoint
- [ ] POST `/api/v1/auth/change-password` endpoint
- [ ] Request validation
- [ ] Integration tests

**Dependencies**: Tickets 3.1.6, 5.1.2, 5.1.3

**Implementation Notes**:
- Validate email for forgot password
- Validate token for reset
- Require authentication for change password

**Files to Create**:
- `backend/go/internal/handler/auth_handler.go` (partial)

---

#### Ticket 5.3.4: Auth Handlers - Session Management
**Description**: Implement session management endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/auth/logout` endpoint
- [ ] GET `/api/v1/auth/sessions` endpoint
- [ ] DELETE `/api/v1/auth/sessions/{id}` endpoint
- [ ] DELETE `/api/v1/auth/sessions` (all) endpoint
- [ ] Integration tests

**Dependencies**: Tickets 3.1.7, 5.1.2

**Implementation Notes**:
- Require authentication
- Support revoking individual or all sessions

**Files to Create**:
- `backend/go/internal/handler/auth_handler.go` (partial)

---

### Phase 5.4: Person Endpoints

#### Ticket 5.4.1: Person Handlers - Profile
**Description**: Implement person profile endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/persons/me` endpoint
- [ ] PATCH `/api/v1/persons/me` endpoint
- [ ] GET `/api/v1/persons/me/organizations` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.1.1, 5.1.2, 3.2.2

**Implementation Notes**:
- Use authentication middleware
- Return current user's profile
- Support profile updates

**Files to Create**:
- `backend/go/internal/handler/person_handler.go` (partial)

---

#### Ticket 5.4.2: Person Handlers - GDPR
**Description**: Implement GDPR compliance endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/persons/me/export` endpoint
- [ ] DELETE `/api/v1/persons/me` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.1.3, 5.1.2, 3.2.2

**Implementation Notes**:
- Export all person data as JSON
- Anonymize person on deletion
- Require authentication

**Files to Create**:
- `backend/go/internal/handler/person_handler.go` (partial)

---

### Phase 5.5: Organization Endpoints

#### Ticket 5.5.1: Organization Handlers - CRUD
**Description**: Implement organization CRUD endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/organizations` endpoint
- [ ] GET `/api/v1/organizations/{id}` endpoint
- [ ] PATCH `/api/v1/organizations/{id}` endpoint
- [ ] DELETE `/api/v1/organizations/{id}` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.2.1, 5.1.2, 3.2.2

**Implementation Notes**:
- Check permissions for all operations
- Validate organization ownership
- Handle soft deletes

**Files to Create**:
- `backend/go/internal/handler/organization_handler.go` (partial)

---

#### Ticket 5.5.2: Organization Handlers - Members
**Description**: Implement organization member management endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/organizations/{id}/members` endpoint
- [ ] POST `/api/v1/organizations/{id}/members` endpoint
- [ ] DELETE `/api/v1/organizations/{id}/members/{person_id}` endpoint
- [ ] PATCH `/api/v1/organizations/{id}/members/{person_id}/wage` endpoint
- [ ] Authorization checks (wage privacy)
- [ ] Integration tests

**Dependencies**: Tickets 4.2.2, 5.1.2, 3.2.2

**Implementation Notes**:
- Check permissions before showing wages
- Support adding/removing members
- Validate wage updates

**Files to Create**:
- `backend/go/internal/handler/organization_handler.go` (partial)

---

#### Ticket 5.5.3: Organization Handlers - Roles & Permissions
**Description**: Implement role and permission management endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/organizations/{id}/roles` endpoint
- [ ] POST `/api/v1/organizations/{id}/roles` endpoint
- [ ] POST `/api/v1/organizations/{id}/roles/{role_id}/assign` endpoint
- [ ] DELETE `/api/v1/organizations/{id}/roles/{role_id}/assign` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 3.2.3, 5.1.2, 3.2.2

**Implementation Notes**:
- Support role creation
- Support role assignment
- Validate permissions

**Files to Create**:
- `backend/go/internal/handler/organization_handler.go` (partial)

---

### Phase 5.6: Meeting Endpoints

#### Ticket 5.6.1: Meeting Handlers - CRUD
**Description**: Implement meeting CRUD endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/meetings` endpoint
- [ ] GET `/api/v1/meetings/{id}` endpoint
- [ ] PATCH `/api/v1/meetings/{id}` endpoint
- [ ] DELETE `/api/v1/meetings/{id}` endpoint
- [ ] GET `/api/v1/meetings` (list) endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.3.2, 4.3.7, 5.1.2, 3.2.2

**Implementation Notes**:
- Support filtering and pagination
- Check organization permissions
- Handle deduplication

**Files to Create**:
- `backend/go/internal/handler/meeting_handler.go` (partial)

---

#### Ticket 5.6.2: Meeting Handlers - Control
**Description**: Implement meeting control endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/meetings/{id}/start` endpoint
- [ ] POST `/api/v1/meetings/{id}/stop` endpoint
- [ ] POST `/api/v1/meetings/{id}/reset` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.3.3, 5.1.2, 3.2.2

**Implementation Notes**:
- Validate meeting state
- Handle state transitions
- Return updated meeting

**Files to Create**:
- `backend/go/internal/handler/meeting_handler.go` (partial)

---

#### Ticket 5.6.3: Meeting Handlers - Updates
**Description**: Implement meeting update endpoints.

**Acceptance Criteria**:
- [ ] PATCH `/api/v1/meetings/{id}/attendee-count` endpoint
- [ ] PATCH `/api/v1/meetings/{id}/average-wage` endpoint
- [ ] PATCH `/api/v1/meetings/{id}/purpose` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.3.4, 5.1.2, 3.2.2

**Implementation Notes**:
- Validate meeting is active
- Handle increment closing/creation
- Return updated meeting

**Files to Create**:
- `backend/go/internal/handler/meeting_handler.go` (partial)

---

#### Ticket 5.6.4: Meeting Handlers - Participants
**Description**: Implement participant management endpoints.

**Acceptance Criteria**:
- [ ] POST `/api/v1/meetings/{id}/participants` endpoint
- [ ] DELETE `/api/v1/meetings/{id}/participants/{person_id}` endpoint
- [ ] GET `/api/v1/meetings/{id}/participants` endpoint
- [ ] Authorization checks
- [ ] Integration tests

**Dependencies**: Tickets 4.3.6, 5.1.2, 3.2.2

**Implementation Notes**:
- Track join/leave times
- Support bulk operations

**Files to Create**:
- `backend/go/internal/handler/meeting_handler.go` (partial)

---

#### Ticket 5.6.5: Meeting Handlers - Cost Queries
**Description**: Implement meeting cost query endpoints.

**Acceptance Criteria**:
- [ ] GET `/api/v1/meetings/{id}/cost` endpoint
- [ ] Real-time cost calculation
- [ ] Cost breakdown
- [ ] Integration tests

**Dependencies**: Tickets 4.3.1, 4.3.7, 5.1.2, 3.2.2

**Implementation Notes**:
- Calculate costs on-the-fly
- Return detailed cost breakdown
- Support active meeting cost calculation

**Files to Create**:
- `backend/go/internal/handler/meeting_handler.go` (partial)

---

## Phase 6: Frontend Implementation

**Goal**: Implement React frontend application.

### Phase 6.1: Frontend Infrastructure

#### Ticket 6.1.1: Accessibility Infrastructure Setup
**Description**: Set up accessibility infrastructure and utilities for ADA compliance.

**Acceptance Criteria**:
- [ ] Accessibility utility functions
- [ ] ARIA label helpers
- [ ] Focus management utilities
- [ ] Keyboard navigation helpers
- [ ] Screen reader announcement utilities
- [ ] Skip navigation component
- [ ] Focus trap component
- [ ] Accessibility testing setup (axe-core, jest-axe)
- [ ] TypeScript types for accessibility

**Dependencies**: Ticket 1.1.2

**Implementation Notes**:
- Use `react-aria` or similar for accessibility primitives
- Set up `@axe-core/react` for development accessibility checks
- Create reusable accessibility utilities
- Document accessibility patterns

**Files to Create**:
- `frontend/react/src/utils/accessibility.ts`
- `frontend/react/src/utils/aria.ts`
- `frontend/react/src/utils/focus.ts`
- `frontend/react/src/components/accessibility/SkipNavigation.tsx`
- `frontend/react/src/components/accessibility/FocusTrap.tsx`
- `frontend/react/src/components/accessibility/Announcer.tsx`
- `frontend/react/src/setupTests.ts` (with accessibility testing)

---

#### Ticket 6.1.2: API Client Setup
**Description**: Set up API client with authentication.

**Acceptance Criteria**:
- [ ] API client service
- [ ] Axios/fetch configuration
- [ ] Request interceptors (auth token)
- [ ] Response interceptors (error handling)
- [ ] Base URL configuration
- [ ] TypeScript types for API responses

**Dependencies**: Ticket 1.1.2

**Implementation Notes**:
- Use axios or fetch with interceptors
- Handle token refresh
- Format errors consistently

**Files to Create**:
- `frontend/react/src/services/api/client.ts`
- `frontend/react/src/services/api/interceptors.ts`
- `frontend/react/src/services/api/types.ts`

---

#### Ticket 6.1.3: State Management Setup
**Description**: Set up state management (Context API or Redux).

**Acceptance Criteria**:
- [ ] Auth context/provider
- [ ] User state management
- [ ] Organization state management
- [ ] Meeting state management
- [ ] State persistence (localStorage)
- [ ] TypeScript types

**Dependencies**: Ticket 6.1.1

**Implementation Notes**:
- Use Context API for simple state
- Consider Redux Toolkit for complex state
- Persist auth state

**Files to Create**:
- `frontend/react/src/context/AuthContext.tsx`
- `frontend/react/src/context/UserContext.tsx`
- `frontend/react/src/context/OrganizationContext.tsx`

---

#### Ticket 6.1.4: Routing Setup
**Description**: Set up React Router with protected routes.

**Acceptance Criteria**:
- [ ] React Router configuration
- [ ] Route definitions
- [ ] Protected route component
- [ ] Public route component
- [ ] Navigation components
- [ ] Route guards

**Dependencies**: Tickets 6.1.2, 6.1.3

**Implementation Notes**:
- Use React Router v6
- Protect routes based on auth
- Handle route redirects

**Files to Create**:
- `frontend/react/src/routes/Router.tsx`
- `frontend/react/src/routes/ProtectedRoute.tsx`
- `frontend/react/src/routes/PublicRoute.tsx`

---

### Phase 6.2: Cookie Consent Implementation

#### Ticket 6.2.1: Cookie Consent UI Component
**Description**: Implement cookie consent banner and preference management UI.

**Acceptance Criteria**:
- [ ] Cookie consent banner component
- [ ] Consent preference modal/dialog
- [ ] Granular cookie category controls
- [ ] Accept all / Reject all buttons
- [ ] Save preferences button
- [ ] Withdraw consent functionality
- [ ] Full ADA compliance (ARIA labels, keyboard navigation, screen reader support)
- [ ] Responsive design
- [ ] Animation/transition support
- [ ] Integration with consent service API

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Use accessible dialog component (react-aria-dialog or similar)
- Support keyboard navigation (Tab, Enter, Escape)
- Provide clear ARIA labels and descriptions
- Announce consent changes to screen readers
- Store consent in both backend and localStorage
- Show banner on first visit or when consent expires

**Files to Create**:
- `frontend/react/src/components/consent/CookieConsentBanner.tsx`
- `frontend/react/src/components/consent/ConsentPreferences.tsx`
- `frontend/react/src/components/consent/ConsentCategory.tsx`
- `frontend/react/src/hooks/useCookieConsent.ts`
- `frontend/react/src/services/api/consent.ts`

---

#### Ticket 6.2.2: Cookie Consent Service Integration
**Description**: Integrate frontend with consent service API.

**Acceptance Criteria**:
- [ ] API client methods for consent operations
- [ ] Session ID management
- [ ] Consent state synchronization
- [ ] Consent history retrieval
- [ ] Error handling
- [ ] Loading states
- [ ] TypeScript types

**Dependencies**: Tickets 6.1.2, 6.2.1

**Implementation Notes**:
- Generate and persist session ID
- Sync consent state with backend
- Handle consent updates
- Support consent withdrawal

**Files to Create**:
- `frontend/react/src/services/api/consent.ts` (partial)
- `frontend/react/src/types/consent.ts`

---

#### Ticket 6.2.3: Cookie Management Utilities
**Description**: Implement utilities for managing cookies based on consent.

**Acceptance Criteria**:
- [ ] Cookie setting/getting utilities
- [ ] Cookie deletion utilities
- [ ] Consent-based cookie management
- [ ] Cookie category filtering
- [ ] Cookie expiration handling
- [ ] Unit tests

**Dependencies**: Ticket 6.2.1

**Implementation Notes**:
- Only set cookies based on consent preferences
- Support cookie categories (necessary, analytics, marketing, functional)
- Automatically remove cookies when consent is withdrawn
- Handle cookie expiration

**Files to Create**:
- `frontend/react/src/utils/cookies.ts`
- `frontend/react/src/utils/cookieManager.ts`
- `frontend/react/src/utils/cookies.test.ts`

---

### Phase 6.3: Authentication UI

#### Ticket 6.3.1: Login Page
**Description**: Implement login page component.

**Acceptance Criteria**:
- [ ] Login form component
- [ ] Email/password inputs
- [ ] OAuth provider buttons
- [ ] Form validation
- [ ] Error handling
- [ ] Loading states
- [ ] Responsive design

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Use form library (react-hook-form)
- Validate inputs
- Handle OAuth redirects
- Show loading during auth

**Files to Create**:
- `frontend/react/src/pages/Login.tsx`
- `frontend/react/src/components/auth/LoginForm.tsx`

---

#### Ticket 6.3.2: Registration Page
**Description**: Implement registration page component.

**Acceptance Criteria**:
- [ ] Registration form
- [ ] Input validation
- [ ] Password strength indicator
- [ ] Terms acceptance
- [ ] Error handling
- [ ] Success handling

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Validate email format
- Check password strength
- Handle registration errors

**Files to Create**:
- `frontend/react/src/pages/Register.tsx`
- `frontend/react/src/components/auth/RegisterForm.tsx`

---

#### Ticket 6.3.3: Password Reset Flow
**Description**: Implement password reset pages.

**Acceptance Criteria**:
- [ ] Forgot password page
- [ ] Reset password page
- [ ] Email input validation
- [ ] Token validation
- [ ] Password reset form
- [ ] Success/error messages

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Handle reset token from URL
- Validate token before showing form
- Confirm password match

**Files to Create**:
- `frontend/react/src/pages/ForgotPassword.tsx`
- `frontend/react/src/pages/ResetPassword.tsx`

---

### Phase 6.4: Dashboard & Organization UI

#### Ticket 6.4.1: Dashboard Page
**Description**: Implement main dashboard page.

**Acceptance Criteria**:
- [ ] Dashboard layout
- [ ] Organization selector
- [ ] Quick stats
- [ ] Recent meetings
- [ ] Navigation menu
- [ ] Responsive design

**Dependencies**: Tickets 6.1.2, 6.1.3

**Implementation Notes**:
- Show user's organizations
- Display meeting statistics
- Quick access to common actions

**Files to Create**:
- `frontend/react/src/pages/Dashboard.tsx`
- `frontend/react/src/components/dashboard/OrganizationSelector.tsx`
- `frontend/react/src/components/dashboard/QuickStats.tsx`

---

#### Ticket 6.4.2: Organization Management Pages
**Description**: Implement organization management UI.

**Acceptance Criteria**:
- [ ] Organization list page
- [ ] Create organization page
- [ ] Organization settings page
- [ ] Member management page
- [ ] Role management page
- [ ] CRUD operations

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Support creating organizations
- Manage members and roles
- Update organization settings

**Files to Create**:
- `frontend/react/src/pages/Organizations.tsx`
- `frontend/react/src/pages/OrganizationDetail.tsx`
- `frontend/react/src/pages/OrganizationMembers.tsx`
- `frontend/react/src/components/organization/OrganizationForm.tsx`
- `frontend/react/src/components/organization/MemberList.tsx`

---

### Phase 6.5: Meeting UI

#### Ticket 6.5.1: Meeting List Page
**Description**: Implement meeting list/history page.

**Acceptance Criteria**:
- [ ] Meeting list component
- [ ] Filtering (date, organization)
- [ ] Pagination
- [ ] Sorting
- [ ] Meeting cards
- [ ] Cost display
- [ ] Responsive design

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Display past meetings
- Show meeting costs
- Support filtering and sorting
- Paginate results

**Files to Create**:
- `frontend/react/src/pages/Meetings.tsx`
- `frontend/react/src/components/meeting/MeetingList.tsx`
- `frontend/react/src/components/meeting/MeetingCard.tsx`
- `frontend/react/src/components/meeting/MeetingFilters.tsx`

---

#### Ticket 6.5.2: Meeting Timer Component
**Description**: Implement meeting timer and cost calculator.

**Acceptance Criteria**:
- [ ] Timer display
- [ ] Real-time cost calculation
- [ ] Attendee count input
- [ ] Average wage input
- [ ] Purpose input
- [ ] Start/stop/reset buttons
- [ ] Cost per second/minute/hour display
- [ ] Increment history display

**Dependencies**: Tickets 6.1.1, 6.1.2

**Implementation Notes**:
- Update cost in real-time
- Handle increment changes
- Display running totals
- Support meeting control

**Files to Create**:
- `frontend/react/src/components/meeting/MeetingTimer.tsx`
- `frontend/react/src/components/meeting/CostDisplay.tsx`
- `frontend/react/src/components/meeting/IncrementHistory.tsx`
- `frontend/react/src/hooks/useMeetingTimer.ts`

---

#### Ticket 6.5.3: Meeting Detail Page
**Description**: Implement meeting detail/view page.

**Acceptance Criteria**:
- [ ] Meeting information display
- [ ] Increment breakdown
- [ ] Participant list
- [ ] Cost breakdown
- [ ] Edit meeting (if permitted)
- [ ] Delete meeting (if permitted)

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3, 6.4.2

**Implementation Notes**:
- Show all meeting details
- Display increment history
- Show participants
- Support editing

**Files to Create**:
- `frontend/react/src/pages/MeetingDetail.tsx`
- `frontend/react/src/components/meeting/MeetingInfo.tsx`
- `frontend/react/src/components/meeting/IncrementBreakdown.tsx`

---

#### Ticket 6.5.4: Create Meeting Page
**Description**: Implement create meeting page.

**Acceptance Criteria**:
- [ ] Create meeting form
- [ ] Organization selection
- [ ] Purpose input
- [ ] External meeting ID (optional)
- [ ] Form validation
- [ ] Success handling

**Dependencies**: Tickets 6.1.1, 6.1.2, 6.1.3

**Implementation Notes**:
- Support creating meetings
- Handle external IDs
- Validate inputs

**Files to Create**:
- `frontend/react/src/pages/CreateMeeting.tsx`
- `frontend/react/src/components/meeting/CreateMeetingForm.tsx`

---

### Phase 6.6: ADA Compliance Implementation

#### Ticket 6.6.1: ADA Compliance - Form Components
**Description**: Ensure all form components are fully ADA compliant.

**Acceptance Criteria**:
- [ ] All inputs have proper labels (visible or aria-label)
- [ ] Error messages associated with inputs (aria-describedby)
- [ ] Required fields clearly marked
- [ ] Keyboard navigation works for all form elements
- [ ] Focus management on form submission
- [ ] Screen reader announcements for errors
- [ ] Color contrast meets WCAG AA standards
- [ ] Form validation accessible to screen readers
- [ ] Accessibility tests pass

**Dependencies**: Tickets 6.1.1, 6.3.1, 6.3.2, 6.3.3

**Implementation Notes**:
- Use semantic HTML form elements
- Associate labels with inputs
- Provide error messages via aria-describedby
- Ensure focus moves to first error on submit
- Test with screen readers (NVDA, JAWS, VoiceOver)

**Files to Update**:
- All form components in `frontend/react/src/components/`
- Add accessibility attributes and tests

---

#### Ticket 6.6.2: ADA Compliance - Navigation & Layout
**Description**: Ensure navigation and layout components are ADA compliant.

**Acceptance Criteria**:
- [ ] Skip navigation link implemented
- [ ] Navigation landmarks (nav, main, aside)
- [ ] Heading hierarchy is logical (h1, h2, h3)
- [ ] Focus indicators visible and clear
- [ ] Keyboard navigation works throughout
- [ ] Focus trap in modals/dialogs
- [ ] ARIA landmarks properly used
- [ ] Screen reader navigation works
- [ ] Accessibility tests pass

**Dependencies**: Tickets 6.1.1, 6.1.4, 6.4.1

**Implementation Notes**:
- Implement skip navigation at top of page
- Use semantic HTML5 elements
- Ensure proper heading hierarchy
- Add ARIA landmarks where needed
- Test keyboard-only navigation

**Files to Update**:
- `frontend/react/src/components/layout/`
- `frontend/react/src/components/navigation/`
- All page components

---

#### Ticket 6.6.3: ADA Compliance - Interactive Components
**Description**: Ensure all interactive components (buttons, links, modals) are ADA compliant.

**Acceptance Criteria**:
- [ ] All buttons have accessible names
- [ ] Icon-only buttons have aria-label
- [ ] Links have descriptive text
- [ ] Modals/dialogs properly announced
- [ ] Modal focus management (trap, return)
- [ ] Loading states announced to screen readers
- [ ] Status messages announced (aria-live)
- [ ] Keyboard shortcuts documented
- [ ] Accessibility tests pass

**Dependencies**: Tickets 6.1.1, All Phase 6.3-6.5 tickets

**Implementation Notes**:
- Use aria-label for icon buttons
- Implement focus trap in modals
- Return focus to trigger after modal closes
- Use aria-live regions for dynamic content
- Provide keyboard shortcuts where appropriate

**Files to Update**:
- All button components
- All modal/dialog components
- All interactive components

---

#### Ticket 6.6.4: ADA Compliance - Data Tables & Lists
**Description**: Ensure data tables and lists are accessible.

**Acceptance Criteria**:
- [ ] Tables have proper headers (th elements)
- [ ] Table headers associated with cells (scope)
- [ ] Table captions where appropriate
- [ ] Lists use proper semantic HTML (ul, ol, dl)
- [ ] List items have accessible names
- [ ] Complex tables have row/column headers
- [ ] Screen reader navigation works
- [ ] Accessibility tests pass

**Dependencies**: Tickets 6.4.2, 6.5.1

**Implementation Notes**:
- Use semantic table elements
- Associate headers with data cells
- Use scope attribute for simple tables
- Use headers/id for complex tables
- Test with screen readers

**Files to Update**:
- `frontend/react/src/components/organization/MemberList.tsx`
- `frontend/react/src/components/meeting/MeetingList.tsx`
- All table components

---

#### Ticket 6.6.5: ADA Compliance - Media & Images
**Description**: Ensure images and media are accessible.

**Acceptance Criteria**:
- [ ] All images have alt text
- [ ] Decorative images have empty alt
- [ ] Complex images have long descriptions
- [ ] Charts/graphs have text alternatives
- [ ] Video/audio have captions/transcripts
- [ ] Color not sole means of conveying information
- [ ] Accessibility tests pass

**Dependencies**: All Phase 6 tickets

**Implementation Notes**:
- Provide meaningful alt text
- Use empty alt for decorative images
- Provide longdesc or aria-describedby for complex images
- Ensure color contrast meets WCAG standards

**Files to Update**:
- All image components
- All chart/graph components (if any)

---

#### Ticket 6.6.6: ADA Compliance - Testing & Validation
**Description**: Set up comprehensive accessibility testing and validation.

**Acceptance Criteria**:
- [ ] Automated accessibility testing (axe-core)
- [ ] Manual testing checklist
- [ ] Screen reader testing (NVDA, JAWS, VoiceOver)
- [ ] Keyboard-only navigation testing
- [ ] Color contrast validation
- [ ] WCAG 2.1 AA compliance verified
- [ ] Accessibility audit report
- [ ] Continuous accessibility checks in CI/CD

**Dependencies**: All Phase 6 tickets

**Implementation Notes**:
- Integrate axe-core into test suite
- Run accessibility tests in CI/CD
- Document manual testing procedures
- Create accessibility testing checklist
- Generate accessibility reports

**Files to Create**:
- `frontend/react/src/__tests__/accessibility.test.tsx`
- `frontend/react/docs/accessibility-testing.md`
- `.github/workflows/accessibility.yml`

---

## Phase 7: Infrastructure & Deployment

**Goal**: Set up deployment infrastructure and CI/CD.

### Phase 7.1: Docker Infrastructure

#### Ticket 7.1.1: Backend Dockerfile
**Description**: Create Dockerfile for Go backend.

**Acceptance Criteria**:
- [ ] Multi-stage Dockerfile
- [ ] Go build stage
- [ ] Minimal runtime image
- [ ] Health check
- [ ] Non-root user
- [ ] Environment variable support
- [ ] Build optimization

**Dependencies**: Ticket 1.1.1

**Implementation Notes**:
- Use multi-stage build
- Minimize image size
- Use alpine base image
- Set up health checks

**Files to Create**:
- `backend/go/Dockerfile`
- `backend/go/.dockerignore`

---

#### Ticket 7.1.2: Frontend Dockerfile
**Description**: Create Dockerfile for React frontend.

**Acceptance Criteria**:
- [ ] Multi-stage Dockerfile
- [ ] Build stage (npm/yarn)
- [ ] Nginx runtime
- [ ] Static file serving
- [ ] SPA routing support
- [ ] Environment variable injection

**Dependencies**: Ticket 1.1.2

**Implementation Notes**:
- Build React app
- Serve with Nginx
- Support SPA routing
- Inject environment variables at build time

**Files to Create**:
- `frontend/react/Dockerfile`
- `frontend/react/nginx.conf`
- `frontend/react/.dockerignore`

---

#### Ticket 7.1.3: Docker Compose Production
**Description**: Create production Docker Compose configuration.

**Acceptance Criteria**:
- [ ] Production docker-compose.yml
- [ ] Service definitions
- [ ] Network configuration
- [ ] Volume definitions
- [ ] Environment variable files
- [ ] Health checks
- [ ] Restart policies

**Dependencies**: Tickets 7.1.1, 7.1.2, 1.1.3

**Implementation Notes**:
- Configure for production
- Set up proper networking
- Use named volumes
- Configure health checks

**Files to Create**:
- `infrastructure/docker/docker-compose.prod.yml`
- `infrastructure/docker/.env.production.example`

---

### Phase 7.2: AWS Infrastructure

#### Ticket 7.2.1: AWS Lambda Configuration
**Description**: Set up AWS Lambda configuration for backend.

**Acceptance Criteria**:
- [ ] Lambda function configuration
- [ ] API Gateway integration
- [ ] Environment variables
- [ ] IAM roles
- [ ] VPC configuration (if needed)
- [ ] Deployment package

**Dependencies**: Ticket 7.1.1

**Implementation Notes**:
- Use AWS SAM or CDK
- Configure API Gateway
- Set up proper IAM permissions
- Configure VPC if accessing RDS

**Files to Create**:
- `infrastructure/aws/lambda/backend-sam.yaml` (or CDK equivalent)
- `infrastructure/aws/lambda/README.md`

---

#### Ticket 7.2.2: RDS Configuration
**Description**: Set up AWS RDS PostgreSQL configuration.

**Acceptance Criteria**:
- [ ] RDS instance configuration
- [ ] Database credentials (Secrets Manager)
- [ ] Backup configuration
- [ ] Security groups
- [ ] Subnet groups
- [ ] Parameter groups

**Dependencies**: None

**Implementation Notes**:
- Use RDS PostgreSQL
- Store credentials in Secrets Manager
- Configure automated backups
- Set up security groups

**Files to Create**:
- `infrastructure/aws/rds/rds.tf` (or CDK/SAM equivalent)
- `infrastructure/aws/rds/README.md`

---

#### Ticket 7.2.3: ElastiCache Configuration
**Description**: Set up AWS ElastiCache for Valkey/Redis.

**Acceptance Criteria**:
- [ ] ElastiCache cluster configuration
- [ ] Security groups
- [ ] Subnet groups
- [ ] Parameter groups
- [ ] Backup configuration

**Dependencies**: None

**Implementation Notes**:
- Use ElastiCache for Redis
- Configure security groups
- Set up backups

**Files to Create**:
- `infrastructure/aws/elasticache/elasticache.tf` (or CDK/SAM equivalent)
- `infrastructure/aws/elasticache/README.md`

---

#### Ticket 7.2.4: S3 Configuration
**Description**: Set up AWS S3 buckets for file storage.

**Acceptance Criteria**:
- [ ] S3 bucket configuration
- [ ] Bucket policies
- [ ] CORS configuration
- [ ] Lifecycle policies
- [ ] Versioning (if needed)

**Dependencies**: None

**Implementation Notes**:
- Create S3 buckets
- Configure access policies
- Set up CORS for frontend
- Configure lifecycle rules

**Files to Create**:
- `infrastructure/aws/s3/s3.tf` (or CDK/SAM equivalent)
- `infrastructure/aws/s3/README.md`

---

#### Ticket 7.2.5: CloudFront Configuration
**Description**: Set up CloudFront for frontend distribution.

**Acceptance Criteria**:
- [ ] CloudFront distribution
- [ ] S3 origin
- [ ] Cache behaviors
- [ ] SSL certificate
- [ ] Custom domain (optional)

**Dependencies**: Ticket 7.2.4

**Implementation Notes**:
- Distribute frontend from S3
- Configure caching
- Set up SSL
- Support custom domain

**Files to Create**:
- `infrastructure/aws/cloudfront/cloudfront.tf` (or CDK/SAM equivalent)
- `infrastructure/aws/cloudfront/README.md`

---

### Phase 7.3: CI/CD Pipeline

#### Ticket 7.3.1: GitHub Actions Workflow
**Description**: Set up CI/CD pipeline with GitHub Actions.

**Acceptance Criteria**:
- [ ] Test workflow (on PR)
- [ ] Build workflow
- [ ] Deploy workflow (on main)
- [ ] Backend tests
- [ ] Frontend tests
- [ ] Linting
- [ ] Security scanning

**Dependencies**: All previous phases

**Implementation Notes**:
- Run tests on every PR
- Build and deploy on main branch
- Run linting and security scans
- Support manual deployment

**Files to Create**:
- `.github/workflows/test.yml`
- `.github/workflows/deploy.yml`

---

#### Ticket 7.3.2: Deployment Automation
**Description**: Automate deployment to AWS.

**Acceptance Criteria**:
- [ ] Backend deployment script
- [ ] Frontend deployment script
- [ ] Database migration automation
- [ ] Rollback capability
- [ ] Deployment notifications

**Dependencies**: Tickets 7.2.1, 7.2.2, 7.3.1

**Implementation Notes**:
- Deploy Lambda functions
- Deploy frontend to S3/CloudFront
- Run migrations automatically
- Support rollback

**Files to Create**:
- `scripts/deploy-backend.sh`
- `scripts/deploy-frontend.sh`
- `scripts/run-migrations.sh`

---

### Phase 7.4: Monitoring & Observability

#### Ticket 7.4.1: Logging Configuration
**Description**: Set up centralized logging.

**Acceptance Criteria**:
- [ ] CloudWatch Logs integration
- [ ] Log aggregation
- [ ] Log retention policies
- [ ] Structured logging
- [ ] Log search capabilities

**Dependencies**: Tickets 2.1.2, 7.2.1

**Implementation Notes**:
- Send logs to CloudWatch
- Use structured logging
- Set retention policies
- Enable log search

**Files to Create**:
- `infrastructure/aws/logging/cloudwatch.tf` (or CDK/SAM equivalent)

---

#### Ticket 7.4.2: Metrics & Monitoring
**Description**: Set up metrics and monitoring.

**Acceptance Criteria**:
- [ ] CloudWatch metrics
- [ ] Custom metrics
- [ ] Dashboards
- [ ] Alarms
- [ ] Uptime monitoring

**Dependencies**: Tickets 7.2.1, 7.4.1

**Implementation Notes**:
- Track API metrics
- Monitor error rates
- Set up uptime checks
- Create dashboards

**Files to Create**:
- `infrastructure/aws/monitoring/cloudwatch-metrics.tf`
- `infrastructure/aws/monitoring/dashboards.json`

---

#### Ticket 7.4.3: Alerting Configuration
**Description**: Set up alerting for critical issues.

**Acceptance Criteria**:
- [ ] Error rate alerts
- [ ] Latency alerts
- [ ] Uptime alerts
- [ ] Database alerts
- [ ] SNS/SES integration
- [ ] Alert routing

**Dependencies**: Ticket 7.4.2

**Implementation Notes**:
- Alert on high error rates
- Alert on high latency
- Alert on downtime
- Route alerts via email/SMS

**Files to Create**:
- `infrastructure/aws/alerting/sns-topics.tf`
- `infrastructure/aws/alerting/alarms.tf`

---

## Summary

This implementation plan breaks down the entire application into 7 phases with detailed tickets. Each ticket is designed to be:
- **Focused**: Single responsibility
- **Testable**: Clear acceptance criteria
- **Dependent**: Dependencies clearly stated
- **Actionable**: Implementation notes provided

The phases build upon each other, starting with foundation (data models, infrastructure) and progressing through business logic, API layer, frontend, and finally deployment.

**Key Compliance Features**:
- **Cookie Consent**: Fully auditable cookie consent system with GDPR/CCPA compliance, including consent history tracking and policy versioning
- **ADA Compliance**: Complete WCAG 2.1 AA compliance with screen reader support, keyboard navigation, and comprehensive accessibility testing

**Estimated Timeline** (rough estimates):
- Phase 1: 2-3 weeks (includes Cookie Consent model)
- Phase 2: 3-4 weeks (includes Consent Repository)
- Phase 3: 2-3 weeks (includes Consent Service)
- Phase 4: 3-4 weeks
- Phase 5: 2-3 weeks (includes Consent API endpoints)
- Phase 6: 4-5 weeks (includes Cookie Consent UI and ADA compliance)
- Phase 7: 2-3 weeks

**Total**: ~18-26 weeks for complete implementation (includes Cookie Consent and ADA compliance)

**Next Steps**:
1. Review and prioritize phases
2. Assign tickets to team members
3. Set up project management tool (Jira, GitHub Projects, etc.)
4. Begin Phase 1 implementation
