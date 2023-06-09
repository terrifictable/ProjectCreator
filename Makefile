GO := go

branch := $(shell git rev-parse --abbrev-ref HEAD)
commit := $(shell git rev-parse --short HEAD)

LDFLAGS := -ldflags "-X main.BRANCH=${branch} -X main.COMMIT=${commit} -X projs/common.dbg=${DBG}"
FLAGS := 


build:
	@$(MAKE) exec ACTION=build	

run:
	@$(MAKE) exec ACTION=run

exec:
	@mkdir -p .logs/
	$(GO) $(ACTION) $(LDFLAGS) $(FLAGS) . $(ARGS) 
	@#| tee .logs/$(ACTION)-$(shell date +"%Y-%m-%d").log
