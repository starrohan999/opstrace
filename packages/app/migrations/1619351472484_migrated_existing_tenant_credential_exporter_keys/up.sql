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