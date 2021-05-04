CFLAGS="-Qunused-arguments"
SHELL=bash
GOOS=linux
GOARCH=amd64

gofiles = $(shell find pkg cmd go.* -print)

.PHONY: python
python: python/libtrivy/libtrivy.so
	PLAT=$$(python -c 'from distutils.util import get_platform; print(get_platform())'); \
	python setup.py bdist_wheel -p $$PLAT

.PHONY: golib
golib: build/libtrivy.so

.PHONY: clean
clean:
	@rm -rf build dist 2> /dev/null || true
	@rm python/libtrivy/*.{so,h,dynlib} 2> /dev/null || true
	@rm -rf python/libtrivy.egg-info 2> /dev/null || true

build/libtrivy.so: $(gofiles)
	./scripts/build.sh $@

python/libtrivy/%: build/%
	cp $< $@

$(gofiles):
	touch $@
