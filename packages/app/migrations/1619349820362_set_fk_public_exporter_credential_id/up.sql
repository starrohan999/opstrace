alter table "public"."exporter" drop constraint "exporter_tenant_credential_fkey",
             add constraint "exporter_credential_id_fkey"
             foreign key ("credential_id")
             references "public"."credential"
             ("id") on update cascade on delete restrict;
