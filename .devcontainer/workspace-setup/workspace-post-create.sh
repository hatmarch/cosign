#!/bin/bash

set -euo pipefail

WORKSPACE_FOLDER=${1}
#SSO_PROFILE=${2}

rsync -a .devcontainer/workspace-setup/vscode/ ${WORKSPACE_FOLDER}/.vscode/ --ignore-existing

echo "Logging into AWS SSO with profile ${SSO_PROFILE:-container}"
aws sso login --profile ${SSO_PROFILE:-container}

#sudo pip install -r requirements.txt