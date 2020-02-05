-- +goose Up

CREATE TABLE gpm_tax(
    year INT NOT NULL PRIMARY KEY,
    percentage DECIMAL(10,4) not null
);

INSERT INTO gpm_tax(year , percentage) VALUES(2019, 0.15);
INSERT INTO gpm_tax(year, percentage) VALUES(2020, 0.15);

CREATE TABLE social_base(
  year INT NOT NULL PRIMARY KEY,
  vdu DECIMAL(10,4) not null,
  count INT NOT NULL,
  percentage DECIMAL(10,4) not null
);

INSERT INTO social_base(year, vdu, count, percentage) VALUES(2019,1136.2 ,43 , 0.9);
INSERT INTO social_base(year, vdu, count, percentage) VALUES(2020,1241.4 ,43 , 0.9);

CREATE TABLE vsd_tax(
    year INT NOT NULL PRIMARY KEY,
    percentage DECIMAL(10,4) not null
);

INSERT INTO vsd_tax(year, percentage) VALUES(2019, 0.1252);
INSERT INTO vsd_tax(year, percentage) VALUES(2020, 0.1252);

CREATE TABLE psd_tax(
    year INT NOT NULL PRIMARY KEY,
    percentage DECIMAL(10,4) not null
);

INSERT INTO psd_tax(year, percentage) VALUES(2019, 0.0698);
INSERT INTO psd_tax(year, percentage) VALUES(2020, 0.0698);


-- +goose Down
DROP TABLE gpm_tax;

DROP TABLE social_base;

DROP TABLE vsd_tax;

DROP TABLE psd_tax;