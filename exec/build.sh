#!/usr/bin/env sh

error() {

  last_error_exit_code="$?"

  if [[ ${allow_next_error} == "true" ]] ; then
    allow_next_error=false
    return 0
  fi

  local parent_lineno="$1"
  error_message="$2"
  local code="${3:-1}"

  if [[ -z "$error_message" ]] ; then
    error_message="Exiting with status ${code}"
  fi

  echo "Error on or near line ${parent_lineno}: ${error_message}" 1>&2

  exit "${code}"
}

trap 'error ${LINENO}' ERR

base_path="$(pwd)"
project_name="$(basename ${base_path})"

echo "Building ${project_name}..."

build_path="${base_path}/build"

if [[ -d "${build_path}" ]] ; then
  rm -rf "${build_path}"
fi

mkdir "${build_path}"

version=$(git describe --tags)
build_date=$(date)
git_commit=$(git rev-parse --short HEAD)

echo "Version: ${version}"
echo "BuildDate: ${build_date}"
echo "GitCommit: ${git_commit}"

pkg="$(head -n 1 go.mod | cut -d' ' -f2)"
pkg_version="${pkg}/version"

go build -ldflags "-X '${pkg_version}.Version=${version}' \
                   -X '${pkg_version}.BuildDate=${build_date}' \
                   -X '${pkg_version}.GitCommit=${git_commit}' \
         -o "${build_path}/metar-ws2811" \
         ./main.go

echo "Done."
