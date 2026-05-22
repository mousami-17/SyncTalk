# Build stage
FROM node:22-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install ALL dependencies (including devDependencies needed for build)
RUN npm install

# Copy source code
COPY . .

# Build arguments for environment variables
ARG VITE_API_URL
ARG VITE_WS_HOST
ARG VITE_WS_PORT

# Create .env file for build
RUN echo "VITE_API_URL=${VITE_API_URL}" > .env && \
    echo "VITE_WS_HOST=${VITE_WS_HOST}" >> .env && \
    echo "VITE_WS_PORT=${VITE_WS_PORT}" >> .env

# Build the application
RUN npm run build

# Production stage - serve with nginx
FROM nginx:1.27-alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy built files from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
