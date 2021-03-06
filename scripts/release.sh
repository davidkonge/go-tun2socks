#!/usr/bin/env bash

# Create GITHUB_API_TOKEN and export it first.

set -x

tag=v1.6
description="Adding routes before packets input to lwIP."
list_assets_cmd="ls -1 build/*.zip"

token=$GITHUB_API_TOKEN
owner=eycorsican
repo=go-tun2socks
base_url=https://api.github.com

content_type_json="Content-Type: application/json"
content_type_zip="Content-Type: application/zip"

api_create_release=/repos/$owner/$repo/releases
api_get_release_by_tag=/repos/$owner/$repo/releases/tags/$tag
api_create_release_data="{\
\"tag_name\": \"${tag}\",\
\"target_commitish\": \"master\",\
\"name\": \"${tag}\",
\"body\": \"${description}\",
\"draft\": false,\
\"prerelease\": false\
}"

# Get the release id by tag name.
release_id=`curl -u $owner:$token -H "$content_type_json" -X GET "${base_url}${api_get_release_by_tag}" | \
python -c "import sys;import json;print(json.loads(\"\".join(sys.stdin.readlines()))[\"id\"])"` 2>/dev/null

# If there is release with the tag exists, delete it first.
if [ $? -eq 0 ]; then
	api_delete_release=/repos/$owner/$repo/releases/$release_id
    curl -u $owner:$token -H "$content_type_json" -X DELETE "${base_url}${api_delete_release}"
fi

# Create a release.
curl -u $owner:$token -H "$content_type_json" -X POST "${base_url}${api_create_release}" --data "${api_create_release_data}"

# Get the release id by tag name.
release_id=`curl -u $owner:$token -H "$content_type_json" -X GET "${base_url}${api_get_release_by_tag}" | \
python -c "import sys;import json;print(json.loads(\"\".join(sys.stdin.readlines()))[\"id\"])"`

# Upload assets.
eval $list_assets_cmd | while read -r asset_name; do
	upload_url=https://uploads.github.com/repos/$owner/$repo/releases/$release_id/assets?name=$(basename $asset_name)
	curl --progress-bar -u $owner:$token -H "$content_type_zip" --data-binary @"$asset_name" -X POST $upload_url
	if [ $? -ne "0" ]; then
		exit 1
	fi
done
