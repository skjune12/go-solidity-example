SOLC = solc
ABIGEN = abigen
CONTRACT = Store
PKGNAME = main
BINARY = ethtest

.PHONY: all
all: clean solc abigen build

.PHONY: solc
solc: $(CONTRACT).sol
	$(SOLC) $(CONTRACT).sol --abi --bin -o .

.PHONY: abigen
abigen: $(CONTRACT).sol
	abigen --sol $(CONTRACT).sol --pkg $(PKGNAME) --out=$(CONTRACT).go

.PHONY: build
build:
	go build -o ${BINARY}

.PHONY: clean
clean:
	rm -rf $(CONTRACT).abi $(CONTRACT).bin $(BINARY)
