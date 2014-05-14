
dump:
	pg_dump --no-tablespaces --clean -x --no-owner -s galleryinfo > database/dump.sql

clean:
	rm -f opengallery

.PHONY: dump clean