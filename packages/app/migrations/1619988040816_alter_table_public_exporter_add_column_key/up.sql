ALTER TABLE "public"."exporter" ADD COLUMN "key" text NULL;
ALTER TABLE "public"."exporter" ADD CONSTRAINT "exporter_key_key" UNIQUE ("key");

-- There are already DNS records using the exporter name, so need to keep those working
UPDATE "public"."exporter" e SET key = e.name;

ALTER TABLE "public"."exporter" ALTER COLUMN "key" SET NOT NULL;

