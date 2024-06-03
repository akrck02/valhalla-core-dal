date=$(date +"%Y%m%d%H%M%S")
version=$(cat version.config)
branch=$(git rev-parse --abbrev-ref HEAD)

versionstring="$version-$date-$branch"
echo "Publishing version $versionstring"

# Build the project
git tag -a $versionstring -m "Release $versionstring"
git push origin --tags