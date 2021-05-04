ALTER TABLE "public"."exporter" ADD COLUMN "key" text NULL;
ALTER TABLE "public"."exporter" ADD CONSTRAINT "exporter_unique_key_for_tenant" UNIQUE ("key", "tenant_id");

-- There are already DNS records using the exporter name, so need to keep those working
UPDATE "public"."exporter" e SET key = e.name;

ALTER TABLE "public"."exporter" ALTER COLUMN "key" SET NOT NULL;

