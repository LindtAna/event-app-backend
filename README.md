# EventApp Backend API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.10.0-008080?style=flat&logo=gin&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-modernc.org%2Fsqlite-003B57?style=flat&logo=sqlite&logoColor=white)
![golang-jwt](https://img.shields.io/badge/golang--jwt-v5.2.1-F28D35?style=flat)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI_3.0-85EA2D?style=flat&logo=swagger&logoColor=black)
![godotenv](https://img.shields.io/badge/godotenv-v1.5.1-222222?style=flat)

## Über das Projekt
Das EventApp Backend ist eine RESTful API, die mit Go (Golang) und dem Gin-Framework entwickelt wurde.
Sie dient als Grundlage für eine Event-Management-Plattform mit Mandantenfähigkeit (Multi-Tenancy).
Benutzer, Events, Kategorien und Teilnahmen sind jeweils einem Tenant (Unternehmen/Organisation) zugeordnet; geschützte Routen und Datenbankabfragen arbeiten streng tenantbezogen.

Benutzer können sich als Tenant-Owner (mit Firmennamen) oder als Member (per Einladungscode) registrieren, ihr Profil bearbeiten, Events im eigenen Tenant verwalten sowie Teilnehmer koordinieren.

---

## Features
- **Mandantenfähigkeit (Multi-Tenancy):** Vollständige Isolierung von Daten pro Tenant (`tenants`-Struktur).
- **Erweiterte Benutzerverwaltung:** Registrierung als `owner` (erstellt Tenant + Admin) oder `member` (über Einladungscode). Aktualisierung des Benutzerprofils.
- **Token-basierte Authentifizierung:** Kurzlebiger Access Token (15 Min.) und Refresh Token (7 Tage) in einem HttpOnly-Cookie (mit adaptiven SameSite/Secure-Attributen).
- **Event-Management:** CRUD-Operationen für Events innerhalb des eigenen Tenants; Filterung nach Ersteller-ID.
- **Teilnehmer-Management:** Benutzer können sich für Events anmelden und wieder abmelden; Abfrage von Teilnehmern pro Event mit tenantweiter Eindeutigkeitsprüfung.
- Einladungsverwaltung: Tenant-Owner können Einladungen erstellen, tenantbezogen auflisten und widerrufen. Mitglieder registrieren sich mit einem gültigen Einladungscode.
- **Zugriffskontrolle & Rollen:** Tenantbezogene Rechte (Owner vs. Member). Nur der Ersteller oder Tenant-Admin darf Events verwalten.
- **API-Dokumentation:** Automatisch generierte und interaktive Swagger-UI.
- **Striktes CORS-Management:** Middleware zur sicheren Verarbeitung von Cross-Origin-Anfragen basierend auf erlaubten `CLIENT_ORIGIN`-Werten.

---

## Technologie-Stack
- **Sprache:** Go (Golang)
- **Web-Framework:** [Gin](https://gin-gonic.com/)
- **Datenbank:** SQLite (via [`modernc.org/sqlite`](https://gitlab.com/cznic/sqlite))–  über `database/sql` und den im Projekt konfigurierten SQLite-Treiber
- **Authentifizierung:** JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt))
- **Umgebungsvariablen:** [`joho/godotenv`](https://github.com/joho/godotenv)
- **API-Dokumentation:** Swagger / OpenAPI ([`swaggo/gin-swagger`](https://github.com/swaggo/gin-swagger))

---

## Repository-Struktur
```
.
├─ cmd
│  ├── api
│  │   ├── main.go          # Einstiegspunkt der Applikation
│  │   ├── server.go        # HTTP-Server Konfiguration
│  │   ├── routes.go        # Definition der API-Endpunkte
│  │   ├── middleware.go    # JWT-Authentifizierung, Tenant-Prüfung und CORS
│  │   ├── auth.go          # Registrierung (Owner/Member), Login, Tokens und Profil
│  │   ├── events.go        # Event- und Teilnehmer-Handler
│  │   ├── invitations.go        # Owner-only HTTP-Handler für Einladungen
│  │   ├── categories.go    # HTTP-Handler zum Abrufen der Tenant-Kategorien
│  │   └── context.go       # Helper zum Extrahieren des authentifizierten Benutzers/Tenants
│  └── migrate
│       └── migrations
│           ├── 000001_create_tenants_table.down.sql
│           ├── 000001_create_tenants_table.up.sql
│           ├── 000002_create_users_table.down.sql
│           ├── ...
│           └── 000006_create_invitations_table.up.sql
│ 
├── internal
│   ├── database          # Datenbankverbindung und Models
│   │   ├── tenants.go           # Tenant-Modell
│   │   ├── attendees.go         # Datenbankzugriffe für Event-Teilnehmer (Tenant-geprüft)
│   │   ├── events.go            # Event-Modell, Validierung, CRUD-Operationen
│   │   ├── models.go            # Zusammenfassung der Datenbankmodelle
│   │   ├── categories.go        # Tenantbezogene Kategorie-Modelle
│   │   ├── event_categories.go  # Transaktionsbasierte Kategorie-Logik
│   │   └── invitations.go       # Einladungen-Modell
│   │   └── users.go             # User-Modell
│   └── env                      # Hilfsfunktionen für Umgebungsvariablen
│       └── env.go
├── docs                  # Automatisch generierte Swagger-Dateien
│   └── testing
│        └── api-manual-tests.md  # Anleitung für manuelle API-Tests & cURL-Beispiele
│
├── go.mod    # Go Module Abhängigkeiten
├── go.sum
└── README.md
```
---

## Umgebungsvariablen & Lokaler Start
#### `/EventApp/.env`

```env
PORT=8080
JWT_SECRET=dein-sehr-geheimes-jwt-secret-mit-mindestens-32-bytes
CLIENT_ORIGIN=http://localhost:5173
```
- `PORT` legt den HTTP-Port fest
- `JWT_SECRET` muss mindestens 32 Bytes lang sein (die Anwendung bricht andernfalls beim Start ab)
- `CLIENT_ORIGIN` definiert die erlaubte(n) Frontend-Origin(s)

---
## Ausführen der Migrationen
Vor dem Start müssen die Datenbankmigrationen ausgeführt werden. Die SQLite-Datei data.db wird dabei im Root-Verzeichnis angelegt

```
go run ./cmd/migrate up
go run ./cmd/api

```

Migrationen lassen sich bei Bedarf mit `go run ./cmd/migrate down` zurücksetzen

---

## API-Dokumentation

Die API ist vollständig mit Swagger dokumentiert. Sobald der Server läuft, kann die interaktive Dokumentation im Browser aufgerufen werden unter:

http://localhost:8080/swagger/index.html

Die im Router registrierten API-Routen verwenden das Präfix `/api/v1`.

### API-Endpunkte

#### Authentifizierung und Benutzerprofil

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | Öffentlich | Tenant-Owner oder Tenant-Member registrieren|
| `POST` | `/api/v1/auth/login` | Öffentlich | Anmelden; Access Token zurückgeben, Refresh-Cookie setzen |
| `POST` | `/api/v1/auth/refresh` | Refresh-Cookie | Neuen Access Token ausstellen |
| `POST` | `/api/v1/auth/logout` | Öffentlich | Refresh-Cookie löschen |
| `PUT` | `/api/v1/users/profile` | Bearer Access Token | Name, Bio und Avatar-URL des eigenen Profils aktualisieren |

#### Registrierungs-Modi (registrationType):
- `owner`: Erstellt einen neuen Tenant und den Administrator. Erfordert companyName, email, password, name
- `member`: Tritt einem bestehenden Tenant bei. Übermittelt den erhaltenen `invitationCode` im Feld invitationCode. Die Einladung muss gültig sein, darf nicht abgelaufen, verwendet oder widerrufen sein, und der zugehörige Tenant muss aktiv sein. Das Anlegen des Mitglieds und das Markieren der Einladung als verwendet erfolgen atomar in einer Transaktion. Tenant-ID, Rolle (member) und Status werden serverseitig festgelegt.

#### Einladungen (Tenant-Owner)

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `POST` | `/api/v1/invitations` | Bearer Access Token, Rolle owner | Erstellt eine Einladung für den eigenen Tenant|
| `GET` | `/api/v1/invitations` | Bearer Access Token, Rolle owner | owner	Listet Einladungen des eigenen Tenants aufn |
| `DELETE` | `/api/v1/auth/refresh` | Bearer Access Token, Rolle owner| Widerruft eine aktive Einladung des eigenen Tenants |

`POST /api/v1/invitations` erwartet ein JSON-Objekt mit expiresAt. Der Wert muss ein gültiger RFC3339-Zeitstempel sein und in der Zukunft liegen. Die Antwort enthält das Einladungsobjekt und invitationCode, das rohe Token, das nur bei der Erstellung zurückgegeben wird. Gespeichert wird in der Datenbank ausschließlich sein SHA-256-Hash; der ursprüngliche Einladungscode kann danach nicht erneut abgerufen werden. Der Code muss deshalb bei der Erstellung an das einzuladende Mitglied weitergegeben werden.

Die Verwaltungsendpunkte für Einladungen sind durch den Bearer Access Token und zusätzlich durch die Rolle owner geschützt. Mitglieder erhalten bei diesen Endpunkten 403 Forbidden. Der Tenant wird ausschließlich aus dem authentifizierten Benutzerkontext abgeleitet; ein tenantId im Request ist nicht erforderlich und wird nicht vertraut.

#### Registrierungs-Modi (registrationType):
- `owner`: Erstellt einen neuen Tenant und den Administrator. Erfordert companyName, email, password, name
- `member`: Tritt einem bestehenden Tenant bei. Erfordert invitationCode, email, password, name
- Passwörter müssen zwischen 8 und 72 Zeichen lang sein. E-Mail-Adressen und Firmennamen müssen eindeutig sein

#### Events(Tenant-geschützt)

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/events` | Bearer Access Token | Alle Events des eigenen Tenants abrufen; optionaler Filter ?ownerId=<id> |
| `GET` | `/api/v1/events/:id` | Bearer Access Token | Einzelnes Tenant-Event abrufen |
| `POST` | `/api/v1/events` | Bearer Access Token | Event im eigenen Tenant erstellen |
| `PUT` | `/api/v1/events/:id` | Bearer Access Token | Eigenes Event aktualisieren |
| `DELETE` | `/api/v1/events/:id` | Bearer Access Token | Eigenes Event löschen; erfolgreiche Antwort ist `204 No Content` |

Der Parameter `ownerId` muss eine gültige Ganzzahl sein; andernfalls antwortet die API mit `400 Bad Request`. Beim Aktualisieren und Löschen muss der angemeldete Benutzer der Ersteller des Events sein.

Validierungen im Event-Modell: title (3–100 Zeichen), description (3–400), location (3–200), Start- und Endzeitpunkt sowie eine Kategorie.

#### Kategorien(Tenant-geschützt)

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/categories` | Bearer Access Token | Alle System- und benutzerdefinierten Kategorien des Tenants abrufen |

Jeder Tenant verfügt über Standardkategorien (Meeting, Workshop, Team-Event, Planning, UI&UX, Konferenz, Weiterbildung, Product Launch) sowie optional benutzerdefinierte Kategorien, die dynamisch beim Erstellen/Aktualisieren von Events verwaltet und bei Nicht-mehr-Verwendung bereinigt werden.

#### Teilnehmer(Tenant-geschützt)

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/events/:id/attendees` | Bearer Access Token | Teilnehmer eines Tenant-Events abrufen |
| `GET` | `/api/v1/attendees/:id/events` | Bearer Access Token | Events abrufen, für die der Benutzer als Teilnehmer eingetragen ist |
| `POST` | `/api/v1/events/:id/attendees/:userId` | Bearer Access Token | Teilnehmer eintragen |
| `DELETE` | `/api/v1/events/:id/attendees/:userId` | Bearer Access Token | Teilnehmer austragen |

Beim Eintragen müssen das Event und der Zielbenutzer existieren. Bereits vorhandene Teilnahmen werden mit `409 Conflict` abgewiesen. Ein Benutzer darf sich selbst eintragen.

Die Teilnehmerliste liefert Benutzerinformationen aus den Feldern `id` und `name`.


### Manuelle API-Tests und cURL-Anleitung

Das Projekt enthält eine detaillierte Dokumentation für manuelle Tests unter **`docs/testing/api-manual-tests.md`**. 
Diese Anleitung deckt Test-Szenarien und verifizierte cURL-Befehle (für Windows `cmd`) ab, darunter:

- Tenant-Isolierung und Eindeutigkeitsprüfungen (Company Name, E-Mail)
- Event- und Kategorienverwaltung sowie Schutz vor ID-Spoofing
- Teilnehmer-Management (Attendees) und Berechtigungsprüfungen
- Authentifizierung, Refresh-Token, Logout und CORS-Verhalten
- Einladungs-Workflow (Erstellung, Nutzung durch Member, Widerruf)
- SQLite-Integritätsprüfungen (Foreign Keys und CASCADE-Löschung)

---

## Datenbank-Schema

Die SQLite-Datenbank wird durch die Migrationen `000001` bis `000006` strukturiert. Die Migrationen definieren die Tabellen `tenants`,`users`,`invitations`, `events`, `attendees` und `categories` einschließlich der unten aufgeführten Datentypen, Standardwerte und Constraints.

### tenants
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`)|
| `name` | `TEXT` | `NOT NULL`; Name des Unternehmens / der Organisation|
| `name_key` | `TEXT` | `NOT NULL`, `UNIQUE`; Normalisierter Name für eindeutige Prüfung|
| `slug` | `TEXT` | `NOT NULL`, `UNIQUE`; URL-freundlicher Bezeichner|
| `status` | `TEXT` | `NOT NULL`, Standardwert 'active', CHECK (status IN ('active', 'suspended'))|
| `created_at` | `TEXT` | `NOT NULL`, Standardwert 'CURRENT_TIMESTAMP'|

### users
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `tenant_id` | `INTEGER` | `NOT NULL`; Fremdschlüssel auf tenants(id) mit ON `DELETE CASCADE`|
| `email` | `TEXT` | `NOT NULL`, `UNIQUE`, `COLLATE NOCASE`; E-Mail-Adresse für den Login|
| `name` | `TEXT` | `NOT NULL`; Anzeigename |
| `password_hash` | `TEXT` | `NOT NULL`; enthält den bcrypt-Hash des Passworts |
| `bio` | `TEXT` | Standardwert `''`; Profilbeschreibung(kein `NOT NULL`-Constraint)|
| `avatar_url` | `TEXT` | Standardwert `''`; Avatar-URL(kein `NOT NULL`-Constraint)|
| `role` | `TEXT` |  `NOT NULL`; Standardwert 'member', `CHECK` (role `IN` ('owner', 'member'))|
| `status` | `TEXT` |  `NOT NULL`; Standardwert 'active', `CHECK` (status `IN` ('active', 'disabled'))|
| `created_at` | `TEXT` |  `NOT NULL`; Standardwert `CURRENT_TIMESTAMP`|

Im Go-Modell ist das Passwort mit json:"-" markiert und wird deshalb nicht in JSON-Antworten serialisiert. Die Kombination aus (tenant_id, id) ist über einen Unique-Constraint abgesichert

### invitations
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `tenant_id` | `INTEGER` | `NOT NULL`; Tenant-Zuordnung|
| `token_hash` | `TEXT` | `NOT NULL`, `UNIQUE`; SHA-256-Hash des Einladungstokens; der Klartext-Token wird nicht gespeichert |
| `expires_at` | `TEXT` | `NOT NULL`; Ablaufzeitpunkt der Einladung |
| `used_at` | `TEXT` | Optional; Zeitpunkt der Verwendung |
| `revoked_at` | `TEXT` | Optional; Zeitpunkt des Widerrufs |
| `created_by` | `INTEGER` |  `NOT NULL`; ID des Erstellers (Verweis auf users)|
| `created_at` | `TEXT` |  `NOT NULL`; Standardwert `CURRENT_TIMESTAMP`|

### events
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `tenant_id` | `INTEGER` | `NOT NULL`; Tenant-Zuordnung |
| `owner_id` | `INTEGER` | `NOT NULL`; Ersteller des Events (Verweis auf users)|
| `title` | `TEXT` | `NOT NULL`; Titel des Events (3–100 Zeichen)|
| `description` | `TEXT` | `NOT NULL`; Beschreibung (3–400 Zeichen)|
| `image_url` | `TEXT` | Standardwert `''`; URL des Event-Bildes |
| `location` | `TEXT` | `NOT NULL`; Veranstaltungsort (3–200 Zeichen)|
| `start_date_time` | `TEXT` | `NOT NULL`; Beginn des Events |
| `end_date_time` | `TEXT` | `NOT NULL`; Ende des Events |
| `category_id` | `TEXT` | `NOT NULL`; ID der zugeordneten Kategorie |
| `url` | `TEXT` | Standardwert `''`; weiterführender Link |

Die Migration deklariert einen Fremdschlüssel von `owner_id` auf `users.id`. 
Die Kategorie wird stattdessen in der Anwendung durch `resolveCategory` validiert bzw. aufgelöst. 

### attendees
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `tenant_id` | `INTEGER` | `NOT NULL`;Tenant-Zuordnung |
| `user_id` | `INTEGER` | `NOT NULL`; verknüpfter Benutzer|
| `event_id` | `INTEGER` | `NOT NULL`; verknüpftes Event |
| `created_at` | `TEXT` | `NOT NULL`; Standardwert `CURRENT_TIMESTAMP` |

Unique-Constraint: UNIQUE (tenant_id, event_id, user_id) verhindert doppelte Anmeldungen desselben Benutzers zum selben Event innerhalb desselben Tenants direkt auf DB-Ebene

### categories
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `tenant_id` | `TEXT` | `NOT NULL`; Fremdschlüssel auf tenants(id) mit `ON DELETE CASCADE`|
| `id` | `TEXT` | Teil des Primärschlüssels; System- oder generierte Hex-ID|
| `name` | `TEXT` | `NOT NULL`, Anzeigename der Kategorie|
| `name_key` | `TEXT` | `NOT NULL`, Eindeutiger Name-Schlüssel pro Tenant|
| `is_default` | `INTEGER` | `NOT NULL`, Standardwert `0`; `CHECK (is_default IN (0, 1))` (Systemkategorien) |

- Primärschlüssel: Zusammengesetzt aus (tenant_id, id).
- Constraints: UNIQUE (tenant_id, name_key) stellt sicher, dass Kategorienamen pro Tenant eindeutig sind.

### Beziehungen und Löschverhalten

- Multi-Tenant Composite Foreign Keys: Um strikte Tenant-Isolierung zu garantieren, verknüpfen zusammengesetzte Fremdschlüssel (wie `EVENTS(tenant_id, owner_id) -> USERS(tenant_id, id)` und `EVENTS(tenant_id, category_id) -> CATEGORIES(tenant_id, id)`) Datensätze plattformübergreifend an denselben Tenant.
- Cascade Deletes:
  - Das Löschen eines Tenants löscht kaskadierend alle zugehörigen Benutzer, Kategorien, Events, Teilnahmen und Einladungen (`ON DELETE CASCADE`)
  - Das Löschen eines Users oder Events löscht abhängige Teilnahmen und Einladungen ebenfalls automatisch
- Category Restrict: events_category_same_tenant_fk verwendet ON DELETE RESTRICT, damit Kategorien, die noch von Events verwendet werden, nicht versehentlich gelöscht werden können.

---

## Architektur

Das Projekt folgt einer klaren Separation of Concerns (SoC), trennt HTTP-Verarbeitung, Authentifizierung und Datenbankzugriffe:

- **Routing & Handling (cmd/api):** Registrieren die Endpunkte, validieren Request-Daten, wenden Zugriffskontrollen an und erzeugen HTTP-Antworten.
- **Middleware (`cmd/api/middleware.go`):** Prüft bei geschützten Routen den Bearer Access Token und stellt `user` sowie `userId` im Gin-Kontext bereit. Die CORS-Middleware behandelt Cross-Origin-Anfragen und OPTIONS-Preflight-Requests.
- **Datenzugriff (`internal/database`):** Kapselt SQL-Abfragen für Benutzer, Events, Teilnehmer und Kategorien.
- Einladungen (`cmd/api/invitations.go`): Stellt Owner-only-Endpunkte zum Erstellen, Auflisten und Widerrufen von Einladungen bereit.
- **Kategorieauflösung (`internal/database/event_categories.go`):** Löst vorhandene Kategorien auf, erstellt bei Bedarf benutzerdefinierte Kategorien und bereinigt ungenutzte Kategorien in Datenbanktransaktionen.
- **Konfiguration (`internal/env`):** Enthält Hilfsfunktionen zum Laden von Konfigurationswerten.


### Ablauf: Authentifizierung

1. Bei `POST /api/v1/auth/register` wird das Passwort mit bcrypt gehasht und der Benutzer gespeichert.
2. Bei `POST /api/v1/auth/login` werden E-Mail und Passwort geprüft. Bei erfolgreicher Anmeldung werden ein Access Token und ein Refresh Token erstellt.
3. Der Access Token wird im JSON-Feld `accessToken` zurückgegeben. Der Refresh Token wird als Cookie namens `refreshToken` gesetzt.
4. Geschützte Endpunkte erwarten `Authorization: Bearer <accessToken>`.
5. `POST /api/v1/auth/refresh` liest das Refresh-Cookie und gibt einen neuen Access Token sowie das Benutzerobjekt zurück.
6. `POST /api/v1/auth/logout` löscht das Refresh-Cookie.

### Ablauf: Event und Kategorie

Beim Erstellen oder Aktualisieren eines Events werden Event-Änderung und Kategorieauflösung in einer Transaktion zusammengefasst. Eine Kategorie wird entweder anhand einer gültigen `categoryId` gefunden oder anhand von `categoryName` wiederverwendet bzw. neu angelegt. Beim Wechsel oder Löschen eines Events wird geprüft, ob die bisherige benutzerdefinierte Kategorie noch von anderen Events verwendet wird.

Ablauf: Einladungen
- Ein Tenant-Owner erstellt über `POST /api/v1/invitations` eine Einladung mit einem Ablaufzeitpunkt in der Zukunft.
- Die API gibt den einmaligen `invitationCode` zurück und speichert ausschließlich dessen SHA-256-Hash.
- Der Owner übermittelt den Einladungscode an das gewünschte Mitglied.
- Das Mitglied registriert sich über `POST /api/v1/auth/register` mit `registrationType: "member" ` und `invitationCode`.
- Die API validiert Einladung und Tenant und erstellt den Benutzer sowie den Verbrauch der Einladung gemeinsam in einer Transaktion.
- Der Owner kann Einladungen über `GET /api/v1/invitations` auflisten und noch gültige, ungenutzte Einladungen über `DELETE /api/v1/invitations/:id` widerrufen.

### System-Architektur

```text
                          Frontend / HTTP-Client
                                   │
┌───────────────────────────────── ▼ ───────────────────────────────────────┐
│            Gin HTTP-Server                                                │
│             - Gin Logger und Recovery                                     │
│            - CORSMiddleware (global; behandelt auch OPTIONS-Preflight)    │
│                                  |                                        │
│  ┌────────────────────────────── ▼ ─────────────────────────────┐         │
│  │ AuthMiddleware (Prüft Bearer Access Token & Tenant-Status)   │         │
│  └──────────────────────────────────────────────────────────────┘         │
│  ┌────────────────────────│───────────────────────────────────┐           │
│  │                        │                                   │           │
│  │  ┌──────────────────── ▼ ─┐        ┌────────────────────── ▼ ──────┐   │
│  │  │ Öffentliche API-Routen │        │ Geschützte API-Routen         │   │
│  │  │ Präfix: /api/v1        │        │ Präfix: /api/protected        │   │
│  │  │                        │        │                               │   │
│  │  │ Auth                   │        │ Benutzerprofil                │   │
│  │  │- POST /auth/register   │        │- PUT /users/profile           │   │
│  │  │(Modus: owner / member) │        │                               │   │
│  │  │- POST /auth/login      │        │Events                         │   │
│  │  │- POST /auth/refresh    │        │- GET /events                  │   │
│  │  │- POST /auth/logout     │        │- GET /events/:id              │   │
│  │  │                        │        │- POST /events                 │   │
│  │  │                        │        │- PUT /events/:id              │   │
│  │  │                        │        │- DELETE /events/:id           │   │
│  │  │                        │        │                               │   │
│  │  │                        │        │Kategorien                     │   │
│  │  │                        │        │- GET /categories              │   │
│  │  │                        │        │                               │   │
│  │  │                        │        │Teilnahmeverwaltung            │   │
│  │  │                        │        │- GET /events/:id/attendees    │   │
│  │  │                        │        │- GET /attendees/:id/events    │   │
│  │  │                        │        │- POST /events/:id/attendees   │   │
│  │  │                        │        │- DELETE /events/:id/attendees │   │
│  │  │                        │        │                               │   │
│  │  │                        │        │Einladungen (Owner)            │   │
│  │  │                        │        │- POST /invitations            │   │
│  │  │                        │        │- GET /invitations             │   │
│  │  │                        │        │- DELETE /invitations/:id      │   │
│  │  └────────────────────────┘        └───────────────────────────────┘   │
│  └─────────────────────────────────────────────────────────────────────┘  │
│                               │                                           │
│              ┌─────────────── ▼ ──────────────────┐                       │
│              │ HTTP-Handler / Application Methods │                       │
│              │           (cmd/api)                │                       │
│              └────────────────────────────────────┘                       │
│                              │                                            │
│ ┌─────────────────────────── ▼ ────────────────────────────────────────┐  │
│ │                Datenbankmodelle & Tenant-Logik                       │  │
│ │                    (internal/database)                               │  │
│ │     ┌───────────────────────────────────────────────────────┐        │  │
│ │     │          │        │         │           │             │        │  │
│ │ ┌───▼────┐ ┌───▼──┐ ┌───▼──┐ ┌────▼────┐ ┌────▼──────┐ ┌────▼─────┐  │  │
│ │ │Tenants │ │ Users│ │Events│ │Attendees│ │Invitations│ │Categories│  │  │
│ │ └────────┘ └──────┘ └──────┘ └─────────┘ └───────────┘ └──────────┘  │  │
│ │                        │                                             │  │
│ │       ┌────────────────▼───────────────────┐                         │  │
│ │       │event_categories.go                 │                         │  │
│ │       │- Tenant-Kategorie auflösen         │                         │  │
│ │       │- Transaktional erstellen/bereinigen│                         │  │
│ │       └────────────────────────────────────┘                         │  │
│ └──────────────────────────────────────────────────────────────────────┘  │
│                               │                                           │
│                 ┌─────────────▼─────────────────┐                         │
│                 │  database/sql + SQLite-Treiber│                         │
│                 │ (modernc.org/sqlite)          │                         │
│                 └───────────────────────────────┘                         │
│                               │                                           │
│┌───────────────────────────── ▼ ────────────────────────────────────────┐ │
││   SQLite-Datenbank: data.db                                            │ │
││                                                                        │ │
││ Tenant-gebundene Tabellen:                                             │ │
││- tenants (id, name, slug, status)                                      │ │
││- users (id, tenant_id, email, name,password, role, status, ...)        │ │
││- events (id, tenant_id, owner_id,title, description, category_id,...)  │ │
││- attendees (id, tenant_id, user_id, event_id, UNIQUE(tenant,event,usr))│ │
││- categories (id, tenant_id, name,is_default)                           │ │
││- invitations (id, tenant_id,token_hash, expires_at, ...)               │ │
│└────────────────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────────────┘
Swagger UI: /swagger/*any (außerhalb des Präfixes /api/v1)
```

Die öffentlichen Routen sind nicht durch `AuthMiddleware` geschützt. Für `POST /api/v1/auth/refresh` wird dennoch das `refreshToken`-Cookie benötigt; `POST /api/v1/auth/logout` löscht dieses Cookie. Die geschützten Routen erwarten einen gültigen Bearer Access Token. `CORSMiddleware` wird global registriert. `gin.Default()` aktiviert zusätzlich die Gin-Logger- und Recovery-Middleware.

---

## Sicherheit und Authentifizierung

- **Access Token:** JWT mit HMAC-Signatur; Ablauf nach 15 Minuten.
- **Refresh Token:** JWT mit HMAC-Signatur; Ablauf nach sieben Tagen im HttpOnly-Cookie (refreshToken). Bei HTTPS mit Secure und SameSite=None, sonst SameSite=Lax.
- Passwortspeicherung: Bcrypt-Hashing (max. 72 Bytes). Passwörter werden niemals in JSON-Antworten serialisiert.
- Einladungen: Verwaltungsendpunkte sind ausschließlich für Tenant-Owner zugänglich. Einladungstoken werden nur bei der Erstellung im Klartext zurückgegeben und in der Datenbank ausschließlich als SHA-256-Hash gespeichert. Einladungen sind einmalig verwendbar und können vor ihrer Verwendung widerrufen werden.
- CORS: Leitet Anfragen nur von konfigurierten CLIENT_ORIGIN-Domains durch und erlaubt Credentials.
Diese Cookie- und CORS-Einstellungen werden für eine produktive HTTPS-Bereitstellung passend  konfiguriert.
---

## Roadmap

[ ] Integration mit dem Frontend (React)

[ ] Hinzufügen von Unit- und Integration-Tests

[ ] Containerisierung des Backends (Docker)