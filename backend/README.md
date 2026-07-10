# La Gueule de Bois — Backend

API Go du projet école « La Gueule de Bois » : soirées, témoignages, votes,
scoring, badges, groupes, signalements. Stack : Go 1.23+, chi (routeur),
pgx/pgxpool (accès Postgres), golang-migrate (migrations SQL), PostgreSQL 16.

Trois couches : `internal/handlers` (HTTP/JSON/validation) →
`internal/services` (logique métier) → `internal/repository` (accès
Postgres via pgx).

## Prérequis

- Go 1.23+ (ce squelette a été développé/testé avec Go 1.25 via
  `go.mod`/toolchain auto — `go build` télécharge le bon toolchain si besoin).
- PostgreSQL 16 accessible (local, Docker, ou VPS).
- [golang-migrate CLI](https://github.com/golang-migrate/migrate) pour
  appliquer les migrations (`brew install golang-migrate` ou binaire du
  repo GitHub).

## Variables d'environnement

| Variable          | Obligatoire | Défaut       | Description                                  |
|-------------------|-------------|--------------|-----------------------------------------------|
| `DATABASE_URL`    | oui         | —            | DSN Postgres, ex. `postgres://user:pass@localhost:5432/gdb?sslmode=disable` |
| `SESSION_SECRET`  | oui         | —            | Secret HMAC de signature des cookies de session (chaîne aléatoire longue) |
| `PORT`            | non         | `8080`       | Port d'écoute HTTP                            |
| `UPLOAD_DIR`      | non         | `./uploads`  | Dossier de stockage des photos uploadées      |
| `RESEND_API_KEY`  | non         | —            | Clé API Resend ; absente = notifications simulées (MockNotifier), utile en dev |

## Lancement local

```bash
# 1. Démarrer Postgres (exemple rapide avec Docker) :
docker run --name gdb-postgres -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gdb -p 5432:5432 -d postgres:16

# 2. Appliquer les migrations :
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/gdb?sslmode=disable" up

# 3. Lancer l'API :
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/gdb?sslmode=disable"
export SESSION_SECRET="change-moi-en-une-longue-chaine-aleatoire"
go run ./cmd/api
```

L'API écoute par défaut sur `http://localhost:8080`.

## Commandes utiles

```bash
go build ./...      # compile tout le module
go vet ./...         # analyse statique
go test ./...         # tests (aucune base de données requise)
gofmt -l .            # vérifie le formatage (vide = OK)
```

## Endpoints exposés

Toutes les routes sont préfixées par `/api`. Les routes marquées 🔒
nécessitent le cookie de session (`RequireAuth`).

| Méthode | Route                                  | UC        | Description |
|---------|-----------------------------------------|-----------|--------------|
| POST    | `/api/auth/register`                    | UC01      | Inscription (crée le compte + connecte) |
| POST    | `/api/auth/login`                       | UC02      | Connexion |
| POST    | `/api/auth/logout`                      | UC03      | Déconnexion |
| GET     | `/api/auth/me` 🔒                        | —         | Utilisateur courant |
| GET     | `/api/users/{id}` 🔒                     | UC05      | Profil public d'un utilisateur |
| GET     | `/api/classement` 🔒                     | UC17      | Classement global |
| GET     | `/api/groupes/{id}/classement` 🔒        | UC20      | Classement restreint à un groupe |
| GET     | `/api/me/badges` 🔒                      | UC15      | Mes badges obtenus / à débloquer |
| GET     | `/api/soirees` 🔒                        | UC10      | Historique de mes soirées |
| POST    | `/api/soirees` 🔒                        | UC06      | Créer une soirée (inclut UC16) |
| GET     | `/api/soirees/{id}` 🔒                   | —         | Détail d'une soirée + ses photos |
| PUT     | `/api/soirees/{id}` 🔒                   | UC07      | Modifier une soirée (inclut UC16) |
| DELETE  | `/api/soirees/{id}` 🔒                   | UC08      | Supprimer une soirée (inclut UC16) |
| POST    | `/api/soirees/{id}/photos` 🔒            | UC06/07   | Upload photo (multipart, jpeg/png/webp, max 5 Mo) |
| POST    | `/api/soirees/{id}/temoins` 🔒           | UC09      | Inviter un témoin (email via Resend) |
| POST    | `/api/soirees/{id}/temoignages` 🔒       | UC11      | Ajouter un témoignage (inclut UC16) |
| GET     | `/api/soirees/{id}/temoignages` 🔒       | —         | Lister les témoignages d'une soirée |
| POST    | `/api/temoignages/{id}/votes` 🔒         | UC12      | Voter (swipe) sur un témoignage |
| GET     | `/healthz`                               | —         | Ping santé |
| GET     | `/uploads/*`                             | —         | Sert les photos en local (nginx en prod) |

### TODO (non implémentés dans ce squelette, marqués 501 dans le routeur)

Modèles, migrations et repositories sont prêts pour ces cas d'utilisation ;
il reste à écrire le service + handler sur le modèle des routes ci-dessus :

- `POST /api/temoignages/{id}/signalements` — UC13 (signaler un témoignage)
- `GET /api/signalements`, `POST /api/signalements/{id}/traiter` — UC22 (modérateur)
- `POST /api/groupes` — UC18 (créer un groupe)
- `POST /api/groupes/{id}/membres` — UC19 (rejoindre un groupe)
- `POST /api/amis/demandes` — UC21 (demande d'ami)

## Format d'erreur uniforme

```json
{ "error": "message lisible", "code": "invalid_input" }
```

Mapping : `invalid_input`/`ErrValidation` → 400 · `forbidden` (non
propriétaire/non membre) → 403 · `not_found` → 404 · `conflict` (doublon :
email pris, déjà voté, déjà signalé...) → 409.

## Choix métier tranchés (fiches UC muettes sur ces points)

- **Score (UC16)** : `+10` par soirée créée, `+5` par témoignage reçu sur
  ses soirées, `+1`/`-1` par vote positif/négatif reçu dessus, plancher à 0.
  Voir `internal/services/scoring.go`.
- **Recalcul du score sur vote (UC12)** : la fiche UC12 ne liste pas UC16
  parmi les cas inclus (contrairement à UC06/07/08/11), mais comme le score
  intègre les votes, `VoteService.Cast` déclenche quand même un recalcul —
  extension mineure et documentée, pas une contradiction.
- **Badges (UC14)** : critère unique = seuil de score (4 badges seedés en
  migration `000002`, de 10 à 300 points). Voir `internal/services/badges.go`.
- **Table `temoin_invitations`** : nécessaire pour vérifier la
  pré-condition « témoin invité » de UC11, non nommée explicitement par la
  fiche UC09.
- **Identifiants** : `BIGSERIAL` pour toutes les entités métier ; UUID
  réservé au nommage des fichiers photo uploadés (anti-collision disque).
- **Sessions** : cookie opaque signé HMAC (`SESSION_SECRET`) plutôt que
  JWT (imposé par la spec) ; TTL de 7 jours.

## Docker

```bash
docker build -t gdb-backend .
docker run -p 8080:8080 \
  -e DATABASE_URL=... -e SESSION_SECRET=... \
  -v gdb-uploads:/app/uploads \
  gdb-backend
```

Le binaire embarque uniquement `/app/api` et `/app/migrations` (build
multi-stage, image finale `distroless/static`).
