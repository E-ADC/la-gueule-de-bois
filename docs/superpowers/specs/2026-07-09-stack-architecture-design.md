# La Gueule de Bois — Spec : stack technique & architecture

Date : 2026-07-09
Statut : validé (brainstorm équipe + Claude)

## Contexte

Projet école ESGI (module UML, groupe de 3). Application sociale « La Gueule de Bois » : enregistrement de soirées, témoignages, scoring, badges, classements, groupes. 22 cas d'utilisation définis (voir `cas-utilisation-final-22 .md` et `fiches-cas-utilisation.md`).

Contrainte forte du sujet : **la conception UML doit être terminée avant toute ligne de code applicatif** (perte de points sinon). Ce spec fige les choix techniques nécessaires pour finir l'UML (notamment le prestataire externe), puis cadre l'implémentation qui suivra.

Deadline : < 1 mois. Hébergement cible : VPS OVH existant (2 vCores, 4 Go RAM, 40 Go disque, Docker installé — voir `vps-setup-status (1).md`).

## Décisions

### Stack

| Composant | Choix | Justification |
|---|---|---|
| Backend | Go 1.23+, routeur **chi**, accès DB **pgx**, migrations **golang-migrate** | Go validé explicitement par le prof, connu de l'équipe, binaire unique léger (VPS 4 Go), typage fort cohérent avec la conception UML |
| Frontend | **React + Vite + TypeScript**, React Router, fetch | SPA choisie par l'équipe ; swipe (UC12) via framer-motion ou CSS transform |
| Direction visuelle | Thème **« Ambrée »** — palette bière (caramel/cuivre/brun torréfié), Helvetica unique, bordures franches, ombres pleines | Choisi parmi 3 propositions le 2026-07-09 ; tokens dans `docs/design/theme-ambree.css` |
| Base de données | **PostgreSQL 16** (conteneur Docker) | Relationnel classique, colle aux entités métier |
| Prestataire externe | **Resend** (API email, SDK Go officiel) | Contrainte sujet « un prestataire extérieur ». 3 000 emails/mois gratuits. Appel API HTTPS (pas de SMTP sortant → pas de blocage port 25) |
| Déploiement | **Docker Compose** sur le VPS : nginx → API Go → PostgreSQL | Simple, reproductible, déjà testé sur le VPS |

### Prestataire externe — intégration UML

Resend devient **acteur secondaire « Service d'email »** sur :
- UC09 (inviter témoin) — notification d'invitation
- UC14 (débloquer badge) — notification de badge
- UC21 (demande d'ami) — notification au destinataire
- UC22 (traiter signalement) — notification à l'auteur du témoignage

À répercuter dans le diagramme de cas d'utilisation **et** les fiches textuelles avant les diagrammes suivants.

Limitation acceptée : sans nom de domaine vérifié, Resend n'envoie qu'aux adresses de l'équipe (mode test) — suffisant pour la démo. Achat domaine OVH (~5 €) optionnel plus tard.

### Architecture backend

Trois couches :

```
handlers/     HTTP, JSON, validation des entrées, codes d'erreur
services/     logique métier : scoring (UC16), badges (UC14), modération (UC22)
repository/   accès PostgreSQL via pgx
```

- Entités métier (alignées sur le diagramme de classes à vérifier) : `User`, `Soiree`, `Temoignage`, `Vote`, `Badge`, `Groupe`, `Signalement`, `Photo`.
- L'envoi d'email passe par une interface `Notifier` (implémentation Resend) : mockable en test, et matérialise proprement l'acteur externe dans les diagrammes de séquence.
- Auth : **sessions cookie HttpOnly** (pas de JWT) — SPA servie sur le même domaine, plus simple et plus sûr.

### Photos de soirées — version minimale

Fonctionnalité incluse, implémentation la plus simple possible :

- Upload de photos sur une soirée (UC06 création / UC07 modification — pas de nouveau cas d'utilisation, enrichissement des fiches existantes).
- `multipart/form-data`, validation type MIME (jpeg/png/webp) + taille max 5 Mo. Fichier écrit tel quel sur disque (volume Docker), renommé en UUID. **Aucun traitement** : pas de redimensionnement, pas de miniatures, pas de galerie avancée.
- Servies en statique par nginx (`/uploads/...`).
- Entité `Photo` dans le diagramme de classes : association Soiree 1—* Photo.

### Structure du repo

```
/backend      code Go
/frontend     app React/Vite
/docs         conception UML + specs
docker-compose.yml
```

### API — gestion d'erreurs

Réponses d'erreur JSON uniformes : `{ "error": "...", "code": "..." }`.
Correspondance avec les exceptions des fiches UC : données invalides → 400, non-propriétaire/non-membre → 403, ressource inexistante → 404, doublon (email pris, déjà voté, déjà signalé) → 409.

### Tests

- Go : tests table-driven sur `services/` en priorité (scoring et badges = cœur métier, UC16/UC14).
- Frontend : lint + build en CI ; tests composants seulement si temps disponible.
- CI GitHub Actions : remplacer les TODO du squelette par `go vet` + `go test ./...` + `go build` côté back, `npm run lint` + `npm run build` côté front.

### Déploiement VPS

Compose avec trois services :
1. **nginx** — sert le build frontend statique + les photos (`/uploads`), proxifie `/api` vers l'API Go, terminaison TLS (Certbot une fois le domaine acheté)
2. **api** — binaire Go
3. **db** — PostgreSQL 16 + volume persistant

## Bonus (si temps en fin de projet, non prioritaire)

- Améliorations photos : miniatures, redimensionnement, galerie
- Cas d'utilisation coupés de la liste à 38 (réintégrables facilement)
- CI/CD complet avec déploiement auto (partie « bonus » confirmée par le prof)

## Hors périmètre

- Stockage cloud (S3 et assimilés)
- Notifications push / temps réel
- App mobile native / PWA

## Ordre d'exécution (rappel contrainte sujet)

1. Intégrer l'acteur « Service d'email » (Resend) dans le diagramme de cas d'utilisation + fiches ; enrichir UC06/UC07 avec les photos
2. Vérifier le diagramme de classes du collègue (7+ classes, cohérence avec les 22 UC, faire apparaître le prestataire)
3. Diagrammes d'objets pré/post pour chaque cas illustré
4. 7+ diagrammes de séquence (boucles/alternatives/références), cohérents avec les UC
5. **Ensuite seulement** : code (backend, frontend), CI, déploiement VPS
