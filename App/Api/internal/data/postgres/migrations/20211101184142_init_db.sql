-- +goose Up
-- +goose StatementBegin

CREATE TABLE "types" (
  "id" SERIAL PRIMARY KEY,
  "name" text NOT NULL
);

CREATE TABLE "pictures" (
  "id" SERIAL PRIMARY KEY,
  "guid" text NOT NULL,
  "content_type" text NOT NULL
);

CREATE TABLE "pets" (
  "id" SERIAL PRIMARY KEY,
  "picture_id" int,
  "type_id" int NOT NULL,
  "name" text,
  "color" text,
  "marks" text,
  CONSTRAINT type_fk FOREIGN KEY ("type_id")
        REFERENCES public.types ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
  CONSTRAINT picture_fk FOREIGN KEY ("picture_id")
        REFERENCES public.pictures ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

CREATE TABLE "pet_breeds" (
  "id" SERIAL PRIMARY KEY,
  "pet_id" int NOT NULL,
  "name" text NOT NULL,
  CONSTRAINT breed_pet_fk FOREIGN KEY ("pet_id")
        REFERENCES public.pets ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE "tags" (
  "id" SERIAL PRIMARY KEY,
  "pet_id" int NOT NULL,
  "shape" text,
  "color" text,
  "text" text,
  CONSTRAINT tag_pet_fk FOREIGN KEY ("pet_id")
        REFERENCES public.pets ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);


CREATE TABLE "postings" (
  "id" SERIAL PRIMARY KEY,
  "guid" text NOT NULL,
  "pet_id" int NOT NULL,
  "date" timestamp with time zone NOT NULL,
  "location" text NOT NULL,
  "name" text,
  "email" text NOT NULL,
   CONSTRAINT postings_pet_fk FOREIGN KEY ("pet_id")
        REFERENCES public.pets ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);

CREATE TABLE "sightings" (
  "id" SERIAL PRIMARY KEY,
  "guid" text NOT NULL,
  "pet_id" int NOT NULL,
  "in_custody" boolean NOT NULL,
  "date" timestamp with time zone NOT NULL,
  "location" text NOT NULL,
  "name" text,
  "email" text NOT NULL,
  CONSTRAINT sightings_pet_fk FOREIGN KEY ("pet_id")
        REFERENCES public.pets ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
);


CREATE TABLE "matches" (
  "postings_id" int NOT NULL,
  "sightings_id" int NOT NULL,
  "last_contacted" timestamp with time zone,
  PRIMARY KEY ("postings_id", "sightings_id"),
  CONSTRAINT matches_sighting_fk FOREIGN KEY ("sightings_id")
        REFERENCES public.sightings ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE,
  CONSTRAINT matches_posting_fk FOREIGN KEY ("postings_id")
        REFERENCES public.postings ("id") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE CASCADE
        
);




-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table public.matches;
drop table public.sightings;
drop table public.postings;
drop table public.tags;
drop table public.pet_breeds;
drop table public.pets;
drop table public.types;
drop table public.pictures;
-- +goose StatementEnd
