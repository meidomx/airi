CREATE TABLE public.airi_task (
                                  task_id bigserial NOT NULL,
                                  task_key varchar(400) NOT NULL,
                                  description varchar(4000) NOT NULL,
                                  "type" smallint NOT NULL,
                                  status smallint NOT NULL,
                                  config text NOT NULL,
                                  time_created bigint NOT NULL,
                                  time_updated bigint NOT NULL,
                                  CONSTRAINT airi_task_pk PRIMARY KEY (task_id)
);
CREATE INDEX airi_task_status_idx ON public.airi_task (status,time_created);
CREATE UNIQUE INDEX airi_task_task_key_idx ON public.airi_task (task_key);
