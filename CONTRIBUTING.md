# Convention de branches

Format : `type(scope):description`

Types autorisés :
- `feat` — nouvelle fonctionnalité
- `fix` — correction de bug
- `style` — changement visuel/CSS sans impact fonctionnel

Exemples :
- `feat(frontend):homepage`
- `fix(api):videos`

## Workflow
1. Créer une branche depuis `dev`
2. Développer, commit, push
3. Ouvrir une Pull Request vers `dev`
4. Une fois `dev` stable → PR vers `preprod`
5. Une fois validé en préprod → PR vers `main`
