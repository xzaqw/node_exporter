BIN_DIR ?= /usr/local/bin
CFG_DIR ?= /etc
SVC_DIR ?= /lib/svc/method

all:
	@mkdir -p $(CFG_DIR)
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(SVC_DIR)
	@cp ./kstat_config.yml $(CFG_DIR)
	@cp ../../node_exporter $(BIN_DIR)/prometheus-node-exporter
	@cp ./node_exporter $(SVC_DIR)
	@cp ./node_exporter.xml $(SVC_DIR)
	sudo svccfg -v import $(SVC_DIR)/node_exporter.xml
