ALTER TABLE "public"."exporter" ADD COLUMN "tenant" text;
ALTER TABLE "public"."exporter" ALTER COLUMN "tenant" DROP NOT NULL;
