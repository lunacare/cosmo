# Use Node.js LTS version
FROM 836236105554.dkr.ecr.us-west-2.amazonaws.com/mirror:node-22-alpine

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Set working directory
WORKDIR /app

# Copy package files first to leverage Docker cache
COPY package.json pnpm-lock.yaml ./

# Install dependencies --frozen-lockfile
RUN pnpm install

# Copy source code and configuration files
COPY . .

# Run tests
CMD ["pnpm", "test"] 