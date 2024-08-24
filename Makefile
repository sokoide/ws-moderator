.PHONY: run moderator gui

run: moderator
	@echo "*****"
	@echo "run 'make gui' first if you haven't done yet"
	@echo "*****"
	./moderator

moderator:
	@echo building moderator...
	go build ./cmd/moderator

gui:
	@echo bulding gui...
	cd gui && npm run build
