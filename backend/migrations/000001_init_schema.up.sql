-- Schéma initial de "La Gueule de Bois".
-- Choix non précisés par les fiches UC, tranchés ici (documentés en ligne) :
--   * identifiants numériques BIGSERIAL (plus simples à manipuler qu'UUID
--     pour un projet école, sauf pour les fichiers uploadés où l'UUID sert
--     de nom de fichier anti-collision, cf. spec photos).
--   * `users.role` : 'user' | 'moderator' (UC22 nécessite un acteur
--     Modérateur, la fiche ne précise pas de table de rôles dédiée).
--   * `users.score` dénormalisé et recalculé par UC16 : évite de recalculer
--     à la volée à chaque lecture (classement UC17/UC20).

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    pseudo        VARCHAR(32)  NOT NULL UNIQUE,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    avatar        TEXT         NOT NULL DEFAULT '',
    bio           TEXT         NOT NULL DEFAULT '',
    score         INTEGER      NOT NULL DEFAULT 0,
    role          VARCHAR(16)  NOT NULL DEFAULT 'user',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    token      TEXT PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

CREATE TABLE soirees (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    titre       VARCHAR(120) NOT NULL,
    date_soiree TIMESTAMPTZ NOT NULL,
    lieu        VARCHAR(180) NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_soirees_user_id ON soirees(user_id);

-- Photo : association Soiree 1—* Photo (spec "photos minimales", pas de
-- redimensionnement ni de galerie avancée).
CREATE TABLE photos (
    id         BIGSERIAL PRIMARY KEY,
    soiree_id  BIGINT NOT NULL REFERENCES soirees(id) ON DELETE CASCADE,
    chemin     TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_photos_soiree_id ON photos(soiree_id);

-- Invitation d'un témoin sur une soirée (UC09). Nécessaire pour vérifier
-- la pré-condition de UC11 ("témoin invité"). Non nommée explicitement
-- dans la fiche, déduite du besoin.
CREATE TABLE temoin_invitations (
    id         BIGSERIAL PRIMARY KEY,
    soiree_id  BIGINT NOT NULL REFERENCES soirees(id) ON DELETE CASCADE,
    invite_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (soiree_id, invite_id)
);

CREATE TABLE temoignages (
    id         BIGSERIAL PRIMARY KEY,
    soiree_id  BIGINT NOT NULL REFERENCES soirees(id) ON DELETE CASCADE,
    auteur_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    contenu    TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_temoignages_soiree_id ON temoignages(soiree_id);

-- Vote (swipe) sur un témoignage, un seul par (temoignage, utilisateur) —
-- règle explicite de la fiche UC12 ("a déjà voté -> vote ignoré").
CREATE TABLE votes (
    id            BIGSERIAL PRIMARY KEY,
    temoignage_id BIGINT NOT NULL REFERENCES temoignages(id) ON DELETE CASCADE,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    valeur        SMALLINT NOT NULL CHECK (valeur IN (-1, 1)),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (temoignage_id, user_id)
);

-- Badges : critère unique = seuil de score (UC14 ne mentionne pas d'autre
-- critère). Table peuplée par une migration de données (voir 000002).
CREATE TABLE badges (
    id          BIGSERIAL PRIMARY KEY,
    code        VARCHAR(64) NOT NULL UNIQUE,
    nom         VARCHAR(120) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    seuil_score INTEGER NOT NULL
);

CREATE TABLE user_badges (
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    badge_id    BIGINT NOT NULL REFERENCES badges(id) ON DELETE CASCADE,
    debloque_le TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, badge_id)
);

CREATE TABLE groupes (
    id          BIGSERIAL PRIMARY KEY,
    nom         VARCHAR(80) NOT NULL UNIQUE,
    createur_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE groupe_membres (
    groupe_id BIGINT NOT NULL REFERENCES groupes(id) ON DELETE CASCADE,
    user_id   BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (groupe_id, user_id)
);

-- Signalement d'un témoignage (UC13), traité par un modérateur (UC22).
CREATE TABLE signalements (
    id            BIGSERIAL PRIMARY KEY,
    temoignage_id BIGINT NOT NULL REFERENCES temoignages(id) ON DELETE CASCADE,
    auteur_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    motif         TEXT NOT NULL,
    statut        VARCHAR(32) NOT NULL DEFAULT 'en_attente',
    traite_par_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    traite_le     TIMESTAMPTZ,
    -- UC13 : "déjà signalé par cet utilisateur -> ignoré"
    UNIQUE (temoignage_id, auteur_id)
);

-- Demande d'ami (UC21). Modèle et table prêts ; handler HTTP en TODO
-- dans ce squelette (cf. README).
CREATE TABLE demandes_amis (
    id               BIGSERIAL PRIMARY KEY,
    demandeur_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    destinataire_id  BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    statut           VARCHAR(16) NOT NULL DEFAULT 'en_attente',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (demandeur_id, destinataire_id)
);
