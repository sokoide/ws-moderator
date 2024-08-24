.PHONY: build run moderator gui

build: gui moderator

run: build
	./moderator

moderator: gui
	@echo building moderator...
	go build ./cmd/moderator

gui:
	@echo bulding gui...
	cd gui; npm run build
