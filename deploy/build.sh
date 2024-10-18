// Copyright 2025 Confidential Containers.
// This file is part of Confidential Containers

// Licensed under the Apache 2.0;
//
// https://github.com/confidential-containers/operator
//


#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

coco_containerd_version=1.6.8.2
official_containerd_version=1.7.7
vfio_gpu_containerd_version=1.7.0.0
nydus_snapshotter_version=v0.13.14

coco_containerd_repo=${coco_containerd_repo:-"https://github.com/confidential-containers/containerd"}
official_containerd_repo=${official_containerd_repo:-"https://github.com/containerd/containerd"}
vfio_gpu_containerd_repo=${vfio_gpu_containerd_repo:-"https://github.com/confidential-containers/containerd"}
nydus_snapshotter_repo=${nydus_snapshotter_repo:-"https://github.com/containerd/nydus-snapshotter"}
containerd_dir="$(mktemp -d -t containerd-XXXXXXXXXX)/containerd"
extra_docker_manifest_flags="${extra_docker_manifest_flags:-}"
registry="wetee/cvm-deploy"

script_dir=$(dirname "$(readlink -f "$0")")

supported_arches=(
	"linux/amd64"
	# "linux/s390x"
)

function setup_env_for_arch() {
	case "$1" in
		"linux/amd64")
			kernel_arch="x86_64"
			golang_arch="amd64"
			;;
		"linux/s390x")
			kernel_arch="s390x"
			golang_arch="s390x"
			;;
		*) echo "$1 is not supported" >/dev/stderr && exit 1 ;;
	esac
}

function main() {
	pushd "${script_dir}"
	local tag

	tag=$(git rev-parse HEAD)

	for arch in "${supported_arches[@]}"; do
		setup_env_for_arch "${arch}"

		echo "Building containerd payload image for ${arch}"
		docker build \
			--build-arg ARCH="${golang_arch}" \
			--build-arg COCO_CONTAINERD_VERSION="${coco_containerd_version}" \
			--build-arg COCO_CONTAINERD_REPO="${coco_containerd_repo}" \
			--build-arg OFFICIAL_CONTAINERD_VERSION="${official_containerd_version}" \
			--build-arg OFFICIAL_CONTAINERD_REPO="${official_containerd_repo}" \
			--build-arg VFIO_GPU_CONTAINERD_VERSION="${vfio_gpu_containerd_version}" \
			--build-arg VFIO_GPU_CONTAINERD_REPO="${vfio_gpu_containerd_repo}" \
			--build-arg NYDUS_SNAPSHOTTER_VERSION="${nydus_snapshotter_version}" \
			--build-arg NYDUS_SNAPSHOTTER_REPO="${nydus_snapshotter_repo}" \
			-t "${registry}:${kernel_arch}-${tag}" \
			--platform="${arch}" \
			.
		docker push "${registry}:${kernel_arch}-${tag}"
	done

	popd
}

main "$@"