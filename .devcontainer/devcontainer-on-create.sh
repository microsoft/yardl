#! /usr/bin/env bash

if [ -s /opt/matlab/latest/licenses/license.host.lic ]; then
    cp /opt/matlab/latest/licenses/license.host.lic /opt/matlab/latest/licenses/license.lic
fi
