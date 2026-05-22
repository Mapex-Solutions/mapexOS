# Stage 1: Builder
FROM node:24-alpine AS builder

WORKDIR /build

# isolated-vm requires native compilation
RUN apk add --no-cache python3 make g++

# Copy workspace root files
COPY package.json package-lock.json ./
COPY configs/ configs/

# Copy required workspace packages
COPY packages/utils/ packages/utils/
COPY packages/infrastructure/ packages/infrastructure/
COPY packages/microservices/ packages/microservices/
COPY packages/script-registry/ packages/script-registry/
COPY packages/validations/ packages/validations/
COPY packages/schemas/ packages/schemas/

# Copy service
COPY services/js-workflow-executor/ services/js-workflow-executor/

# Install all dependencies (workspace resolution)
RUN npm ci

# Build packages (they export from lib/ which is gitignored)
RUN npm run build --workspace=packages/utils && \
    npm run build --workspace=packages/validations && \
    npm run build --workspace=packages/schemas && \
    npm run build --workspace=packages/infrastructure && \
    npm run build --workspace=packages/microservices && \
    npm run build --workspace=packages/script-registry

# Build service
RUN npm run build --workspace=services/js-workflow-executor

# Prune dev dependencies
RUN npm prune --omit=dev

# Replace workspace symlinks with actual built packages
RUN cd node_modules/@mapexos && \
    for pkg in utils infrastructure microservices script-registry validations schemas; do \
        rm -rf $pkg && mkdir -p $pkg && \
        cp -r /build/packages/$pkg/lib $pkg/lib && \
        cp /build/packages/$pkg/package.json $pkg/package.json; \
    done

# Stage 2: Production
FROM node:24-alpine

WORKDIR /app

RUN apk add --no-cache tini

# Copy built service and resolved node_modules
COPY --from=builder /build/services/js-workflow-executor/dist ./dist
COPY --from=builder /build/services/js-workflow-executor/package.json ./
COPY --from=builder /build/node_modules ./node_modules

# Non-root user (node:24-alpine already has node user with uid 1000)
RUN mkdir -p /tmp/mapexos/cache && chown -R node:node /app /tmp/mapexos
USER node

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["node", "--no-node-snapshot", "dist/main.js"]
