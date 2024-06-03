# github-dump-issue-comments

Sometimes you just want to quickly search all the comments in a really long
issue and the github UI makes that quite annoying. This tool is for those times.


## Usage

```
Usage: github-dump-issue-comments [<url>|<owner> <repo> <issue_number>]
```

Example:

```
$ github-dump-issue-comments  https://github.com/AsahiLinux/linux/issues/72 | head
[2022-12-07 03:36:00 +0000 UTC] asahilina:
This is a tracker bug for general GPU issues, like:

* Apps that crash after startup
* Rendering glitches
* GPU fault/timeout errors

When making a comment on this bug, please run the `asahi-diagnose` command and attach the file it saves to your comment. Please tell us what you were doing when the problem happened, what desktop environment and window system you use, and any other details about the issue.

The purpose of this bug is to collect reports of app issues in one place, so we have somewhere to look when figuring out what to work on. Since the driver is still a work-in-progress and lots of things are not expected to work, please don't expect a timely response to reports. We're working on it!
```
