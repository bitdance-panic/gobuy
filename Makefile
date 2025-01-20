.PHONY: gen tidy test lint all

all: gen tidy test lint

gen:
	./scripts/gen.sh

tidy:
	./scripts/tidy.sh

test:
	./scripts/test.sh

lint:
	./scripts/lint.sh