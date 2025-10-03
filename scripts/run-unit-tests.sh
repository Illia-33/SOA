#!/bin/bash
set -e

(cd $PWD/services/accounts; make test)
(cd $PWD/services/posts; make test)
(cd $PWD/services/stats; make test)
