# Tracker Git Branch Resource

An input-only resource that links [Pivotal Tracker][tracker] to git feature branches.

* On check, the resource will find finished and delivered Tracker stories and return the latest refs from the git branch corresponding to the Tracker story.
* On input, the resource will clone the repository and checkout the appropriate ref.

The git branches are identified by the presence of a story ID in the branch name.

[tracker]: https://www.pivotaltracker.com
