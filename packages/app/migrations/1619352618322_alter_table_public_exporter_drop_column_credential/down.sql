ALTER TABLE "public"."exporter" ADD COLUMN "credential" text;
ALTER TABLE "public"."exporter" ALTER COLUMN "credential" DROP NOT NULL;
