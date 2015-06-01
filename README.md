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
  type: tracker-git
  source:
    token: TRACKER_API_TOKEN
    project_id: "TRACKER_PROJECT_ID"
    tracker_url: https://www.pivotaltracker.com
    repos:
      - git@github.com:you/your_repo
    private_key: GITHUB_PRIVATE_KEY
```

#### Source Configuration

* `token`: *Required.* Your API token, which can be found on your profile page.
* `project_id`: *Required.* Your project ID, which can be found in the URL of your project. Make sure that your `project_id` is a string because it will converted to JSON when given to the resource and JSON doesn't like integers.
* `tracker_url`: *Optional.*
* `repos`: *Required.* The location of the repositories which will contain the branches corresponding to Tracker stories.
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
