-- +goose Up

CREATE TABLE incomes(
    id int not null auto_increment primary key,
    amount DECIMAL(10,4) not null,
    earned TIMESTAMP not null,
    year INT NOT NULL,
    origin VARCHAR(50),
    comment TEXT,

    INDEX idx_by_year(year)
) DEFAULT CHARACTER SET = utf8;

-- +goose Down
DROP TABLE incomes;