CC      	= go
PROGRAM		= outbox-relay
BUILDDIR	= build

.PHONY: build clean cleanBuild

build: $(PROGRAM)

$(PROGRAM):
	$(CC) build

clean:
	rm -f $(PROGRAM)

cleanBuild:
	rm -rf $(BUILDDIR)/*


