alter table "public"."exporter" drop constraint "exporter_tenant_id_fkey",
          add constraint "exporter_tenant_fkey"
          foreign key ("tenant")
          references "public"."tenant"
          ("name")
          on update cascade
          on delete cascade;
