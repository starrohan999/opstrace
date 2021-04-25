-- ALTER TABLE "public"."tenant" ADD COLUMN "url_slug" text NULL;
-- ALTER TABLE "public"."tenant" ADD CONSTRAINT "tenant_url_slug_key" UNIQUE ("url_slug");

-- -- There are already DNS records using the tenant name, so need to keep those working
-- UPDATE "public"."tenant" t SET url_slug = t.name;
