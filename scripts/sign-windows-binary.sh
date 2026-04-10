#!/usr/bin/env bash

set -euo pipefail

if [[ "$SKIP_SIGNING" == "true" ]]; then
    >&2 echo "Skipping signing of windows binary as SKIP_SIGNING is set";
    exit 0;
fi

# Only sign windows binary if fully configured.
# Test variables set by joining with | between and looking for || showing at least one variable is empty.
# Move the binary to a temporary location and sign it there to avoid the target being up-to-date if signing fails.
if [[ "|$AZURE_SIGNING_CLIENT_ID|$AZURE_SIGNING_CLIENT_SECRET|$AZURE_SIGNING_TENANT_ID|$AZURE_SIGNING_ACCOUNT_ENDPOINT|$AZURE_SIGNING_ACCOUNT_NAME|$AZURE_SIGNING_CERT_PROFILE_NAME|" == *"||"* ]]; then
    >&2 echo "Can't sign windows binaries as required configuration not set: AZURE_SIGNING_CLIENT_ID, AZURE_SIGNING_CLIENT_SECRET, AZURE_SIGNING_TENANT_ID, AZURE_SIGNING_ACCOUNT_ENDPOINT, AZURE_SIGNING_ACCOUNT_NAME, AZURE_SIGNING_CERT_PROFILE_NAME";
    >&2 echo "To rebuild with signing delete the unsigned windows exe file and rebuild with the fixed configuration";
    if [[ "$CI" == "true" ]]; then
        >&2 echo "Signing windows binary is required in CI";
        exit 1;
    fi
    >&2 echo "Skipping signing of windows binary as not in CI";
    exit 0;
fi

file="dist/build-provider-sign-windows_windows_$GORELEASER_ARCH/pulumi-resource-provider-boilerplate.exe";
>&2 echo "Moving $file to $file.unsigned";
mv "$file" "$file.unsigned";

>&2 echo "Logging in to Azure";
az login --service-principal \
    --username "$AZURE_SIGNING_CLIENT_ID" \
    --password "$AZURE_SIGNING_CLIENT_SECRET" \
    --tenant "$AZURE_SIGNING_TENANT_ID";

if [[ $? -ne 0 ]]; then
    >&2 echo "Failed to login to Azure";
    exit 1;
fi

>&2 echo "Getting Trusted Signing access token";
ACCESS_TOKEN="$(az account get-access-token --resource "https://codesigning.azure.net" | jq -r .accessToken)";

if [[ -z "$ACCESS_TOKEN" ]]; then
    >&2 echo "Failed to get access token";
    exit 1;
fi

# jsign expects the keystore as a bare endpoint hostname (no scheme, no trailing slash).
ENDPOINT_HOST="${AZURE_SIGNING_ACCOUNT_ENDPOINT#https://}"
ENDPOINT_HOST="${ENDPOINT_HOST#http://}"
ENDPOINT_HOST="${ENDPOINT_HOST%/}"

>&2 echo "Signing $file.unsigned with Azure Trusted Signing ($ENDPOINT_HOST, $AZURE_SIGNING_ACCOUNT_NAME/$AZURE_SIGNING_CERT_PROFILE_NAME)";
java -jar bin/jsign-7.4.jar \
    --storetype TRUSTEDSIGNING \
    --keystore "$ENDPOINT_HOST" \
    --storepass "$ACCESS_TOKEN" \
    --alias "$AZURE_SIGNING_ACCOUNT_NAME/$AZURE_SIGNING_CERT_PROFILE_NAME" \
    "$file.unsigned";

>&2 echo "Moving $file.unsigned to $file";
mv "$file.unsigned" "$file";

>&2 echo "Logging out from Azure";

if ! az logout; then
    >&2 echo "Failed to logout from Azure, ignoring error";
fi

>&2 echo "Signing of windows binary complete";
