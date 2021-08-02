
BEGIN;

CREATE TABLE IF NOT EXISTS lsif_dependency_repos (
    name text NOT NULL,
    version text NOT NULL,
    scheme text NOT NULL,
    PRIMARY KEY (name, version, scheme)
);

COMMIT;
