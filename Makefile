.PHONY: schema schema-check

schema:
	go run ./tools/schema/ -o system/schema/schema.sql

schema-check:
	@tmp=$$(mktemp); \
	go run ./tools/schema/ -o "$$tmp"; \
	if ! diff -q "$$tmp" system/schema/schema.sql > /dev/null 2>&1; then \
		echo "ERROR: system/schema/schema.sql is out of date. Run 'make schema'."; \
		rm "$$tmp"; \
		exit 1; \
	fi; \
	rm "$$tmp"; \
	echo "OK: schema is up to date."
