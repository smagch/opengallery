
dump:
	pg_dump --no-tablespaces --no-owner --clean -s galleryinfo > database/dump.sql

.PHONY: dump