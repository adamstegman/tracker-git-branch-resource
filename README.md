# Tracker Git Branch Resource

An input-only resource that links [Pivotal Tracker][tracker] stories to git feature branches.

* On check, the resource will find finished and delivered Tracker stories and return the latest refs from the git branch corresponding to the Tracker story.
* On input, the resource will clone the repository and checkout the appropriate ref.

The git branches are identified by the presence of a story ID in the branch name.

[tracker]: https://www.pivotaltracker.com

## Configuration

### Resource

``` yaml
- name: tracker
  type: tracker-git-branch
  source:
    token: TRACKER_API_TOKEN
    projects:
      - "TRACKER_PROJECT_ID"
    tracker_url: https://www.pivotaltracker.com
    repo: git@github.com:you/your_repo
    private_key: GITHUB_PRIVATE_KEY
```

#### Source Configuration

* `token`: *Required.* Your API token, which can be found on your profile page.
* `projects`: *Required.* Your Tracker project IDs, which can be found in the URL of your project.
  Make sure that each value is a string because it will converted to JSON when given to the resource and JSON doesn't like integers.
* `tracker_url`: *Optional.*
* `repo`: *Required.* The location of the repository which will contain the branches corresponding to Tracker stories.
* `private_key`: *Optional.* Private key to use when pulling/pushing.
    Example:
    ```
    private_key: |
      -----BEGIN RSA PRIVATE KEY-----
      MIIEowIBAAKCAQEAtCS10/f7W7lkQaSgD/mVeaSOvSF9ql4hf/zfMwfVGgHWjj+W
      <Lots more text>
      DWiJL+OFeg9kawcUL6hQ8JeXPhlImG6RTUffma9+iGQyyBMCGd1l
      -----END RSA PRIVATE KEY-----
    ```

You'll need a seperate resource for each Tracker project.

## Development

Run `scripts/test` to execute the tests using [Ginkgo][].
The branches in this project are fixtures for the tests, so your working directory must not contain any uncommitted changes before running this script.

[ginkgo]: http://onsi.github.io/ginkgo/
