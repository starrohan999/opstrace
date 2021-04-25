CREATE EXTENSION IF NOT EXISTS pgcrypto;
ALTER TABLE "public"."tenant" ADD COLUMN "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid();

ALTER TABLE "public"."credential" ADD COLUMN "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid();
ALTER TABLE "public"."credential" ADD COLUMN "tenant_id" uuid;

ALTER TABLE "public"."exporter" ADD COLUMN "id" uuid NOT NULL UNIQUE DEFAULT gen_random_uuid();
ALTER TABLE "public"."exporter" ADD COLUMN "tenant_id" uuid;
ALTER TABLE "public"."exporter" ADD COLUMN "credential_id" uuid;

UPDATE public.credential c
  SET tenant_id = t.id
  FROM public.tenant t
  WHERE c.tenant = t.name;

UPDATE public.exporter e
  SET tenant_id = t.id
  FROM public.tenant t
  WHERE e.tenant = t.name;

UPDATE public.exporter e
  SET credential_id = c.id
  FROM public.credential c
  WHERE e.tenant_id = c.tenant_id AND e.credential = c.name;

ALTER TABLE "public"."credential" ALTER COLUMN "tenant_id" SET NOT NULL;
ALTER TABLE "public"."exporter" ALTER COLUMN "tenant_id" SET NOT NULL;
ALTER TABLE "public"."exporter" ALTER COLUMN "credential_id" SET NOT NULL;