
dump:
	pg_dump --no-tablespaces --clean -x --no-owner -s galleryinfo > database/dump.sql

clean:
	rm -f open-gallery-info

.PHONY: dump clean