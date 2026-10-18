# EventApp Backend API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.10.0-008080?style=flat&logo=gin&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-modernc.org%2Fsqlite-003B57?style=flat&logo=sqlite&logoColor=white)
![golang-jwt](https://img.shields.io/badge/golang--jwt-v5.2.1-F28D35?style=flat)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI_3.0-85EA2D?style=flat&logo=swagger&logoColor=black)
![godotenv](https://img.shields.io/badge/godotenv-v1.5.1-222222?style=flat)

## Über das Projekt
Das EventApp Backend ist eine RESTful API, die mit Go (Golang) und dem Gin-Framework entwickelt wurde.
Sie dient als Grundlage für eine Event-Management-Plattform.
Benutzer können sich registrieren und anmelden, ihr Profil bearbeiten, Events erstellen und verwalten sowie sich als Teilnehmer für Events ein- und austragen.

Events verfügen über einen Ersteller, Veranstaltungsdaten und eine Kategorie. Neben fest vorgegebenen Systemkategorien unterstützt die API benutzerdefinierte Kategorien. Die API stellt außerdem Endpunkte zum Abrufen der Event-Teilnehmer und der Events eines Teilnehmers bereit.

---

## Features
- **Benutzerverwaltung:** Registrierung, Login und Aktualisierung des Benutzerprofils
- **Token-basierte Authentifizierung:** Kurzlebiger Access Token und Refresh Token in einem HttpOnly-Cookie.
- **Event-Management:** CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen) für Events; öffentliche Filterung nach Ersteller-ID.
- **Teilnehmer-Management:** Benutzer können sich für Events anmelden und wieder abmelden; Abfrage von Teilnehmern pro Event.
- **Zugriffskontrolle:** Nur der Ersteller darf sein Event ändern oder löschen. Benutzer können sich selbst anmelden oder abmelden.
- **API-Dokumentation:** Automatisch generierte und interaktive Swagger-UI.
- **CORS-Support:** Middleware zur Verarbeitung von Cross-Origin-Anfragen.

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
│  │   ├── main.go       # Einstiegspunkt der Applikation
│  │   ├── server.go     # HTTP-Server Konfiguration
│  │   ├── routes.go     # Definition der API-Endpunkte
│  │   ├── middleware.go    # JWT-Authentifizierung und CORS
│  │   ├── auth.go          # Registrierung, Login, Tokens und Profil
│  │   ├── events.go     # Event- und Teilnehmer-Handler
│  │   ├── categories.go # HTTP-Handler zum Abrufen aller Kategorien
│  │   └── context.go    # Helper zum Extrahieren des authentifizierten Benutzers aus dem │  └── migrate
│       └── migrations
│           ├── 000001_create_users_table.up.sql
│           ├── 000004_create_categories_table.up.sql
│           ├── 000004_create_categories_table.up.sql
│           └── ..usw
│ 
├── internal
│   ├── database          # Datenbankverbindung und Model
│   │   ├── attendees.go         # Datenbankzugriffe für Event-Teilnehmer(CRUD & M:N Verknüpfungen)
│   │   ├── events.go            # Event-Modell, Struct-Validierung, CRUD-Datenbankoperationen für Events
│   │   ├── models.go            # Zusammenfassung der Datenbankmodelle
│   │   ├── categories.go        # Category-Modell und Datenbankabfrage aller Kategorien
│   │   ├── event_categories.go  # Transaktionsbasierte Event-/Kategorie-Logik
│   │   └── users.go             # User-Modell und Profilzugriffe
│   └── env                      # Hilfsfunktionen für Umgebungsvariablen
│       └── env.go
├── docs                  # Automatisch generierte Swagger-Dateien
├── go.mod                # Go Module Abhängigkeiten
├── go.sum
└── README.md
```
Die Dateien zur Kategorienfunktion haben unterschiedliche Aufgaben:
- `cmd/api/categories.go` enthält den Gin-HTTP-Handler `getAllCategories`. Er ruft `app.models.Categories.GetAll(...)` auf und gibt die Kategorien als JSON mit Status `200 OK` zurück; bei einem Fehler antwortet er mit `500 Internal Server Error`.

- `internal/database/categories.go` enthält das Datenmodell `Category`, `CategoryModel` und die Datenbankabfrage `GetAll`. Die Kategorien werden nach `is_default` (Systemkategorien zuerst) und anschließend nach Namen sortiert.

- `internal/database/event_categories.go` enthält die separate Geschäfts- und Datenbanklogik zum Auflösen, Erstellen und Bereinigen von Kategorien innerhalb von Transaktionen, wenn Events erstellt, aktualisiert oder gelöscht werden.
---

##  Umgebungsvariablen
#### `/EventApp/.env`

```env
PORT=8080
JWT_SECRET=dein-sehr-geheimes-jwt-secret
CLIENT_ORIGIN=http://localhost:5173
```
- `PORT` legt den HTTP-Port fest
- `JWT_SECRET` wird zum Signieren und Prüfen der JWTs verwendet. Für eine reale Bereitstellung muss ein ausreichend zufälliger, geheimer Wert gesetzt werden
---

## API-Dokumentation

Die API ist vollständig mit Swagger dokumentiert. Sobald der Server läuft, kann die interaktive Dokumentation im Browser aufgerufen werden unter:

http://localhost:8080/swagger/index.html

Die im Router registrierten API-Routen verwenden das Präfix `/api/v1`.

### API-Endpunkte

#### Authentifizierung und Benutzerprofil

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `POST` | `/api/v1/auth/register` | Öffentlich | Benutzerkonto erstellen |
| `POST` | `/api/v1/auth/login` | Öffentlich | Anmelden; Access Token und Benutzerobjekt zurückgeben, Refresh Token als Cookie setzen |
| `POST` | `/api/v1/auth/refresh` | Refresh-Cookie | Neuen Access Token aus dem Refresh Token ausstellen |
| `POST` | `/api/v1/auth/logout` | Öffentlich; Cookie wird gelöscht | Refresh-Cookie entfernen |
| `PUT` | `/api/v1/users/profile` | Bearer Access Token | Name, Bio und Avatar-URL des angemeldeten Benutzers aktualisieren |

Bei der Registrierung werden `email`, `password` und `name` erwartet. Die E-Mail muss gültig sein; das Passwort muss mindestens acht Zeichen lang sein. Bei der Profilaktualisierung ist der Name erforderlich und muss mindestens zwei Zeichen lang sein. `bio` und `avatarUrl` sind optionale Textfelder.

#### Events

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/events` | Öffentlich | Alle Events abrufen; optionaler Filter `?ownerId=<id>` |
| `GET` | `/api/v1/events/:id` | Öffentlich | Einzelnes Event abrufen |
| `POST` | `/api/v1/events` | Bearer Access Token | Event erstellen; der angemeldete Benutzer wird als Ersteller gesetzt |
| `PUT` | `/api/v1/events/:id` | Bearer Access Token | Eigenes Event aktualisieren |
| `DELETE` | `/api/v1/events/:id` | Bearer Access Token | Eigenes Event löschen; erfolgreiche Antwort ist `204 No Content` |

Der Parameter `ownerId` muss eine gültige Ganzzahl sein; andernfalls antwortet die API mit `400 Bad Request`. Beim Aktualisieren und Löschen muss der angemeldete Benutzer der Ersteller des Events sein.

Ein Event verwendet im API-Modell diese Felder:

| JSON-Feld | Bedeutung | Validierung im Event-Modell |
|---|---|---|
| `id` | Event-ID | Wird vom Backend gesetzt |
| `ownerId` | ID des Erstellers | Wird beim Erstellen aus dem authentifizierten Benutzer übernommen |
| `owner` | Objekt mit `id` und `name` des Erstellers | Wird beim Abrufen mitgeladen, sofern ein Ersteller vorhanden ist |
| `title` | Titel des Events | Erforderlich, 3–100 Zeichen |
| `description` | Beschreibung | Erforderlich, 3–400 Zeichen |
| `imageUrl` | Bild-URL | Optional |
| `location` | Veranstaltungsort | Erforderlich, 3–200 Zeichen |
| `startDateTime` | Startdatum und -uhrzeit | Erforderlich |
| `endDateTime` | Enddatum und -uhrzeit | Erforderlich |
| `categoryId` | ID einer bestehenden Kategorie | Wird bei der Kategorieauflösung geprüft |
| `categoryName` | Name einer Kategorie, z. B. für eine neue benutzerdefinierte Kategorie | Optional; darf nicht gleichzeitig mit `categoryId` angegeben werden |
| `url` | Weiterführender Link | Optional |

Für `categoryId` und `categoryName` gilt: Bei Event-Erstellung und -Aktualisierung muss genau eine der beiden Angaben verwendet werden. Eine nicht vorhandene `categoryId` oder eine fehlende/ungültige Kategorie führt zu `400 Bad Request`. Die Datenbanklogik speichert die Kategorie-ID am Event.

#### Kategorien

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/categories` | Öffentlich | Alle System- und benutzerdefinierten Kategorien abrufen |

Die Kategorien werden in der separaten Tabelle `categories` gespeichert. Die Systemkategorien stehen in der Antwort vor den benutzerdefinierten Kategorien; innerhalb der jeweiligen Gruppe wird alphabetisch ohne Beachtung der Groß-/Kleinschreibung sortiert.

Folgende Systemkategorien werden durch die Migration angelegt und nicht automatisch gelöscht:

- Meeting
- Workshop
- Team-Event
- Planning
- UI&UX
- Konferenz
- Weiterbildung
- Product Launch

Beim Erstellen oder Aktualisieren eines Events kann eine vorhandene `categoryId` oder ein `categoryName` angegeben werden. Bei einem neuen Namen legt das Backend eine benutzerdefinierte Kategorie an. Die Namen sind in der Datenbank ohne Beachtung der Groß-/Kleinschreibung eindeutig. Wenn ein Event seine Kategorie wechselt oder gelöscht wird, entfernt das Backend eine benutzerdefinierte Kategorie, sobald kein Event mehr auf sie verweist. Systemkategorien bleiben bestehen.

Die Kategorieauflösung und die zugehörige Event-Änderung werden jeweils innerhalb derselben Datenbanktransaktion ausgeführt.

#### Teilnehmer

| Methode | Endpunkt | Authentifizierung | Beschreibung |
|---|---|---|---|
| `GET` | `/api/v1/events/:id/attendees` | Öffentlich | Teilnehmer eines Events abrufen |
| `GET` | `/api/v1/attendees/:id/events` | Öffentlich | Events abrufen, für die der Benutzer als Teilnehmer eingetragen ist |
| `POST` | `/api/v1/events/:id/attendees/:userId` | Bearer Access Token | Teilnehmer eintragen |
| `DELETE` | `/api/v1/events/:id/attendees/:userId` | Bearer Access Token | Teilnehmer austragen |

Beim Eintragen müssen das Event und der Zielbenutzer existieren. Bereits vorhandene Teilnahmen werden mit `409 Conflict` abgewiesen. Ein Benutzer darf sich selbst eintragen; der Ersteller des Events darf auch einen anderen Benutzer eintragen. Beim Austragen darf der angemeldete Benutzer seine eigene Teilnahme beenden oder als Ersteller einen Teilnehmer entfernen.

Die Teilnehmerliste liefert Benutzerinformationen aus den Feldern `id`, `name` und `email`. Das Passwort wird nicht in JSON-Antworten serialisiert.


---

## Datenbank-Schema

Die SQLite-Datenbank wird durch die Migrationen `000001` bis `000004` strukturiert. Die Migrationen definieren die Tabellen `users`, `events`, `attendees` und `categories` einschließlich der unten aufgeführten Datentypen, Standardwerte und Constraints.

### `users`

| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `email` | `TEXT` | `NOT NULL`, `UNIQUE`; E-Mail-Adresse für den Login |
| `name` | `TEXT` | `NOT NULL`; Anzeigename |
| `password` | `TEXT` | `NOT NULL`; enthält den bcrypt-Hash des Passworts |
| `bio` | `TEXT` | Standardwert `''`; Profilbeschreibung. Die Migration setzt für dieses Feld kein `NOT NULL`. |
| `avatar_url` | `TEXT` | Standardwert `''`; Avatar-URL. Die Migration setzt für dieses Feld kein `NOT NULL`. |

Im Go-Modell ist das Passwort mit `json:"-"` markiert und wird deshalb nicht in JSON-Antworten serialisiert.

### `events`
| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `owner_id` | `INTEGER` | `NOT NULL`; Ersteller des Events; Fremdschlüssel auf `users(id)` mit `ON DELETE CASCADE` |
| `title` | `TEXT` | `NOT NULL`; Titel des Events |
| `description` | `TEXT` | `NOT NULL`; Beschreibung |
| `image_url` | `TEXT` | Standardwert `''`; URL des Event-Bildes |
| `location` | `TEXT` | `NOT NULL`; Veranstaltungsort |
| `start_date_time` | `TEXT` | `NOT NULL`; Beginn des Events |
| `end_date_time` | `TEXT` | `NOT NULL`; Ende des Events |
| `category_id` | `TEXT` | Standardwert `''`; ID der zugeordneten Kategorie |
| `url` | `TEXT` | Standardwert `''`; weiterführender Link |

Die Migration deklariert einen Fremdschlüssel von `owner_id` auf `users.id`. 
Die Kategorie wird stattdessen in der Anwendung durch `resolveCategory` validiert bzw. aufgelöst. 

### `attendees`

Quelle: `000003_create_attendees_table.up.sql`

| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `INTEGER` | Primärschlüssel (`PRIMARY KEY`), automatisch erhöht (`AUTOINCREMENT`) |
| `user_id` | `INTEGER` | `NOT NULL`; Fremdschlüssel auf `users(id)` mit `ON DELETE CASCADE` |
| `event_id` | `INTEGER` | `NOT NULL`; Fremdschlüssel auf `events(id)` mit `ON DELETE CASCADE` |

Die Migration definiert Fremdschlüssel auf `users` und `events`, aber **keinen zusammengesetzten Unique-Constraint** für `(user_id, event_id)`. 
Die API prüft vor dem Eintragen, ob die Teilnahme bereits existiert.

### `categories`

| Feld | Typ | Eigenschaften / Zweck |
|---|---|---|
| `id` | `TEXT` | Primärschlüssel (`PRIMARY KEY`); Systemkategorien verwenden feste IDs, benutzerdefinierte Kategorien zufällig erzeugte hexadezimale IDs |
| `name` | `TEXT` | `NOT NULL`, `UNIQUE`, `COLLATE NOCASE`; Kategorienamen sind ohne Beachtung der Groß-/Kleinschreibung eindeutig |
| `is_default` | `INTEGER` | `NOT NULL`, Standardwert `0`; `CHECK (is_default IN (0, 1))` begrenzt den Wert auf `0` oder `1` und kennzeichnet Systemkategorien |


Die Anwendung erstellt, aktualisiert und löscht Events zusammen mit der zugehörigen Kategorienauflösung innerhalb von Datenbanktransaktionen.
Wenn eine benutzerdefinierte Kategorie von keinem Event mehr verwendet wird, wird sie entfernt; Systemkategorien werden nicht automatisch gelöscht. 

### Beziehungen und Löschverhalten

- `events.owner_id` verweist auf `users.id`; die Migration deklariert `ON DELETE CASCADE`.
- `attendees.user_id` verweist auf `users.id` und `attendees.event_id` auf `events.id`; beide Fremdschlüssel deklarieren `ON DELETE CASCADE`.
- `events.category_id` wird logisch mit `categories.id` verknüpft und durch die Anwendungslogik geprüft, ist in der vorliegenden Migration aber nicht als SQL-Fremdschlüssel definiert
- Die Tabelle `attendees` besitzt einen eigenen `id`-Primärschlüssel. Die Kombination aus Event und Benutzer wird von der API auf doppelte Einträge geprüft.

---

## Architektur

Das Projekt folgt einer klaren Separation of Concerns (SoC), trennt HTTP-Verarbeitung, Authentifizierung und Datenbankzugriffe:
:

- **Routing & Handling (cmd/api):** Registrieren die Endpunkte, validieren Request-Daten, wenden Zugriffskontrollen an und erzeugen HTTP-Antworten.
- **Middleware (`cmd/api/middleware.go`):** Prüft bei geschützten Routen den Bearer Access Token und stellt `user` sowie `userId` im Gin-Kontext bereit. Die CORS-Middleware behandelt Cross-Origin-Anfragen und OPTIONS-Preflight-Requests.
- **Datenzugriff (`internal/database`):** Kapselt SQL-Abfragen für Benutzer, Events, Teilnehmer und Kategorien.
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

### System-Architektur
```
                                      Frontend / HTTP-Client
                                              │
┌──────────────────────────────────────────── ▼ ────────────────────────────────────────┐
│                        Gin HTTP-Server (gin.Default)
                    - Gin Logger und Recovery
    - CORSMiddleware (global; behandelt auch OPTIONS-Preflight)                            │
│  ┌─────────────────────────┐      ┌───────────────────────────┐  ┌──────────────────┐ │
│  │       Routes            │      │ Application Methods       │  │    Middleware    │ │
│  │                         │      │                           │  │                  │ │
│  │ Öffentliche API-Routen: │────▸ │ Auth                      │  │ - CORSMiddleware │ │
│  │ Präfix: /api/v1         │────▸ │ Auth                      │  │ - CORSMiddleware │ │
│  │ Auth                    │      │  - registerUser           │  │   (Global)       │ │
│  │ -POST /auth/register    │      │  - registerUser           │  │   (Global)       │ │
│  │ -POST /auth/login       │      │  - login                  │  │                 
│  │ -POST /auth/refresh     │      │  - login                  │  │                  │ │
│  │ -POST /auth/logout      │      │  - login                  │  │                  │ │
│  │ Events                  │      │                           │  │                  │ │
│  │ -GET /events            │      │ Events (Read)             │  │ - AuthMiddleware │ │
│  │ -GET /events/:id        │      │  - getAllEvents           │  │   (JWT für       │ │
│  │ Kategorien              │      │  - getEvent               │  │   geschützte     │ │
│  │ - GET/categories        │      │  - getEvent               │  │   geschützte     │ │
│  │ Teilnehmer              │      │                           │  │   Routen)        │ │
│  │-GET /events/:id/        │      │                           │  │   Routen)        │ │
│  │attendees                │      │ Events (Write)            │  │                  │ │
│  │-GET /attendees/:id/     │      │  - createEvent            │  │                  │ │
│  │ events                  │      │  - createEvent            │  │                  │ │
│  │                         │      │  - createEvent            │  │                  │ │
│  │ Geschützte API-Routen:  │      │  - updateEvent            │  │                  │ │
│  │  Präfix: /api/v1        │      │  - updateEvent            │  │                  │ │
│  │ AuthMiddleware          │      │  - deleteEvent            │  │                  │ │
│  │ -prüft Authorization:   │      │                           │  │                  │ │
│  │ - DELETE /events/:id    │      │ Attendees                 │  │                  │ │
│  │ - POST /events/:id/     │      │  - getAttendeesForEvent   │  │                  │ │
│  │ attendees/:userId       │      │  - GetAttendeesByEvent    │  │                  │ │
│  │ - DELETE /events/:id/   │      │  - getEventsByAttendee    │  │                  │ │
│  │ attendees/:userId       │      │  - deleteAttendeeFromEvent│  │                  │ │
│  │                         │      │  - addAttendeeToEvent     │  │                  │ │
│  └─────────────────────────┘      └───────────────────────────┘  └──────────────────┘ │
│                                                 │                                     │
│                                                 ▼                                     │
│                       SQLite Driver (database/sql + modernc.org/sqlite)               │
└─────────────────────────────────────────────────┼─────────────────────────────────────┘
                                                  │
                                                  ▼
┌───────────────────────────────────────────────────────────────────────────────────────┐
│                               SQLite Database (data.db)                               │
│                                                                                       │
│  Tabellen:                                                                            │
│  - users (id, name, email, password_hash, created_at)                                 │
│  - events (id, title, description, date, creator_id, created_at, etc.)                │
│  - event_attendees (event_id, user_id) - Join Table für M:N-Beziehung                 │
└───────────────────────────────────────────────────────────────────────────────────────┘
```
### System-Architektur

```text
                Frontend / HTTP-Client
                        │
┌────────────────────── ▼ ────────────────────────────────────────────────────────────────┐
│                    Gin HTTP-Server                                                      │
│                  - Gin Logger und Recovery                                              │
│                 - CORSMiddleware (global; behandelt auch OPTIONS-Preflight)             │
│                                |                                                        │
│  ┌──────────────────────────── ▼ ────────────────────────────────────────────────────┐  │
│  │           ┌─────────────────────────────────────┐                                 │  │
│  │  ┌─────── ▼ ──────────────────┐        ┌─────── ▼ ─────────────────────────────┐  │  │
│  │  │ Öffentliche API-Routen     │        │ Geschützte API-Routen                 │  │  │
│  │  │ Präfix: /api/v1            │        │ Präfix: /api/v1                       │  │  │
│  │  │                            │        │                                       │  │  │
│  │  │ Auth                       │        │ AuthMiddleware                        │  │  │
│  │  │- POST /auth/register       │        │ - prüft  Authorization:               │  │  │
│  │  │- POST /auth/login          │        │  Bearer <accessToken>                 │  │  │
│  │  │- POST /auth/refresh        │        │- prüft JWT und Benutzer               │  │  │
│  │  │- POST /auth/logout         │        │                                       │  │  │
│  │  │                            │        │Benutzerprofil                         │  │  │
│  │  │ Events                     │        │- PUT /users/profile                   │  │  │
│  │  │- GET /events               │        │                                       │  │  │
│  │  │- GET /events/:id           │        │Events                                 │  │  │
│  │  │                            │        │- POST /events                         │  │  │
│  │  │Kategorien                  │        │- PUT /events/:id                      │  │  │
│  │  │- GET /categories           │        │- DELETE /events/:id                   │  │  │
│  │  │                            │        │                                       │  │  │
│  │  │Teilnehmer                  │        │Teilnahmeverwaltung                    │  │  │
│  │  │- GET /events/:id/attendees │        │ - POST /events/:id/  attendees/:userId│  │  │
│  │  │- GET /attendees/:id/events │        │- DELETE /events/:id/attendees/:userId │  │  │
│  │  └────────────────────────────┘        └───────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────────────────────────────────┘  │
│                                       │                                                 │
│                      ┌─────────────── ▼ ──────────────────┐                             │
│                      │ HTTP-Handler / Application Methods │                             │
│                      │           (cmd/api)                │                             │
│                      └────────────────────────────────────┘                             │
│                                      │                                                  │
│               ┌───────────────────── ▼ ──────────────────────┐                          │
│               │          Datenbankmodelle (Models)           │                          │
│               │              (internal/database)             │                          │
│               │    ┌─────────────────────────────────┐       │                          │
│               │    │         │          │            │       │                          │
│               │┌───▼──┐  ┌───▼──┐  ┌────▼────┐  ┌────▼─────┐ │                          │
│               ││ Users│  │Events│  │Attendees│  │Categories│ │                          │
│               │└──────┘  └──────┘  └─────────┘  └──────────┘ │                          │
│               │             │                                │                          │
│               │       ┌─────▼──────────────────────────────┐ │                          │
│               │       │event_categories.go                 │ │                          │
│               │       │- Kategorie auflösen                │ │                          │
│               │       │ - benutzerdefinierte Kategorie     │ │                          │
│               │       │ transaktional erstellen/bereinigen │ │                          │
│               │       └────────────────────────────────────┘ │                          │
│               └──────────────────────────────────────────────┘                          │
│                                       │                                                 │
│                        ┌──────────────▼────────────────┐                                │
│                        │  database/sql + SQLite-Treiber│                                │
│                        │ (modernc.org/sqlite)          │                                │
│                        └───────────────────────────────┘                                │
│                                        │                                                │
│                        ┌───────────────▼────────────────────┐                           │
│                        │   SQLite-Datenbank: data.db        │                           │
│                        │                                    │                           │
│                        │ Tabellen und wesentliche Felder:   │                           │
│                        │ - users (id, email, name, password,│                           │
│                        │   bio, avatar_url)                 │                           │
│                        │ - events (id, owner_id, title,     │                           │
│                        │   description,image_url, location, │                           │
│                        │   start_date_time, end_date_time,  │                           │
│                        │   category_id, url)                │                           │
│                        │ - attendees (id, user_id, event_id)│                           │
│                        │ - categories (id, name, is_default)│                           │
│                        └────────────────────────────────────┘                           │
└─────────────────────────────────────────────────────────────────────────────────────────┘
Swagger UI: /swagger/*any (außerhalb des Präfixes /api/v1)
```

Die öffentlichen Routen sind nicht durch `AuthMiddleware` geschützt. Für `POST /api/v1/auth/refresh` wird dennoch das `refreshToken`-Cookie benötigt; `POST /api/v1/auth/logout` löscht dieses Cookie. Die geschützten Routen erwarten einen gültigen Bearer Access Token. `CORSMiddleware` wird global registriert. `gin.Default()` aktiviert zusätzlich die Gin-Logger- und Recovery-Middleware.

---

## Sicherheit und Authentifizierung

- **Access Token:** JWT mit HMAC-Signatur; Ablauf nach 15 Minuten.
- **Refresh Token:** JWT mit HMAC-Signatur; Ablauf nach sieben Tagen. Das Cookie heißt `refreshToken`, ist `HttpOnly`, wird für den Pfad `/api/v1/auth` und die Domain `localhost` gesetzt und hat derzeit `Secure=false`.
- **Geschützte Routen:** Müssen einen `Authorization`-Header im Format `Bearer <token>` senden. Die Middleware prüft Signatur und Gültigkeit des Tokens sowie, ob der referenzierte Benutzer noch in der Datenbank existiert.
- **Passwortspeicherung:** Passwörter werden mit bcrypt gehasht. Das Passwortfeld wird nicht serialisiert.
- **Autorisierung:** Nur der Event-Ersteller darf sein Event ändern oder löschen. Für das Eintragen und Austragen von Teilnehmern gelten die oben beschriebenen Regeln.
- **CORS-Hinweis:** Die aktuell bereitgestellte Middleware übernimmt einen nichtleeren `Origin`-Header als `Access-Control-Allow-Origin` und setzt `Access-Control-Allow-Credentials: true`. Sie beschränkt Origins nicht auf `CLIENT_ORIGIN`. Außerdem wird das Refresh-Cookie derzeit für die Domain `localhost` und mit `Secure=false` gesetzt. 
Diese Cookie- und CORS-Einstellungen werden für eine produktive HTTPS-Bereitstellung passend zkonfiguriert.
---

## Roadmap

[ ] Integration mit dem Frontend (React)

[ ] Hinzufügen von Unit- und Integration-Tests

[ ] Containerisierung des Backends (Docker)