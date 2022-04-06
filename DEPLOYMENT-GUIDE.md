# Deployment Guide

In order to run the indexer the first step is to build the container images.

For this the command bellow allows building and tagging the containers.
Replace `<name>` with:

* api
* jobwatcher
* parsingdispatcher
* watcher

`docker build . -f Dockerfile.<name> -t indexer:<name>`


(Related issue: https://github.com/NFT-com/indexer/issues/47)
** TODO **
