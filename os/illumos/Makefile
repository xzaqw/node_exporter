BIN_DIR ?= /usr/local/bin
CFG_DIR ?= /etc
SVC_DIR ?= /lib/svc/method

BIN_PATH := $(BIN_DIR)/prometheus-node-exporter
MANIFEST_PATH := $(SVC_DIR)/prometheus-node-exporter.xml
SMF_METHOD_PATH := $(SVC_DIR)/method/prometheus-node-exporter
KSTAT_CONFIG_PATH := $(CFG_DIR)/kstat_config.yml 

all:
	@mkdir -p $(CFG_DIR)
	@mkdir -p $(BIN_DIR)
	@mkdir -p $(SVC_DIR)/method
	@cp ./kstat_config.yml $(KSTAT_CONFIG_PATH)
	@cp ../../node_exporter $(BIN_PATH)
	@cp ./prometheus-node-exporter $(SMF_METHOD_PATH)
	@cp ./prometheus-node-exporter.xml $(MANIFEST_PATH)
	sudo svccfg -v import $(MANIFEST_PATH)
