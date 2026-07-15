# Fiches textuelles des cas d'utilisation — La Gueule de Bois

Formalisme utilisé (cours "Modélisation avec UML") : cas, résumé, acteur primaire, acteur secondaire, pré-conditions, résultats, description, exceptions.

---

### UC01 — S'inscrire

| Champ | Contenu |
|---|---|
| **résumé** | Un visiteur crée un compte utilisateur pour accéder à l'application |
| **acteur primaire** | Visiteur |
| **acteur secondaire** | — |
| **pré-conditions** | Le visiteur n'est pas connecté ; l'email/pseudo n'est pas déjà utilisé |
| **résultats** | Un compte utilisateur est créé et activé |
| **description** | 1. Le visiteur accède au formulaire d'inscription<br>2. Il saisit pseudo, email, mot de passe<br>3. Le système vérifie l'unicité de l'email/pseudo<br>4. Le système crée le compte<br>5. Le système connecte automatiquement l'utilisateur |
| **exceptions** | Email/pseudo déjà utilisé → inscription refusée |

---

### UC02 — Se connecter

| Champ | Contenu |
|---|---|
| **résumé** | Un visiteur s'authentifie pour accéder à son compte |
| **acteur primaire** | Visiteur |
| **acteur secondaire** | — |
| **pré-conditions** | Le visiteur possède un compte existant |
| **résultats** | L'utilisateur est authentifié et accède à son espace |
| **description** | 1. Le visiteur saisit identifiant/mot de passe<br>2. Le système vérifie les identifiants<br>3. Le système ouvre une session utilisateur |
| **exceptions** | Identifiants invalides → accès refusé |

---

### UC03 — Se déconnecter

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur met fin à sa session |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | L'utilisateur est connecté |
| **résultats** | La session est fermée, retour à l'état visiteur |
| **description** | 1. L'utilisateur clique sur « se déconnecter »<br>2. Le système ferme la session |
| **exceptions** | — |

---

### UC04 — Modifier son profil

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur met à jour ses informations personnelles |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | Le profil est mis à jour |
| **description** | 1. L'utilisateur accède à son profil<br>2. Il modifie les champs souhaités (pseudo, avatar, bio…)<br>3. Le système valide et enregistre les modifications |
| **exceptions** | Donnée invalide (ex. pseudo déjà pris) → modification refusée |

---

### UC05 — Consulter le profil d'un autre utilisateur

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur consulte les informations publiques d'un autre membre |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté ; profil cible existant et visible |
| **résultats** | Affichage du profil consulté |
| **description** | 1. L'utilisateur recherche/sélectionne un autre utilisateur<br>2. Le système affiche les informations publiques (score, badges, historique visible) |
| **exceptions** | Profil inexistant ou privé → accès refusé |

---

### UC06 — Créer une soirée

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur enregistre une nouvelle soirée |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Système de scoring |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | La soirée est créée ; le score de l'utilisateur est recalculé |
| **description** | 1. L'utilisateur saisit les informations de la soirée (date, lieu, participants…)<br>2. Il peut joindre des photos (optionnel)<br>3. Le système enregistre la soirée et ses photos<br>4. Le système inclut **UC16** (Recalculer le score de l'utilisateur) |
| **exceptions** | Champs obligatoires manquants → création refusée ; photo invalide (format ou taille) → photo rejetée, création poursuivie sans elle |

---

### UC07 — Modifier une soirée

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur met à jour une soirée qu'il a créée |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Système de scoring |
| **pré-conditions** | La soirée existe et appartient à l'utilisateur |
| **résultats** | La soirée est mise à jour ; le score est recalculé |
| **description** | 1. L'utilisateur sélectionne la soirée à modifier<br>2. Il modifie les champs souhaités et peut ajouter/retirer des photos<br>3. Le système enregistre les modifications<br>4. Le système inclut **UC16** |
| **exceptions** | Utilisateur non propriétaire de la soirée → modification refusée |

---

### UC08 — Supprimer une soirée

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur retire une soirée qu'il a créée |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Système de scoring |
| **pré-conditions** | La soirée existe et appartient à l'utilisateur |
| **résultats** | La soirée est supprimée ; le score est recalculé |
| **description** | 1. L'utilisateur sélectionne la soirée à supprimer<br>2. Le système demande confirmation<br>3. Le système supprime la soirée<br>4. Le système inclut **UC16** |
| **exceptions** | Utilisateur non propriétaire → suppression refusée |

---

### UC09 — Inviter un témoin à une soirée

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur associe un autre utilisateur comme témoin d'une soirée |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Service d'email (prestataire Resend) ; Utilisateur invité (notifié) |
| **pré-conditions** | La soirée existe et appartient à l'utilisateur ; le témoin invité possède un compte |
| **résultats** | Le témoin est associé à la soirée (notification envoyée) |
| **description** | 1. L'utilisateur sélectionne une soirée<br>2. Il choisit un utilisateur à inviter comme témoin<br>3. Le système enregistre l'invitation<br>4. Le système notifie le témoin par email via le **Service d'email** |
| **exceptions** | Utilisateur invité inexistant → invitation refusée |

---

### UC10 — Consulter l'historique de ses soirées

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur visualise la liste de ses soirées passées |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | Affichage de la liste des soirées de l'utilisateur |
| **description** | 1. L'utilisateur accède à son historique<br>2. Le système affiche la liste chronologique des soirées |
| **exceptions** | — |

---

### UC11 — Ajouter un témoignage sur une soirée

| Champ | Contenu |
|---|---|
| **résumé** | Un témoin rédige un témoignage sur une soirée à laquelle il a été invité |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Système de scoring |
| **pré-conditions** | L'utilisateur est témoin invité de la soirée |
| **résultats** | Le témoignage est publié ; le score est recalculé |
| **description** | 1. Le témoin sélectionne la soirée<br>2. Il rédige et soumet son témoignage<br>3. Le système publie le témoignage<br>4. Le système inclut **UC16** |
| **exceptions** | Utilisateur non invité comme témoin → ajout refusé |

---

### UC12 — Swiper/voter sur un témoignage

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur réagit (vote) à un témoignage publié |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Le témoignage existe et est visible par l'utilisateur |
| **résultats** | Le vote est enregistré |
| **description** | 1. L'utilisateur consulte un témoignage<br>2. Il swipe/vote (positif ou négatif)<br>3. Le système enregistre le vote |
| **exceptions** | L'utilisateur a déjà voté sur ce témoignage → vote ignoré |

---

### UC13 — Signaler un témoignage

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur signale un témoignage qu'il juge inapproprié |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Le témoignage existe |
| **résultats** | Un signalement est créé, **en attente de traitement par un Modérateur (voir UC22)** |
| **description** | 1. L'utilisateur sélectionne un témoignage<br>2. Il choisit « signaler » et précise un motif<br>3. Le système enregistre le signalement |
| **exceptions** | Témoignage déjà signalé par cet utilisateur → signalement ignoré |

---

### UC14 — Débloquer un badge *(cas interne, sans acteur — inclus par UC16)*

| Champ | Contenu |
|---|---|
| **résumé** | Le système attribue automatiquement un badge à un utilisateur qui remplit les conditions, à l'issue d'un recalcul de score |
| **acteur primaire** | — *(aucun acteur externe : ce cas n'est jamais déclenché directement, uniquement inclus par UC16)* |
| **acteur secondaire** | Service d'email (prestataire Resend) ; Utilisateur (notifié) |
| **pré-conditions** | UC16 vient de s'exécuter |
| **résultats** | Un nouveau badge est attribué à l'utilisateur ; notification envoyée |
| **description** | 1. Le système évalue les critères de badges avec le score mis à jour<br>2. Si un critère est atteint, le système attribue le badge<br>3. Le système notifie l'utilisateur par email via le **Service d'email** |
| **exceptions** | Aucun critère atteint → aucun badge attribué |

---

### UC15 — Consulter ses badges

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur visualise les badges qu'il a obtenus |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | Affichage de la liste des badges obtenus |
| **description** | 1. L'utilisateur accède à sa page badges<br>2. Le système affiche les badges obtenus et à débloquer |
| **exceptions** | — |

---

### UC16 — Recalculer le score d'un utilisateur *(cas interne, sans acteur — inclus par UC06/07/08/11, inclut UC14)*

| Champ | Contenu |
|---|---|
| **résumé** | Cas système, jamais déclenché directement par un acteur : uniquement inclus par UC06, UC07, UC08 et UC11 après toute action affectant les soirées/témoignages d'un utilisateur |
| **acteur primaire** | — *(aucun acteur externe)* |
| **acteur secondaire** | — |
| **pré-conditions** | Le cas appelant (UC06/07/08/11) est en cours d'exécution |
| **résultats** | Le score de l'utilisateur est mis à jour |
| **description** | 1. Le système récupère les données de l'utilisateur<br>2. Il applique l'algorithme de calcul du score<br>3. Il met à jour le score stocké<br>4. Le système inclut **UC14** (Débloquer un badge) |
| **exceptions** | — |

---

### UC17 — Consulter le classement global

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur consulte le classement de tous les utilisateurs |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | Affichage du classement global trié par score |
| **description** | 1. L'utilisateur accède au classement<br>2. Le système affiche la liste des utilisateurs triée par score décroissant |
| **exceptions** | — |

---

### UC18 — Créer un groupe

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur crée un nouveau groupe d'amis |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Utilisateur connecté |
| **résultats** | Le groupe est créé, l'utilisateur en devient membre |
| **description** | 1. L'utilisateur saisit le nom du groupe<br>2. Le système crée le groupe<br>3. Le système ajoute l'utilisateur comme membre |
| **exceptions** | Nom de groupe déjà utilisé → création refusée |

---

### UC19 — Rejoindre un groupe

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur intègre un groupe existant |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | Le groupe existe |
| **résultats** | L'utilisateur devient membre du groupe |
| **description** | 1. L'utilisateur recherche/sélectionne un groupe<br>2. Il demande à rejoindre (ou saisit un code d'invitation)<br>3. Le système l'ajoute comme membre |
| **exceptions** | Utilisateur déjà membre → action ignorée |

---

### UC20 — Consulter le classement d'un groupe *(généralisé par UC17)*

| Champ | Contenu |
|---|---|
| **résumé** | Cas particulier de UC17, restreint aux membres d'un groupe |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | — |
| **pré-conditions** | L'utilisateur est membre du groupe |
| **résultats** | Affichage du classement filtré sur les membres du groupe |
| **description** | 1. L'utilisateur sélectionne un groupe dont il est membre<br>2. Le système affiche le classement restreint aux membres de ce groupe (même logique de tri que UC17) |
| **exceptions** | Utilisateur non membre du groupe → accès refusé |

---

### UC21 — Envoyer une demande d'ami

| Champ | Contenu |
|---|---|
| **résumé** | L'utilisateur propose une relation d'ami à un autre utilisateur |
| **acteur primaire** | Utilisateur |
| **acteur secondaire** | Service d'email (prestataire Resend) ; Utilisateur destinataire (notifié) |
| **pré-conditions** | Utilisateur connecté ; destinataire existant |
| **résultats** | Une demande d'ami est envoyée, en attente de réponse |
| **description** | 1. L'utilisateur sélectionne un autre utilisateur<br>2. Il envoie une demande d'ami<br>3. Le système enregistre la demande<br>4. Le système notifie le destinataire par email via le **Service d'email** |
| **exceptions** | Demande déjà envoyée ou déjà amis → action ignorée |

---

### UC22 — Traiter un signalement

| Champ | Contenu |
|---|---|
| **résumé** | Le modérateur examine un témoignage signalé et décide d'une action |
| **acteur primaire** | Modérateur |
| **acteur secondaire** | Service d'email (prestataire Resend) ; Utilisateur (auteur du témoignage, notifié) |
| **pré-conditions** | Au moins un signalement est en attente (créé via **UC13**) |
| **résultats** | Le signalement est clos ; le témoignage est conservé ou supprimé |
| **description** | 1. Le modérateur consulte la liste des signalements en attente<br>2. Il examine le témoignage signalé et le motif<br>3. Il décide : rejeter le signalement ou supprimer le témoignage<br>4. Le système applique la décision et notifie l'auteur si besoin, par email via le **Service d'email** |
| **exceptions** | Signalement déjà traité → action ignorée |
