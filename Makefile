setup-dev-environment: pre-commit

# Install git pre-commit hooks.
# See https://pre-commit.com for details.
pre-commit:
	pre-commit install

# Generate tests for Gherkin feature definition.
# See https://github.com/hedhyw/gherkingen for details.
#
# Arguments:
# - path: Relative path to Gherkin feature file.
generate-feature-tests:
	docker run --rm -it --read-only --network none \
		--volume "$$(pwd)":/host/:ro \
		hedhyw/gherkingen:v4.0.0 \
		-- /host/"$(path)"
