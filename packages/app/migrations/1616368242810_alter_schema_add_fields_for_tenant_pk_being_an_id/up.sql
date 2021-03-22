CREATE EXTENSION IF NOT EXISTS pgcrypto;
ALTER TABLE "public"."tenant" ADD COLUMN "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid();
ALTER TABLE "public"."exporter" ADD COLUMN "tenant_id" uuid NOT NULL;
ALTER TABLE "public"."credential" ADD COLUMN "tenant_id" uuid NOT NULL;
