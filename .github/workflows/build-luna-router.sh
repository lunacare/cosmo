---
name: Build Luna Router
on:
  push:
    branches:
      - main
  workflow_dispatch:
    inputs:
      BUILD_AS_TEST:
        type: boolean
        description: Builds latest-test instead of latest
        required: true
env:
  BUILD_AS_TEST: ${{ inputs.BUILD_AS_TEST }}

jobs:
  build:
    name: Build Luna Router Docker Image
    runs-on:
      group: lunacare-k8s-group
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      - name: Build Docker Image
        run: |-
          export GITHUB_ACTIONS=true
          if [[ github.event_name == "workflow_dispatch" ]]; then
            export BUILD_AS_TEST=$BUILD_AS_TEST
          elif [[ github.event_name == "push" && github.ref == "refs/heads/main" ]]; then
            export BUILD_AS_TEST=false
          fi

          bash ./router/build-and-deploy.sh "$BUILD_AS_TEST"

