ALTER TABLE "public"."credential" ADD COLUMN "key" text NULL;
ALTER TABLE "public"."credential" ADD CONSTRAINT "credential_unique_key_for_tenant" UNIQUE ("key", "tenant_id");

-- There are already DNS records using the credential name, so need to keep those working
UPDATE "public"."credential" c SET key = c.name;

ALTER TABLE "public"."credential" ALTER COLUMN "key" SET NOT NULL;

