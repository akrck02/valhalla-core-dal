# !/bin/bash
static_valhalla_url=github.com/akrck02
project=$1
project="valhalla-$project"
version=$2

# check if the project is empty
if [ -z "$project" ]; then
    echo "Project name is required"
    exit 1
fi

# check if the version is empty
if [ -z "$version" ]; then
    echo "Version is required"
    exit 1
fi

# download the project dependency
echo "Downloading $project@$version from $static_valhalla_url"
go get $static_valhalla_url/$project@$version