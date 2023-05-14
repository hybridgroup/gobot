# Creating a changelog automatically

## Install and configure tool

We using <https://github.com/git-chglog/git-chglog>, so refer to this side for installation instructions.

## Usage

Example for a new release "v2.0.0":

```sh
git fetch --tags
git-chglog --no-case --next-tag 2.0.0 v1.16.0.. > .chglog/chglog_tmp.md
```

## Compare

If unsure about any result of running git-chglog, just use:
`git log  --since=2022-05-02 --pretty="- %s`

## Manual adjustment

Because there is no commit style guide yet, some manual work is needed to bring the items in the correct position.
Please refer to the current CHANGELOG.md to find the correct way. In general we try to use this style:

* titles will be converted to lower case
* titles are lexical ordered
* each platform has its own title (e.g. **raspi**)
* fixes has the title **bugfix**
* title **api**, **drivers** or **example** is used for changes below related folder
* title **core** is used for changes of common code, e.g. utilities, system
* further special titles **build**, **docs** and **test** can be used

## Finalization

After all work is done in the temporary changelog file, the content can be moved to the real one and the "chglog_tmp.md"
file can be removed.

## TODO's

* introduce a commit style guide
* convert the changelog format to a more common style, see <https://github.com/git-chglog/example-type-scope-subject/tree/master>
