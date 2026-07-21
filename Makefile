# Fusion root Makefile.
# Targets here are convenience wrappers used by CI and local dev.
.PHONY: verify-fresh-bundle web-build web-test web-typecheck

web-build:
	npm run build -w @dcraft-fusion/web

web-test:
	npm test -w @dcraft-fusion/web

web-typecheck:
	npm run typecheck -w @dcraft-fusion/web

# Fails if apps/web/dist is missing or older than the newest source file
# under apps/web/src or packages/domain/src. Used by CI to guarantee the
# SPA bundle is always rebuilt from source (never shipped stale).
verify-fresh-bundle:
	@test -d apps/web/dist || { echo "::error:: apps/web/dist missing — run npm run build -w @dcraft-fusion/web" >&2; exit 1; }
	@NEWEST_SRC="$$(find apps/web/src packages/domain/src -type f -printf '%T@ %p\n' | sort -nr | head -1 | cut -d' ' -f2-)"; \
	NEWEST_DIST="$$(find apps/web/dist -type f -printf '%T@ %p\n' | sort -nr | head -1 | cut -d' ' -f2-)"; \
	if [ -n "$$NEWEST_SRC" ] && [ -n "$$NEWEST_DIST" ] && \
	   [ "$$(stat -c %Y "$$NEWEST_DIST")" -lt "$$(stat -c %Y "$$NEWEST_SRC")" ]; then \
	  echo "::error:: apps/web/dist is older than newest src/ — bundle not rebuilt" >&2; exit 1; \
	fi; \
	echo "OK: apps/web/dist is fresh relative to src/"
