#! /bin/bash

# Function to display usage
usage() {
  echo "Usage: docker_build.sh [options]"
  echo ""
  echo "options:"
  echo "-p | --push: call push_images() function"
  echo "-t [tag]: tag (default: local)"
  exit 1
}

# Build images

function build {
  echo "Building Docker images with tag: $TAG"
  docker build -t "gateway/edv:$TAG" -f edv.Dockerfile ../.. || exit 1
  docker build  -t "gateway/master:$TAG" -f master.Dockerfile ../.. || exit 1
}

function build_multiplatform {
  docker build --platform linux/amd64,linux/arm64 -t "gateway/edv:$TAG" -f edv.Dockerfile ../.. || exit 1
  docker build --platform linux/amd64,linux/arm64  -t "gateway/master:$TAG" -f master.Dockerfile ../.. || exit 1
}

# Tag and push images if the tag is not "local"
function push_images {
  echo "[1] Build multiplatform images to tag and push the images with tag: $TAG"
  echo "[2] Push images with tag: $TAG"
  read -p "Do you want to continue? (y/n) " choice

  build_multiplatform
  case "$choice" in
    y|Y )
      echo "Tagging and pushing images to remote repository with tag: $TAG"
      docker image tag "gateway/edv:$TAG" "us-docker.pkg.dev/gateway-protocol/dfs/edv:$TAG"
      docker push "us-docker.pkg.dev/gateway-protocol/dfs/edv:$TAG"

      docker image tag "gateway/master:$TAG" "us-docker.pkg.dev/gateway-protocol/dfs/master:$TAG"
      docker push "us-docker.pkg.dev/gateway-protocol/dfs/master:$TAG"
      ;;
    n|N )
      echo "Tagging and pushing process aborted."
      ;;
    * )
      echo "Invalid input. Tagging and pushing process aborted."
      ;;
  esac

  echo "Docker image build and push process completed."
}

# Default tag name
TAG="local"

# Parse options
while [[ "$1" != "" ]]; do
  should_push=1
    case $1 in
      -p | --push )
        should_push=0
        ;;
      -t )
        shift
        if [[ -z "$1" || "$1" == -* ]]; then
          echo "Error: -t option requires an argument."
          usage
        fi
        TAG="$1"
        echo "Tag is set to $TAG"
        ;;
      --tag=* )
        TAG="${1#*=}"
        echo "Tag is set to $TAG"
        ;;
      -h | --help )
        usage
        ;;
      * )
        echo "Error: Invalid option '$1'"
        usage
        ;;
    esac
  shift
done

test $(test $should_push && push_images) || build