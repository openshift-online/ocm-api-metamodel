#
# Copyright (c) 2019 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

.PHONY: \
	clean \
	cmds \
	fmt \
	lint \
	$(NULL)

cmds: generate
	for cmd in $$(ls cmd); do \
		go build -mod=readonly -o "$${cmd}" "./cmd/$${cmd}" || exit 1; \
		cp "$${cmd}" "$${HOME}/bin/."; \
	done

generate: antlr
	cd pkg/language && \
	java \
		-jar ../../$? \
		-package language \
		*.g4

antlr:
	wget \
		--output-document=$@ \
		https://www.antlr.org/download/antlr-4.7.2-complete.jar

fmt:
	gofmt -s -l -w \
		./cmd \
		./pkg \
		./tests \
		$(NULL)

test: cmds
	./ocm-metamodel-tool generate \
		--model=/files/go/src/gitlab.cee.redhat.com/service/ocm-api-model/model \
		--base=gitlab.cee.redhat.com/service/ocm-api-metamodel/tests/api \
		--output=tests/api \
		--docs=tests/docs
	ginkgo tests

clean:
	rm -rf \
		$$(ls cmd) \
		antlr \
		tests/api \
		$(NULL)
