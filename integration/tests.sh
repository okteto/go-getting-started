#!/usr/bin/env bash

(
	# Runs integration tests against a pre deployed instance of this application
	# We use okteto tests to run this tests so we go through the internal network
	# and not the ingress

	set -e # make any error fail the script
	set -u # make unbound variables fail the script

	# SC2039: In POSIX sh, set option pipefail is undefined
	# shellcheck disable=SC2039
	set -o pipefail # make any pipe error fail the script

	expected="Hello world!"
	actual=$(curl -s "http://hello-world.${OKTETO_NAMESPACE}:8080")

	if [ "$actual" != "$expected" ]; then
		echo "Expected 'Hello world!' but got '$actual'"
		exit 1
	fi

	echo "OK"
)
