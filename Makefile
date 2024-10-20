# Prepare local environment for development.
setup-dev-environment: setup-pre-commit generate-jwt-secret

# Install git pre-commit hooks.
# See https://pre-commit.com for details.
setup-pre-commit:
	@pre-commit install

# Generate tests for Gherkin feature definition.
# See https://github.com/hedhyw/gherkingen for details.
#
# Arguments:
# - path: Relative path to Gherkin feature file.
generate-feature-tests:
	@docker run --rm -it --read-only --network none \
		--volume "$$(pwd)":/host/:ro \
		hedhyw/gherkingen:v4.0.0 \
		-- /host/"$(path)"

generate-jwt-secret:
	@JWT_SECRET="$$(openssl rand -base64 24)" && echo "JWT_SECRET=$$JWT_SECRET" > .env

# Start development infrastructure.
dev-infra-start:
	make dev-infra-stop
	@docker compose up --detach --wait

# Stop development infrastructure.
dev-infra-stop:
	@docker compose down

# Create new user.
# arguments:
#   - name: User name.
#   - pass: User password.
user:
	@go run ./feature/discussion/server/cmd/cli/main.go create user -name="$(name)" -pass="$(pass)"

# Test the application.
test:
	#make dev-infra-start
	go test ./... -race -count=1 -timeout 10s
