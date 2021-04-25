ALTER TABLE "public"."exporter" DROP COLUMN "credential_id";
ALTER TABLE "public"."exporter" DROP COLUMN "tenant_id";
ALTER TABLE "public"."exporter" DROP COLUMN "id";

ALTER TABLE "public"."credential" DROP COLUMN "tenant_id";
ALTER TABLE "public"."credential" DROP COLUMN "id";

ALTER TABLE "public"."tenant" DROP COLUMN "id";
