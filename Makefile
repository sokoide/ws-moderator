.PHONY: run moderator test gui

run: moderator
	@echo "*****"
	@echo "run 'make gui' first if you haven't done yet"
	@echo "*****"
	./moderator -logLevel DEBUG

moderator:
	@echo building moderator...
	go build ./cmd/moderator

test:
	@echo testing...
	go test -v ./cmd/moderator/...
	go test -v ./pkg/...

gui:
	@echo bulding gui...
	cd gui && npm run build
