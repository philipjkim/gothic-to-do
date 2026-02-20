.PHONY: dev build css templ

# 개발 서버 (air + tailwind watch)
dev:
	@echo "Starting development server..."
	@make -j2 air css-watch

air:
	air

css-watch:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --watch

# 빌드
build: templ css
	go build -o ./bin/server ./cmd/server

templ:
	templ generate

css:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css --minify
