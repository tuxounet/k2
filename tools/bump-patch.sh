
#!/bin/bash
#set -e

# Check for changes
if git diff-index --quiet HEAD --; then
    echo "No changes to commit."
    exit 0
fi


# Fetch the latest tag from GitLab
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

# Extract the major, minor, and patch version numbers
IFS='.' read -r -a version_parts <<< "$latest_tag"
major=${version_parts[0]}
minor=${version_parts[1]}
patch=${version_parts[2]}

# Increment the patch version
new_patch=$((patch + 1))

# Create the new version tag
new_tag="$major.$minor.$new_patch"


echo $new_tag > version.txt

# Add all changes
git add .

# Commit with an automatic message
commit_message="Bump version to $new_tag"
git commit -m "$commit_message"

# Push changes
git push origin $(git rev-parse --abbrev-ref HEAD)

echo "Changes committed and pushed with message: $commit_message"


# Tag the new version in Git
git tag $new_tag
git push origin $new_tag

echo "New version tag created: $new_tag"

