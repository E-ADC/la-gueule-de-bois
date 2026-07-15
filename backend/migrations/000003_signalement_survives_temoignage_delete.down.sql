ALTER TABLE signalements DROP CONSTRAINT signalements_temoignage_id_fkey;
ALTER TABLE signalements
    ADD CONSTRAINT signalements_temoignage_id_fkey
    FOREIGN KEY (temoignage_id) REFERENCES temoignages(id) ON DELETE CASCADE;
ALTER TABLE signalements ALTER COLUMN temoignage_id SET NOT NULL;
