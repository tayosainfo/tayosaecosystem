# SACCO Target Architecture And Workflow

This document reflects the active direction after removing `auth-service`:
- InsForge handles authentication/session/OTP.
- Platform services handle onboarding, SACCO membership, and Kibiina flows.
- API Gateway is the enforced entry point for clients.
- Uganda geo hierarchy uses `uganda_geo_data_2025-11-26.csv` as onboarding reference data.

## High-Level Architecture

```mermaid
flowchart LR
    U[Customer] --> W[Web App\nReact + Vite]
    U --> M[Mobile App\nFlutter app/mobile_app]

    W --> IA[InsForge Auth\nSession + OTP]
    M --> IA

    W --> GW[api-gateway-service]
    M --> GW

    GW --> US[user-service]
    GW --> KB[kibiina-service]
    GW --> AFF[affiliate-service]
    GW --> FEE[fee-service]
    GW --> LOAN[loan-credit-service]
    GW --> NOTIF[notification-service]
    GW --> AUDIT[audit-log-service]
    GW --> OBJ[object-storage-service]

    subgraph Infra
      DB[(PostgreSQL)]
    end

    US --> DB
    KB --> DB
    AFF --> DB
    FEE --> DB
    LOAN --> DB
    NOTIF --> DB
    AUDIT --> DB
```

## Registration And Onboarding Workflow

```mermaid
sequenceDiagram
    autonumber
    participant C as Customer
    participant FE as Web/Mobile Client
    participant IA as InsForge Auth
    participant GW as API Gateway
    participant US as User Service + Onboarding
    participant AFF as Affiliate Service
    participant KB as Kibiina Service
    participant OBJ as Object Storage
    participant NOTIF as Notification Service

    C->>FE: Phase 1 quick signup data
    FE->>IA: signUp / signIn / OTP verify
    IA-->>FE: Session token

    FE->>GW: Start onboarding session
    GW->>US: Persist Phase 1 profile baseline
    US-->>GW: Profile created
    GW-->>FE: Continue to KYC

    C->>FE: Submit Phase 2 KYC data
    FE->>GW: Upload KYC payload + files
    GW->>OBJ: Store ID docs + selfie
    GW->>US: Save KYC metadata + compliance declarations
    US->>NOTIF: KYC submission confirmation
    NOTIF-->>C: SMS/Email confirmation

    C->>FE: Submit Phase 3 SACCO membership setup
    FE->>GW: Membership details
    GW->>US: Save membership setup
    US-->>GW: Membership status
    GW-->>FE: Membership ready

    C->>FE: Enter referral code (optional)
    FE->>GW: referralCode + new member id
    GW->>AFF: Validate referral and create pending reward
    AFF->>NOTIF: Send referral acknowledgement
    NOTIF-->>C: Referral received message

    C->>FE: Create or join Kibiina (Phase 4)
    FE->>GW: Parish + village + kibiina preferences
    GW->>US: Validate Parish/Village from geo dataset
    GW->>KB: Create/join village-level merry-go-round group
    KB->>US: Initialize trust score seed + group metadata
    GW-->>FE: Kibiina profile confirmed
```

## Geographic And Group Model Rules

- Parish is the SACCO boundary unit.
- A single parish-level SACCO can contain unlimited village-level Kibiina groups.
- New villages can be added over time by ingesting updates to `uganda_geo_data_2025-11-26.csv`.
- Onboarding address selectors should cascade: `District -> County -> Sub-County -> Parish -> Village`.
- Kibiina creation must require both parish and village IDs to enforce jurisdiction and reporting.

