# La Gueule de Bois — Liste finale : 22 cas d'utilisation

## Acteurs

**Primaires** (généralisation) : `Visiteur ← Utilisateur ← Modérateur`

Pas d'acteur « Système de scoring » : ce n'est pas une entité externe au sens du cours, juste un traitement interne. UC16 et UC14 n'ont donc aucun acteur associé — ils ne sont atteints que par `include`.

---

## Liste des 22 cas d'utilisation

| ID | Cas d'utilisation | Acteur(s) | Relations |
|---|---|---|---|
| UC01 | S'inscrire | Visiteur | — |
| UC02 | Se connecter | Visiteur | — |
| UC03 | Se déconnecter | Utilisateur | — |
| UC04 | Modifier son profil | Utilisateur | — |
| UC05 | Consulter le profil d'un autre utilisateur | Utilisateur | — |
| UC06 | Créer une soirée | Utilisateur | **Inclut** UC16 |
| UC07 | Modifier une soirée | Utilisateur | **Inclut** UC16 |
| UC08 | Supprimer une soirée | Utilisateur | **Inclut** UC16 |
| UC09 | Inviter un témoin à une soirée | Utilisateur | — |
| UC10 | Consulter l'historique de ses soirées | Utilisateur | — |
| UC11 | Ajouter un témoignage sur une soirée | Utilisateur | **Inclut** UC16 |
| UC12 | Swiper/voter sur un témoignage | Utilisateur | — |
| UC13 | Signaler un témoignage | Utilisateur | *(voir postcondition : déclenche UC22, pas de relation UML formelle)* |
| UC14 | Débloquer un badge | *(aucun, cas interne)* | **Inclut par** UC16 |
| UC15 | Consulter ses badges | Utilisateur | — |
| UC16 | Recalculer le score d'un utilisateur | *(aucun, cas interne)* | *(cas système)* — **Inclut** UC14 |
| UC17 | Consulter le classement global | Utilisateur | — |
| UC18 | Créer un groupe | Utilisateur | — |
| UC19 | Rejoindre un groupe | Utilisateur | — |
| UC20 | Consulter le classement d'un groupe | Utilisateur | **Généralisation** : UC20 est un cas particulier de UC17 |
| UC21 | Envoyer une demande d'ami | Utilisateur | — |
| UC22 | Traiter un signalement | Modérateur | déclenché par UC13 *(lien narratif, sans relation UML formelle)* |

---

## Ce qui a été volontairement coupé (par rapport à la liste à 38)

Retiré pour garder la liste simple, mais réintégrable facilement en bonus si vous avez du temps en fin de projet : visibilité de profil, suppression de compte, commentaires sur témoignage, filtrage du classement par période, exclusion de membre de groupe, quitter un groupe, accepter/refuser une demande d'ami, notifications détaillées, bannissement, recherche, paramètres de notification.

## Relations UML à représenter sur le diagramme

- **Généralisation (acteurs)** : `Visiteur ← Utilisateur ← Modérateur`
- **Généralisation (cas d'utilisation)** : UC20 (Consulter le classement d'un groupe) généralisé par UC17 (Consulter le classement global) — même logique que l'exemple du cours « Consulter sur Internet » spécialise « Consulter comptes »
- **Inclusions** : UC06, UC07, UC08, UC11 → incluent tous UC16 (le score se recalcule à chaque action sur une soirée/témoignage) ; UC16 → inclut UC14 (le badge est débloqué automatiquement après recalcul, sans acteur déclencheur externe)
- **Pas de relation UML entre UC13 et UC22** : `extend` sert à insérer un comportement optionnel dans la *même* interaction (même acteur), ce qui n'est pas le cas ici (deux acteurs différents, deux moments différents). Le lien « signaler peut mener à un traitement » est documenté dans la **postcondition textuelle de UC13**, pas comme relation graphique.

Cette version est plus lisible pour un diagramme propre, tout en gardant assez de relations (inclusion/généralisation) pour satisfaire l'exigence du sujet.

## Prochaine étape

On peut maintenant attaquer le **diagramme de classes** (7+ classes, attributs, méthodes, relations), qui doit rester cohérent avec ces 22 cas. Prêt à enchaîner ?
