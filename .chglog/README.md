# Creating a changelog automatically

## Install and configure tool

We using <https://github.com/git-chglog/git-chglog>, so refer to this side for installation instructions.

It is possible to test the tool by `git-chglog --init` without overriding anything.

## Usage

Example for a new release "v2.5.0":

```sh
# optional update tool by: go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
git checkout release
git pull
git fetch --tags
git checkout dev
git pull upstream dev
git checkout -b rel/prepare_for_release_v250
git-chglog --config .chglog/config_gobot.yml --no-case --next-tag v2.5.0 v2.4.0.. > .chglog/chglog_tmp.md
```

## Compare

If unsure about any result of running git-chglog, just use:
`git log  --since=2024-11-05 --pretty="- %s"`

## Manual adjustment

Most likely some manual work is needed to bring the items in the correct position. We use the style from
["keep a changelog"](https://keepachangelog.com/en/1.1.0/), together with the [standard template](https://github.com/git-chglog/example-type-scope-subject/blob/master/CHANGELOG.standard.md).
The changelog will be generated based on the commit messages, so please follow the
[Convention for Pull Request Descriptions](../CONTRIBUTING.md).

An example for the following commits:

* type(scope): description
* i2c(PCF8583): added
* gpio(HD44780): fix wrong constants
* raspi(PWM): refactor usage
* docs(core): usage of Kernel driver
* or alternative: core(docs): usage of Kernel driver
* build(style): adjust rule for golangci-lint

```md
### Build

* **style**: adjust rule for golangci-lint

### Docs

* **core**: usage of Kernel driver

### I2c

* **PCF8583**: added


### Gpio

* **HD44780**: fix wrong constants

### Raspi

* **PWM**: refactor usage

### Type

* **scope:** description
```

If in doubt, please refer to the current CHANGELOG.md to find the correct way.

## Finalization

After all work is done in the temporary changelog file, the content can be moved to the real one and the "chglog_tmp.md"
file can be removed.
