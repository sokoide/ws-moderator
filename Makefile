.PHONY: run moderator test gui fonts

run: moderator
	@echo "*****"
	@echo "run 'make gui' first if you haven't done yet"
	@echo "*****"
	./moderator -logLevel DEBUG

moderator: fonts
	@echo building moderator...
	go build ./cmd/moderator

test:
	@echo testing...
	go test -v ./cmd/moderator/...
	go test -v ./pkg/...

gui:
	@echo bulding gui...
	cd gui && npm run build

fonts: NotoSansJP-Regular.ttf NotoSansJP-Bold.ttf

NotoSansJP-Regular.ttf:
	@echo 'Please donload NotoSansJP-Bold.ttf and NotoSansJP-Bold.ttf from https://fonts.google.com/noto/specimen/Noto+Sans+JP' into the repo root folder
	exit 1

NotoSansJP-Bold.ttf:
	@echo 'Please donload NotoSansJP-Bold.ttf and NotoSansJP-Bold.ttf from https://fonts.google.com/noto/specimen/Noto+Sans+JP' into the repo root folder
	exit 1
