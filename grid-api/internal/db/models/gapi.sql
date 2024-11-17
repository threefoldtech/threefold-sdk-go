CREATE TYPE "contract_state" AS ENUM (
  'created',
  'deleted',
  'grace_period',
  'out_of_funds'
);

CREATE TYPE "discount" AS ENUM (
  'none',
  'default',
  'bronze',
  'silver',
  'gold'
);

CREATE TYPE "NodeState" AS ENUM (
  'Up',
  'Down',
  'Standby'
);

CREATE TYPE "grid_version" AS ENUM (
  'One'
);

CREATE TYPE "certification" AS ENUM (
  'Gold',
  'NotCertified',
  'DIY',
  'Certified'
);

CREATE TABLE "twin" (
  "id" integer PRIMARY KEY,
  "grid_version" grid_version,
  "account_id" varchar(100),
  "public_key" varchar(100),
  "relay" varchar(100)
);

CREATE TABLE "farm" (
  "id" integer PRIMARY KEY,
  "name" varchar(100),
  "twin_id" integer,
  "stellar_address" varchar(100),
  "dedicated" boolean,
  "pricing_policy_id" integer,
  "certification" certification,
  "grid_version" grid_version
);

CREATE TABLE "node" (
  "id" integer PRIMARY KEY,
  "twin_id" integer,
  "farm_id" integer,
  "total_hru" integer,
  "total_cru" integer,
  "total_sru" integer,
  "total_mru" integer,
  "serial_number" integer,
  "created_at" timestamp,
  "extra_fee" integer,
  "grid_version" grid_version,
  "certification" certification
);

CREATE TABLE "node_contract" (
  "id" integer PRIMARY KEY,
  "node_id" integer,
  "twin_id" integer,
  "contract_resource_id" integer,
  "solution_provider_id" integer,
  "state" contract_state,
  "created_at" timestamp,
  "deployment_data" text,
  "deployment_hash" text,
  "ips_num" integer
);

CREATE TABLE "rent_contract" (
  "id" integer PRIMARY KEY,
  "node_id" integer,
  "twin_id" integer,
  "contract_resource_id" integer,
  "solution_provider_id" integer,
  "state" contract_state,
  "created_at" timestamp
);

CREATE TABLE "name_contract" (
  "id" integer PRIMARY KEY,
  "twin_id" integer,
  "contract_resource_id" integer,
  "solution_provider_id" integer,
  "state" contract_state,
  "created_at" timestamptz
);

CREATE TABLE "ContractBill" (
  "contract_id" integer,
  "discount" discount,
  "timestamp" timestamp,
  "amount" integer
);

CREATE TABLE "node_info" (
  "node_id" integer,
  "has_ipv6" boolean,
  "num_workloads" integer,
  "upload_speed" numeric,
  "download_speed" numeric
);

CREATE TABLE "dmi" (
  "node_id" integer,
  "bios" jsonb,
  "baseboard" jsonb,
  "processor" jsonb,
  "memory" jsonb
);

CREATE TABLE "gpu" (
  "node_id" integer,
  "id" text,
  "vendor" text,
  "device" text,
  "contract" integer
);

CREATE TABLE "location" (
  "node_id" integer,
  "langitude" integer,
  "latitude" integer,
  "country" varchar(100),
  "city" varchar(100),
  "region" varchar(20)
);

CREATE TABLE "PublicConfig" (
  "node_id" integer,
  "ipv4" varchar(50),
  "ipv6" varchar(50),
  "gw4" varchar(50),
  "gw6" varchar(50),
  "domain" varchar(50)
);

CREATE TABLE "Interfaces" (
  "node_id" integer,
  "name" varchar(50),
  "mac" varchar(50),
  "ips" varchar(50)
);

CREATE TABLE "power" (
  "node_id" integer,
  "state" NodeState,
  "target" NodeState,
  "status" NodeState,
  "Healthy" boolean,
  "last_uptime_report" timestamp,
  "total_uptime" integer
);

CREATE TABLE "contract_resource" (
  "id" integer PRIMARY KEY,
  "hru" integer,
  "sru" integer,
  "mru" integer,
  "cru" integer
);

CREATE TABLE "solution_provider" (
  "id" integer PRIMARY KEY,
  "description" text,
  "link" text,
  "approved" boolean,
  "providers" text
);

CREATE TABLE "pricing_policy" (
  "id" integer PRIMARY KEY,
  "grid_version" grid_version,
  "name" varchar(100),
  "foundation_account" varchar(100),
  "certified_sales_account" varchar(100),
  "dedication_discount" integer,
  "su_value" integer,
  "su_unit" varchar(10),
  "cu_value" integer,
  "cu_unit" varchar(10),
  "nu_value" integer,
  "nu_unit" varchar(10),
  "ipu_value" integer,
  "ipu_unit" varchar(10)
);

CREATE TABLE "public_ip" (
  "farm_id" integer,
  "contract_id" integer,
  "gateway" varchar(25),
  "ip" varchar(25)
);

CREATE INDEX "idx_twin_relay" ON "twin" ("relay");

CREATE INDEX "idx_farm_twin" ON "farm" ("twin_id");

CREATE INDEX "idx_node_farm" ON "node" ("farm_id");

CREATE INDEX "idx_node_contract_node" ON "node_contract" ("node_id");

CREATE INDEX "idx_node_contract_twin" ON "node_contract" ("twin_id");

CREATE INDEX "idx_node_contract_state" ON "node_contract" ("state");

CREATE INDEX "idx_rent_contract_node" ON "rent_contract" ("node_id");

CREATE INDEX "idx_rent_contract_twin" ON "rent_contract" ("twin_id");

CREATE INDEX "idx_rent_contract_state" ON "rent_contract" ("state");

CREATE INDEX "idx_name_contract_twin" ON "name_contract" ("twin_id");

CREATE INDEX "idx_name_contract_state" ON "name_contract" ("state");

ALTER TABLE "farm" ADD FOREIGN KEY ("twin_id") REFERENCES "twin" ("id");

ALTER TABLE "farm" ADD FOREIGN KEY ("pricing_policy_id") REFERENCES "pricing_policy" ("id");

ALTER TABLE "node" ADD FOREIGN KEY ("twin_id") REFERENCES "twin" ("id");

ALTER TABLE "node" ADD FOREIGN KEY ("farm_id") REFERENCES "farm" ("id");

ALTER TABLE "node_contract" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "node_contract" ADD FOREIGN KEY ("twin_id") REFERENCES "twin" ("id");

ALTER TABLE "node_contract" ADD FOREIGN KEY ("contract_resource_id") REFERENCES "contract_resource" ("id");

ALTER TABLE "node_contract" ADD FOREIGN KEY ("solution_provider_id") REFERENCES "solution_provider" ("id");

ALTER TABLE "rent_contract" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "rent_contract" ADD FOREIGN KEY ("twin_id") REFERENCES "twin" ("id");

ALTER TABLE "rent_contract" ADD FOREIGN KEY ("contract_resource_id") REFERENCES "contract_resource" ("id");

ALTER TABLE "rent_contract" ADD FOREIGN KEY ("solution_provider_id") REFERENCES "solution_provider" ("id");

ALTER TABLE "name_contract" ADD FOREIGN KEY ("twin_id") REFERENCES "twin" ("id");

ALTER TABLE "name_contract" ADD FOREIGN KEY ("contract_resource_id") REFERENCES "contract_resource" ("id");

ALTER TABLE "name_contract" ADD FOREIGN KEY ("solution_provider_id") REFERENCES "solution_provider" ("id");

ALTER TABLE "ContractBill" ADD FOREIGN KEY ("contract_id") REFERENCES "node_contract" ("id");

ALTER TABLE "node_info" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "dmi" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "gpu" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "location" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "PublicConfig" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "Interfaces" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "power" ADD FOREIGN KEY ("node_id") REFERENCES "node" ("id");

ALTER TABLE "public_ip" ADD FOREIGN KEY ("farm_id") REFERENCES "farm" ("id");

CREATE TABLE "node_contract_public_ip" (
  "node_contract_id" integer,
  "public_ip_contract_id" integer,
  PRIMARY KEY ("node_contract_id", "public_ip_contract_id")
);

ALTER TABLE "node_contract_public_ip" ADD FOREIGN KEY ("node_contract_id") REFERENCES "node_contract" ("id");

ALTER TABLE "node_contract_public_ip" ADD FOREIGN KEY ("public_ip_contract_id") REFERENCES "public_ip" ("contract_id");

