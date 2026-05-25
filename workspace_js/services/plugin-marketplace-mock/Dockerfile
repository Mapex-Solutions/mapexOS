# =============================================================================
# Plugin Marketplace Mock — single-stage Node image.
#
# Mock CDN that serves the plugin registry + per-plugin manifests over
# HTTP. Runs as a tsx-managed TypeScript process; no transpile step is
# required. The build context is this service's own folder.
#
# Usage (from this service folder):
#   docker build -t mapexos/plugin-marketplace-mock:dev .
# =============================================================================

FROM node:24-alpine

WORKDIR /app

# Install dependencies first so this layer caches across source-only
# changes. The repo doesn't keep a package-lock for this mock, so npm
# install is used instead of npm ci.
COPY package.json ./
RUN npm install --no-audit --no-fund

# Source + static plugin assets + tsconfig (tsx resolves the
# @modules / @shared aliases at runtime from tsconfig.json).
COPY tsconfig.json ./
COPY src/ ./src/
COPY public/ ./public/

# Non-root, mirror the node user that the official image ships.
RUN chown -R node:node /app
USER node

EXPOSE 3099

CMD ["npm", "start"]
