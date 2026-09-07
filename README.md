# EventApp Backend API

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.10.0-008080?style=flat&logo=gin&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-modernc.org%2Fsqlite-003B57?style=flat&logo=sqlite&logoColor=white)
![golang-jwt](https://img.shields.io/badge/golang--jwt-v5.2.1-F28D35?style=flat)
![Swagger](https://img.shields.io/badge/Swagger-OpenAPI_3.0-85EA2D?style=flat&logo=swagger&logoColor=black)
![godotenv](https://img.shields.io/badge/godotenv-v1.5.1-222222?style=flat)

## Über das Projekt
Das EventApp Backend ist eine RESTful API, die in Go(Golang) entwickelt wurde.
Sie dient als Grundlage für eine Event-Management-Plattform, auf der Benutzer sich registrieren, Events erstellen, verwalten und sich als Teilnehmer für Events eintragen können.
Der Fokus dieses Backends liegt auf sauberer Architektur, hoher Performance durch das Gin-Framework und einer reibungslosen Integration mit dem Frontend.

---

## Features
- **Benutzerverwaltung:** Registrierung und Login.
- **Event-Management:** CRUD-Operationen (Erstellen, Lesen, Aktualisieren, Löschen) für Events.
- **Teilnehmer-Management:** Benutzer können sich für Events anmelden und wieder abmelden; Abfrage von Teilnehmern pro Event.
- **API-Dokumentation:** Automatisch generierte und interaktive Swagger-UI.
- **Sicherheit:** Geschützte Routen mittels JSON Web Tokens(JWT).
- **CORS-Support:** Vorkonfiguriertes Cross-Origin Resource Sharing für die lokale Frontend-Entwicklung.

---

## Technologie-Stack
- **Sprache:** Go (Golang)
- **Web-Framework:** [Gin](https://gin-gonic.com/)
- **Datenbank:** SQLite (via [`modernc.org/sqlite`](https://gitlab.com/cznic/sqlite) – CGO-frei für einfaches Cross-Compiling)
- **Authentifizierung:** JWT ([`golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt))
- **Umgebungsvariablen:** [`joho/godotenv`](https://github.com/joho/godotenv)
- **Dokumentation:** Swagger / OpenAPI ([`swaggo/gin-swagger`](https://github.com/swaggo/gin-swagger))


---

## Repository-Struktur
```
.
├─ cmd
│  └──api
│     ├── main.go       # Einstiegspunkt der Applikation
│     ├── server.go     # HTTP-Server Konfiguration
│     ├── routes.go     # Definition der API-Endpunkte
│     ├── middleware.go # Auth- und CORS-Middlewares 
│     ├── auth.go       # Auth-Handler (Registrierung, Login & JWT-Erstellung)
│     ├── events.go     # Handler für Event-CRUD-Operationen und Teilnehmer-Verwaltung
│     └── context.go    # Helper zum Extrahieren des authentifizierten Benutzers aus dem Gin-Kontext
│ 
├── internal
│   ├── database          # Datenbankverbindung und Model
│   │   ├── attendees.go  # Datenbankoperationen für Event-Teilnehmer(CRUD & M:N Verknüpfungen)
│   │   ├── events.go     # Event-Modell, Struct-Validierung, CRUD-Datenbankoperationen für Events
│   │   ├── models.go     # Haupt-Datenbank-Wrapper
│   │   └── users.go      # User-Modell und Auth-Datenbankmethoden
│   └── env               # Hilfsfunktionen zum Laden der .env Variablen
│       └── env.go
├── docs                  # Automatisch generierte Swagger-Dateien
├── go.mod                # Go Module Abhängigkeiten
├── go.sum
└── README.md
```
---

##  Umgebungsvariablen
#### `/EventApp/.env`

```env
PORT=8080
JWT_SECRET=dein-sehr-geheimes-jwt-secret
CLIENT_ORIGIN=http://localhost:5173
```

---

## API-Dokumentation

Die API ist vollständig mit Swagger dokumentiert. Sobald der Server läuft, kann die interaktive Dokumentation im Browser aufgerufen werden unter:

http://localhost:8080/swagger/index.html


### Wichtige Endpunkte (Auszug):

- **POST /api/v1/auth/register** - Neuen Benutzer registrieren
- **POST /api/v1/auth/login** - Einloggen & JWT erhalten
- **GET /api/v1/events**  - Alle Events abrufen
- **POST /events (Auth required)** - Neues Event erstellen
- **POST /events/:id/attendees/:userId (Auth required)** - An Event teilnehmen

---

## Datenbank-Schema

Das Projekt verwendet eine relationale SQLite-Datenbank mit folgenden logischen Hauptkomponenten:
- Users: Speichert Benutzerdaten (ID, Name, E-Mail, Passwort-Hash)
- Events: Speichert Event-Details (ID, Ersteller-ID, Titel, Beschreibung, Datum, Veranstaltungsort)
- Event_Attendees (Join Table): Verknüpft Benutzer (als Teilnehmer) mit Events (M:N-Beziehung)

### `users`
Speichert alle registrierten Benutzer und deren Zugangsdaten.

| Feld | Typ | Eigenschaften | Beschreibung |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY, AUTOINCREMENT | Eindeutige Benutzer-ID |
| `email` | TEXT | UNIQUE, NOT NULL | E-Mail-Adresse für den Login |
| `password` | TEXT | NOT NULL | Gehashtes Passwort (sicher verschlüsselt) |
| `created_at`| DATETIME| DEFAULT CURRENT_TIMESTAMP | Zeitpunkt der Registrierung |

### `events`
Speichert die Kerninformationen zu den erstellten Events.

| Feld | Typ | Eigenschaften | Beschreibung |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY, AUTOINCREMENT | Eindeutige Event-ID |
| `title` | TEXT | NOT NULL | Titel der Veranstaltung |
| `description`| TEXT | NOT NULL | Ausführliche Beschreibung des Events |
| `location` | TEXT | NOT NULL | Veranstaltungsort (falls zutreffend) |
| `date` | DATETIME| NOT NULL | Geplantes Datum und Uhrzeit |
| `user_id` | INTEGER | FOREIGN KEY | ID des Erstellers (verweist auf `users.id`) |

### `event_attendees` (Join-Tabelle)
Eine M:N-Verknüpfungstabelle, die regelt, welche Benutzer an welchen Events teilnehmen.

| Feld | Typ | Eigenschaften | Beschreibung |
|---|---|---|---|
| `event_id` | INTEGER | FOREIGN KEY | Verweist auf `events.id` |
| `user_id` | INTEGER | FOREIGN KEY | Verweist auf `users.id` |

*(Hinweis: Der Primary Key dieser Tabelle ist ein Composite Key aus `event_id` und `user_id`. Dies stellt auf Datenbankebene sicher, dass sich ein Benutzer nicht mehrfach für dasselbe Event anmelden kann.)*

---

## Architektur

Das Projekt folgt einer klaren Separation of Concerns (SoC):

- **Routing & Handling (cmd/api):** Nimmt HTTP-Requests entgegen, validiert Inputs und steuert die Responses (Gin Context).

- **Business Logic & Data Access (internal/database):** Kapselt alle SQL-Abfragen und Datenbank-Interaktionen in Form von Models (app.models).

- **Konfiguration (internal/env):** Zentralisiertes Einlesen und Fallbacks für Konfigurationswerte.

### System-Architektur
```
┌──────────────────────────────────────────── ▼ ────────────────────────────────────────┐
│                             Backend (Go + Gin Framework)                              │
│  ┌─────────────────────────┐      ┌───────────────────────────┐  ┌──────────────────┐ │
│  │       Routes            │      │ Application Methods       │  │    Middleware    │ │
│  │                         │      │                           │  │                  │ │
│  │ Unprotected:            │────▸ │ Auth                      │  │ - CORSMiddleware │ │
│  │ - /auth/register        │      │  - registerUser           │  │   (Global)       │ │
│  │ - /auth/login           │      │  - login                  │  │                  │ │
│  │                         │      │                           │  │                  │ │
│  │ - /events               │      │ Events (Read)             │  │ - AuthMiddleware │ │
│  │ - /events/:id           │      │  - getAllEvents           │  │   (JWT für       │ │
│  │ - /events/:id/attendees │      │  - getEvent               │  │   geschützte     │ │
│  │ - /attendees/:id/events │      │                           │  │   Routen)        │ │
│  │ - /swagger/*any         │      │ Events (Write)            │  │                  │ │
│  │                         │      │  - createEvent            │  │                  │ │
│  │ Protected (JWT):        │      │  - updateEvent            │  │                  │ │
│  │ - POST /events          │      │  - deleteEvent            │  │                  │ │
│  │ - PUT /events/:id       │      │                           │  │                  │ │
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

---

## Sicherheit & Authentifizierung

- **Stateless Auth:** Die Authentifizierung erfolgt über JWT. Der Token muss bei geschützten Routen im Header als Authorization: Bearer <token> mitgesendet werden.

- **Passwort-Sicherheit:** Passwörter werden vor dem Speichern in der Datenbank sicher gehasht.

- **CORS:** Geregelt über eine eigene Middleware. Für die lokale Entwicklung werden Anfragen vom in .env definierten CLIENT_ORIGIN akzeptiert.

---

## Roadmap

[ ] Integration mit dem Frontend (React)

[ ] Hinzufügen von Unit- und Integration-Tests

[ ] Containerisierung des Backends (Docker)