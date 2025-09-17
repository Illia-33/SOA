#!/bin/bash
set -e

(cd $PWD/services/accounts; make proto)
(cd $PWD/services/posts; make proto)
(cd $PWD/services/stats; make proto)
