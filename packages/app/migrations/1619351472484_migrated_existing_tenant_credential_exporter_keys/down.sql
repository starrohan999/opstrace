ALTER TABLE "public"."exporter" ALTER COLUMN "tenant_id" DROP NOT NULL;
ALTER TABLE "public"."credential" ALTER COLUMN "tenant_id" DROP NOT NULL;