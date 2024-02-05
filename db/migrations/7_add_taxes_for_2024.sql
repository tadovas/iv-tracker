-- +goose Up

INSERT INTO gpm_tax(year, percentage) VALUES(2024, 0.15);

INSERT INTO social_base(year, vdu, count, percentage) VALUES(2024,1902.70 ,43 , 0.9);

INSERT INTO vsd_tax(year, percentage) VALUES(2024, 0.1252);

INSERT INTO psd_tax(year, percentage) VALUES(2024, 0.0698);


-- +goose Down

DELETE FROM gpm_tax WHERE year=2024;

DELETE FROM social_base WHERE year=2024;

DELETE FROM vsd_tax WHERE year=2024;

DELETE FROM psd_tax WHERE year=2024;
