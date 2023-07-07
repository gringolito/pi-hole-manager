.PHONY: clean
clean:
	@rm -f .phm-debuilder debian/build/*

.PHONY: deb
deb: .phm-debuilder
	@docker run --rm -it -v $(CURDIR):/src phm-debuilder

.phm-debuilder:
	@docker build -t phm-debuilder debian
	@touch $@
