-- +goose Up

INSERT INTO gpm_tax(year, percentage) VALUES(2021, 0.15);

INSERT INTO social_base(year, vdu, count, percentage) VALUES(2021,1352.7 ,43 , 0.9);

INSERT INTO vsd_tax(year, percentage) VALUES(2021, 0.1252);

INSERT INTO psd_tax(year, percentage) VALUES(2021, 0.0698);


-- +goose Down

DELETE FROM gpm_tax WHERE year=2021;

DELETE FROM social_base WHERE year=2021;

DELETE FROM vsd_tax WHERE year=2021;

DELETE FROM psd_tax WHERE year=2021;
