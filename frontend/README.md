# La Gueule de Bois — Frontend

SPA React + Vite + TypeScript. Interface en français. Thème visuel « Ambrée »
(voir `docs/design/theme-ambree.css`, copié dans `src/theme/theme-ambree.css`).

## Installation

```bash
npm install
```

## Développement

```bash
npm run dev
```

Sert l'app sur `http://localhost:5173`. Les requêtes vers `/api` et `/uploads`
sont proxifiées vers `http://localhost:8080` (API Go en local, cf.
`vite.config.ts`) — démarrer le backend séparément.

Tant que le backend n'est pas démarré, les pages affichent leur état vide
(pas d'erreur bloquante), l'auth reste sur l'écran de connexion.

## Build

```bash
npm run build
```

Vérifie les types (`tsc -b`) puis produit le bundle statique dans `dist/`
(prévu pour être servi par nginx, cf. spec de déploiement).

## Lint

```bash
npm run lint
```

## Structure

```
src/
  api/         client fetch (`client.ts`), types des entités (`types.ts`),
               un module par ressource (auth, soirees, temoignages, classement,
               groupes, badges, users)
  auth/        AuthContext (session cookie httpOnly, /api/auth/*)
  components/  Layout (header + nav), ProtectedRoute, états loading/erreur/vide
  pages/       une page par écran (voir App.tsx pour la liste des routes et
               les TODO des cas d'utilisation pas encore couverts)
  theme/       tokens CSS « Ambrée » (copie de docs/design/theme-ambree.css)
```

Le swipe (UC12) n'est volontairement pas implémenté dans ce lot.
