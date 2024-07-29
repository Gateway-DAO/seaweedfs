#! /bin/bash

#!/bin/bash

# Function to display usage information
usage() {
  echo "Usage: $0 [-t TAG_NAME]"
  echo "Options:"
  echo "  -t TAG_NAME    Custom tag name for the Docker images (default: local)"
  echo "  -h             Show this help message"
  exit 1
}

# Default tag name
TAG_NAME="local"

# Parse command-line arguments
while getopts ":t:h" opt; do
  case $opt in
    t)
      TAG_NAME="$OPTARG"
      ;;
    h)
      usage
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      usage
      ;;
  esac
done

# Build images
echo "Building Docker images with tag: $TAG_NAME"
docker build --platform linux/amd64,linux/arm64 -t "gateway/edv:$TAG_NAME" -f edv.Dockerfile ../.. || exit 1
docker build --platform linux/amd64,linux/arm64  -t "gateway/master:$TAG_NAME" -f master.Dockerfile ../.. || exit 1

# Tag and push images if the tag is not "local"
if [ "$TAG_NAME" != "local" ]; then
  echo "You are about to tag and push the images with tag: $TAG_NAME"
  read -p "Do you want to continue? (y/n) " choice
  case "$choice" in
    y|Y )
      echo "Tagging and pushing images to remote repository with tag: $TAG_NAME"
      docker image tag "gateway/edv:$TAG_NAME" "us-docker.pkg.dev/gateway-protocol/dfs/edv:$TAG_NAME"
      docker push "us-docker.pkg.dev/gateway-protocol/dfs/edv:$TAG_NAME"

      docker image tag "gateway/master:$TAG_NAME" "us-docker.pkg.dev/gateway-protocol/dfs/master:$TAG_NAME"
      docker push "us-docker.pkg.dev/gateway-protocol/dfs/master:$TAG_NAME"
      ;;
    n|N )
      echo "Tagging and pushing process aborted."
      ;;
    * )
      echo "Invalid input. Tagging and pushing process aborted."
      ;;
  esac
fi

echo "Docker image build and push process completed."
