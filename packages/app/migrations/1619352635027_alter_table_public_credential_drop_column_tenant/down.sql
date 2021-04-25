ALTER TABLE "public"."credential" ADD COLUMN "tenant" text;
ALTER TABLE "public"."credential" ALTER COLUMN "tenant" DROP NOT NULL;
