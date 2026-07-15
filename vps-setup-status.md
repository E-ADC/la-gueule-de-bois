# VPS OVH — État du setup

## Infos serveur

- **Modèle** : VPS-1 2027 — 2 vCores, 4 Go RAM, 40 Go stockage
- **IPv4** : 164.132.100.28
- **IPv6** : 2001:41d0:404:200::42ee
- **OS** : Ubuntu (image cloud OVH, user par défaut `ubuntu`)
- **Hostname** : vps-33a05772

## Accès SSH

- Connexion par **clé uniquement** (`~/.ssh/id_ed25519`), password auth désactivé
- **Port SSH : 2222** (pas 22 — changé pour réduire le bruit des bots)
- `PermitRootLogin no` (explicite, dans `/etc/ssh/sshd_config.d/99-custom-port.conf`)
- `ssh.socket` désactivé, `ssh.service` classique utilisé à la place (le socket activation par défaut d'Ubuntu écrasait la config de port)

Commande de connexion :
```bash
ssh -i ~/.ssh/id_ed25519 -p 2222 ubuntu@164.132.100.28
```

### Points de config non-standard à connaître
- `/etc/ssh/sshd_config.d/50-cloud-init.conf` : généré par cloud-init, avait `PasswordAuthentication yes` — corrigé en `no`
- `/etc/ssh/sshd_config.d/99-custom-port.conf` : fichier custom contenant `Port 2222` + `PermitRootLogin no`
- L'ordre de lecture des `Include` dans `sshd_config` est important : les fichiers `.d/` sont lus tôt (ligne 24), donc tout override doit être fait dans le bon fichier `.d/`, pas dans `sshd_config` principal

## Firewall (UFW)

- Actif, policy par défaut : deny incoming / allow outgoing
- Ports ouverts : **2222/tcp** (SSH), **80/tcp** (HTTP), **443/tcp** (HTTPS)
- IPv4 et IPv6 couverts

## Fail2ban

- Actif, jail `sshd` configuré sur le **port 2222** (pas 22)
- `maxretry = 3`, `bantime = 3600`, `findtime = 600`
- Config dans `/etc/fail2ban/jail.local`

## Mises à jour automatiques

- `unattended-upgrades` installé et configuré

## Docker

- Docker + Docker Compose (plugin) installés via `get.docker.com`
- User `ubuntu` ajouté au groupe `docker` (pas besoin de sudo pour les commandes docker)
- Testé et fonctionnel (`docker run hello-world` OK)

## Pas encore fait

- Nginx / reverse proxy
- Nom de domaine (aucun configuré pour l'instant côté OVH)
- Certificat TLS (Certbot)
- Déploiement de l'application — stack choisie (Go + React + PostgreSQL + Resend, voir `docs/superpowers/specs/2026-07-09-stack-architecture-design.md`), implémentation en cours
- Configuration base de données

## Contexte projet

Projet d'école (module UML, groupe de 3, ESGI, sujet libre) — application **"La Gueule de Bois"**. Ce VPS est destiné à héberger le déploiement final.

Évaluation : 30% sur qualité UML/dev/soutenance. Le sujet impose explicitement que **la conception UML soit terminée avant le début du code** (perte de points sinon).

## Contraintes obligatoires du sujet (validées avec le prof)

- **Un prestataire extérieur** obligatoire — ✅ choisi : **Resend** (API email), acteur secondaire « Service d'email » rattaché à UC09 (inviter témoin), UC14 (débloquer badge), UC21 (demande d'ami), UC22 (traiter signalement).
- **Minimum 3 entités métier** internes avec relations
- Diagramme de cas d'utilisation : min. 22 use cases (fait, voir ci-dessous)
- Diagramme de classes : min. 7 classes métier avec attributs/méthodes/relations (fait par un collègue, cohérence à vérifier)
- Diagrammes d'objets pré/post-séquence pour chaque cas illustré (pas encore fait)
- Diagrammes de séquence : min. 7 scénarios avec boucles/alternatives/références (pas encore fait)
- Cohérence obligatoire entre cas d'utilisation et séquences

**Stack technique** : tranchée le 2026-07-09 — **Go (chi + pgx) / React + Vite + TypeScript / PostgreSQL 16 / Resend / Docker Compose**. Détails dans `docs/superpowers/specs/2026-07-09-stack-architecture-design.md`.

⚠️ Un autre document de sujet plus ancien (SUJET_PROJET.pdf) imposait une stack précise (frontend framework + Node/PHP/Rust/Django + repo GitHub structuré + CI/CD) — le prof a confirmé à l'oral que cette partie devient **bonus**, pas obligatoire. Garder une trace écrite de cette confirmation si possible.

## Conception UML — état d'avancement

- ✅ 22 cas d'utilisation définis et documentés (fiches textuelles complètes : résumé, acteurs, pré-conditions, résultats, description, exceptions)
- ✅ Acteurs : généralisation `Visiteur ← Utilisateur ← Modérateur`
- ✅ Relations d'inclusion : UC06/07/08/11 → incluent UC16 (recalcul score) → inclut UC14 (déblocage badge)
- ✅ Généralisation de cas : UC20 (classement groupe) généralisé par UC17 (classement global)
- ✅ Prestataire externe (Resend, acteur « Service d'email ») intégré dans les fiches **et** le diagramme de cas d'utilisation (`diagramme_de_cas_d_utilisation.eddx`, 2026-07-10)
- ⚠️ Diagramme de classes — **fait par un collègue**, pas encore vérifié côté cohérence : à contrôler que les 7+ classes couvrent bien toutes les entités des 22 UC (soirée, témoignage, badge, score, groupe, utilisateur, signalement...) et que le prestataire externe y apparaît
- ❌ Diagrammes d'objets — non commencés
- ❌ Diagrammes de séquence — non commencés

Fichiers sources dans le repo : `cas-utilisation-final-22 .md`, `fiches-cas-utilisation.md`, `diagramme_de_cas_d_utilisation.eddx` (+ export PDF `Diagrammes de Cas d'Utilisation.pdf` — à régénérer, il date d'avant l'ajout du prestataire)

## Repo GitHub

- **URL** : https://github.com/E-ADC/la-gueule-de-bois (public)
- **Branche par défaut** : `dev`
- **Structure de branches** :
  - `main` — protégée (PR obligatoire + 1 approval + no force push), déploiement prod
  - `preprod` — protégée (mêmes règles), staging
  - `dev` — libre (commit/push/merge/rollback directs), branche de travail principale
  - Branches éphémères : convention `type(scope):description` (types : `feat`, `fix`, `style`), documentée dans `CONTRIBUTING.md`
- **CI** : GitHub Actions (`.github/workflows/ci.yml`) — squelette fonctionnel (jobs lint/build/test), contient des `TODO` à remplacer par les vraies commandes une fois la stack choisie
- **Issues actives** : ~~choix stack~~ ✅, ~~intégration prestataire externe~~ ✅, vérification diagramme de classes, diagrammes d'objets, diagrammes de séquence, complétion CI, déploiement VPS

## Prochaine étape immédiate

Implémentation (backend Go, frontend React, docker-compose) sur branches `feat/` créées depuis `dev` — jamais directement sur `dev`/`preprod`/`main`. En parallèle côté UML : vérifier le diagramme de classes du collègue, puis diagrammes d'objets et de séquence (rappel contrainte sujet : l'UML doit être terminé avant le code pour ne pas perdre de points).
