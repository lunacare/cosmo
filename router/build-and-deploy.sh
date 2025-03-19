#!/usr/bin/env bash
set -euo pipefail

# Get the root directory of the repository
ROOT_DIR=$(git rev-parse --show-toplevel)
cd "$ROOT_DIR"

if [[ $# -ne 1 ]]; then
        echo "1 --- BUILD_AS_TEST - 'true' or 'false'"
        exit 1
fi
BUILD_AS_TEST=$1

ECR_REPO="lunacare-cosmo-router"
ECR_URL="836236105554.dkr.ecr.us-west-2.amazonaws.com/$ECR_REPO"

VERSION=$( grep -Eo '\[[0-9]+\.[0-9]+\.[0-9]+\]' ./router/CHANGELOG.md | tr -d '[]' | sort -V | tail -n1 )
VERSION="$VERSION-lunacare-$(head -n 1 ./router/VERSION-LUNACARE)"
VERSION_TAG="$ECR_URL:$VERSION"
LATEST_TAG="$ECR_URL:latest"
LATEST_TEST_TAG="$ECR_URL:latest-test"

GITHUB_ACTIONS="${GITHUB_ACTIONS:-"false"}"
if [[ "$GITHUB_ACTIONS" == "false" ]]; then
        CMD_DESCRIBE_IMAGES="aws ecr describe-images --repository-name $ECR_REPO --image-ids imageTag=$VERSION --profile live"
else
        CMD_DESCRIBE_IMAGES="aws ecr describe-images --repository-name $ECR_REPO --image-ids imageTag=$VERSION"
fi

# Execute the command and capture the error message
echo "running command: $CMD_DESCRIBE_IMAGES"
if output=$(eval "$CMD_DESCRIBE_IMAGES" 2>&1); then
    echo "VERSION $VERSION_TAG already exists. Exiting without building images"
    exit 1
else
    # Check for specific error messages
    if [[ $output == *"ImageNotFoundException"* ]]; then
        echo "Tag for VERSION $VERSION not found, proceeding to build it."
    elif [[ $output == *"InvalidParameterException"* ]]; then
        echo "Error: Invalid parameter when describing ecr images $VERSION."
        exit 1
    else
        echo "An unexpected error occurred: $output"
        exit 1
    fi
fi

if [ "$BUILD_AS_TEST" == "true" ]; then
        TAG_ARGS="-t ${LATEST_TEST_TAG}"
        TAGS_TO_PUSH="${LATEST_TEST_TAG}"
else
        TAG_ARGS="-t ${LATEST_TAG} -t ${VERSION_TAG}"
        TAGS_TO_PUSH="${LATEST_TAG} ${VERSION_TAG}"
fi

IFS=' ' read -r -a TAGS_TO_PUSH_ARRAY <<< "$TAGS_TO_PUSH"
for TAG in "${TAGS_TO_PUSH_ARRAY[@]}"; do
        echo "Set tag to push: $TAG"
done

export ECR_REPO="$ECR_REPO"
export TAG_ARGS="$TAG_ARGS"
export TAGS_TO_PUSH="$TAGS_TO_PUSH"

IFS=' ' read -r -a TAG_ARGS_ARRAY <<< "$TAG_ARGS"

# Create a new builder instance if it doesn't exist
BUILDER_NAME="cosmo-multiplatform-builder"

# Remove existing builder if it exists
docker buildx rm "${BUILDER_NAME}" || true

# Create new builder and bootstrap it
docker buildx create --name "${BUILDER_NAME}" --driver docker-container --bootstrap

# Use the builder
docker buildx use "${BUILDER_NAME}"
docker buildx inspect --bootstrap

COMMIT_SHA=$(git rev-parse HEAD)
DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Build and push the images
docker buildx build \
        --platform linux/arm64 \
        --build-arg TARGETOS=linux \
        --build-arg TARGETARCH=arm64 \
        --build-arg COMMIT_SHA=${COMMIT_SHA} \
        --build-arg DATE=${DATE} \
        --build-arg VERSION=${VERSION} \
        "${TAG_ARGS_ARRAY[@]}" \
        --push \
        -f ./router/custom-luna.Dockerfile \
        --progress plain \
        ./router

docker buildx build \
        --platform linux/amd64 \
        --build-arg TARGETOS=linux \
        --build-arg TARGETARCH=amd64 \
        --build-arg COMMIT_SHA=${COMMIT_SHA} \
        --build-arg DATE=${DATE} \
        --build-arg VERSION=${VERSION} \
        "${TAG_ARGS_ARRAY[@]}" \
        --push \
        -f ./router/custom-luna.Dockerfile \
        --progress plain \
        ./router

