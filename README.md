# saasus-sdk-go

This is a SaaSus SDK for Go

## How to use

See the documentation [Implementation sample](https://github.com/saasus-platform/implementation-sample-api-go/)

## SDK version

The SDK uses the version committed in `VERSION` to identify itself in the
`User-Agent` header, including when built from a local source checkout.

The API release workflow sets the `VERSION` environment variable when running
`generate.sh`. For a manual release, run `sh update-version.sh vX.Y.Z` before
creating the tag and commit both `VERSION` and the generated Go version file.
CI verifies that the files agree and that the release tag matches them.
