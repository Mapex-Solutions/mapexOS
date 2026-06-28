# =============================================================================
# MapexOS Frontend — Multi-stage Build
# =============================================================================
#
# Build context: workspace_js/  (the monorepo JS root)
#
# Runtime env vars override __PLACEHOLDER__ strings via entrypoint.
#
# Usage:
#   docker build -t mapexos/frontend:latest -f apps/mapexOS/Dockerfile .
#
# =============================================================================

# Stage 1: Builder
FROM node:24-alpine AS builder

WORKDIR /build

# Copy workspace root files
COPY package.json package-lock.json ./
COPY configs/ configs/

# Copy required workspace packages (apis → schemas, validations, utils)
COPY packages/utils/ packages/utils/
COPY packages/validations/ packages/validations/
COPY packages/schemas/ packages/schemas/
COPY packages/apis/ packages/apis/

# Copy stub package.json for workspaces not needed but referenced in lockfile
# This prevents npm ci from failing on missing workspace references
COPY packages/infrastructure/package.json packages/infrastructure/package.json
COPY packages/microservices/package.json packages/microservices/package.json
COPY packages/script-registry/package.json packages/script-registry/package.json
RUN mkdir -p packages/infrastructure/src packages/microservices/src packages/script-registry/src

# Copy the frontend app
COPY apps/mapexOS/ apps/mapexOS/

# Install all dependencies (workspace resolution)
RUN npm ci

# Build workspace packages that the app imports from lib/
RUN npm run build --workspace=packages/utils && \
    npm run build --workspace=packages/validations && \
    npm run build --workspace=packages/schemas && \
    npm run build --workspace=packages/apis

# Set API URL placeholders (replaced at runtime by entrypoint)
ARG MAPEXOS_API_BASE_URL=__MAPEXOS_API_BASE_URL__
ARG HTTP_GATEWAY_API_BASE_URL=__HTTP_GATEWAY_API_BASE_URL__
ARG ASSETS_API_BASE_URL=__ASSETS_API_BASE_URL__
ARG ROUTER_API_BASE_URL=__ROUTER_API_BASE_URL__
ARG EVENTS_API_BASE_URL=__EVENTS_API_BASE_URL__
ARG TRIGGERS_API_BASE_URL=__TRIGGERS_API_BASE_URL__
ARG WORKFLOW_API_BASE_URL=__WORKFLOW_API_BASE_URL__
ARG VAULT_API_BASE_URL=__VAULT_API_BASE_URL__
ARG JS_EXECUTOR_API_BASE_URL=__JS_EXECUTOR_API_BASE_URL__

ENV MAPEXOS_API_BASE_URL=$MAPEXOS_API_BASE_URL \
    HTTP_GATEWAY_API_BASE_URL=$HTTP_GATEWAY_API_BASE_URL \
    ASSETS_API_BASE_URL=$ASSETS_API_BASE_URL \
    ROUTER_API_BASE_URL=$ROUTER_API_BASE_URL \
    EVENTS_API_BASE_URL=$EVENTS_API_BASE_URL \
    TRIGGERS_API_BASE_URL=$TRIGGERS_API_BASE_URL \
    WORKFLOW_API_BASE_URL=$WORKFLOW_API_BASE_URL \
    VAULT_API_BASE_URL=$VAULT_API_BASE_URL \
    JS_EXECUTOR_API_BASE_URL=$JS_EXECUTOR_API_BASE_URL

# Build the Quasar SPA
WORKDIR /build/apps/mapexOS
RUN npx quasar build --mode spa

# Stage 2: Production (nginx)
FROM nginx:1.27-alpine

# Copy nginx config
COPY apps/mapexOS/nginx.conf /etc/nginx/conf.d/default.conf

# Copy entrypoint script
COPY apps/mapexOS/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

# Copy built SPA
COPY --from=builder /build/apps/mapexOS/dist/spa/ /usr/share/nginx/html/

EXPOSE 80

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
