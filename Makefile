.PHONY: gen tidy test lint all

all: gen tidy test lint

gen:
	sh ./scripts/gen.sh

tidy:
	sh ./scripts/tidy.sh

test:
	sh ./scripts/test.sh

lint:
	sh ./scripts/lint.sh