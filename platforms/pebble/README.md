# Gobot for pebble

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for Pebble smart watch (http://getpebble.com/).

It uses the Pebble 2.0 SDK, and requires the 2.0 iOS or Android app,
and that the "Chomps" app (https://github.com/hybridgroup/chomps)
has been installed on the Pebble watch.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

[![Build Status](https://secure.travis-ci.org/hybridgroup/gobot-pebble.png?branch=master)](http://travis-ci.org/hybridgroup/gobot-pebble)

## Installing

* Install running: ```go get github.com/hybridgroup/gobot-pebble``
* Install Pebble 2.0 iOS or Android app. (If you haven't already)
* Follow README to install and configure "Chomps" on your watch: https://github.com/hybridgroup/chomps

## Using

* Before running the example, make sure configuration settings match with your program,
in example, api host is your computer IP, robot name is 'pebble', robot api port is 8080 and command is PublishEventC

    your example code here...
## Supported Features

    * We support event detection of 3 main pebble buttons.

## Upcoming Features

* Accelerometer support
* Pushing data to pebble watch

## Documentation

We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing

* All patches must be provided under the Apache 2.0 License
* Please use the -s option in git to "sign off" that the commit is your work and you are providing it under the Apache 2.0 License
* Submit a Github Pull Request to the appropriate branch and ideally discuss the changes with us in IRC.
* We will look at the patch, test it out, and give you feedback.
* Avoid doing minor whitespace changes, renamings, etc. along with merged content. These will be done by the maintainers from time to time but they can complicate merges and should be done seperately.
* Take care to maintain the existing coding style.
* Add unit tests for any new or changed functionality
* All pull requests should be "fast forward"
* If there are commits after yours use “git rebase -i <new_head_branch>”
* If you have local changes you may need to use “git stash”
* For git help see [progit](http://git-scm.com/book) which is an awesome (and free) book on git

## License

Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.
