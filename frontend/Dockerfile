# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the rest of the code
COPY . .

# Add build arguments for Dynamic URL config
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

# Build the Vite project
RUN npm run build

# Final stage - Nginx
FROM nginx:alpine

# Copy the build output to replace the default nginx contents.
COPY --from=builder /app/dist /usr/share/nginx/html

# Replace default config with our proxy config template
# We use a template so envsubst can inject variables at runtime
COPY nginx.conf /etc/nginx/templates/default.conf.template

EXPOSE 80

# Nginx alpine image has built-in support for templates in /etc/nginx/templates/
# It will automatically run envsubst and output to /etc/nginx/conf.d/
CMD ["nginx", "-g", "daemon off;"]
