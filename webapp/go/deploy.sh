#!/bin/bash

set -x

./deploy_body.sh | notify_slack
