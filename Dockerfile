# Copyright 2017-2021 The Usacloud Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.17 AS builder
MAINTAINER Usacloud Authors <sacloud.users@gmail.com>

RUN set -x
RUN apt update && apt install -y zip

ADD . /go/src/github.com/sacloud/packer-plugin-sakuracloud

WORKDIR /go/src/github.com/sacloud/packer-plugin-sakuracloud
RUN make tools build
# ======

FROM hashicorp/packer:light
MAINTAINER Usacloud Authors <sacloud.users@gmail.com>

COPY --from=builder /go/src/github.com/sacloud/packer-plugin-sakuracloud/packer-plugin-sakuracloud /bin/