-- UC22 : traiter un signalement en supprimant le témoignage signalé
-- déclenchait une suppression en cascade du signalement lui-même (ON DELETE
-- CASCADE), empêchant le MarkTraite qui suit de le retrouver (0 ligne
-- affectée -> 404) et perdant l'historique de modération. Le signalement
-- doit survivre à la suppression du témoignage qu'il visait.
ALTER TABLE signalements ALTER COLUMN temoignage_id DROP NOT NULL;
ALTER TABLE signalements DROP CONSTRAINT signalements_temoignage_id_fkey;
ALTER TABLE signalements
    ADD CONSTRAINT signalements_temoignage_id_fkey
    FOREIGN KEY (temoignage_id) REFERENCES temoignages(id) ON DELETE SET NULL;
