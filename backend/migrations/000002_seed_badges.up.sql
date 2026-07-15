-- Jeu de badges par défaut (UC14/UC15). La fiche ne fournit aucun barème :
-- choix arbitraire mais raisonnable, seuils croissants sur le score
-- dénormalisé de l'utilisateur (cf. commentaire 000001 et
-- internal/services/scoring.go pour le détail du calcul du score).
INSERT INTO badges (code, nom, description, seuil_score) VALUES
    ('premiere-cuite',          'Première Cuite',          'Première soirée enregistrée avec succès.', 10),
    ('habitue-du-bar',          'Habitué du Bar',           'Score cumulé de 50 points ou plus.',        50),
    ('legende-de-la-soiree',    'Légende de la Soirée',     'Score cumulé de 150 points ou plus.',       150),
    ('roi-de-la-gueule-de-bois','Roi de la Gueule de Bois', 'Score cumulé de 300 points ou plus.',       300);
