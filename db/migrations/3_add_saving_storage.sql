-- +goose Up
CREATE TABLE tax_savings(
    id int not null auto_increment primary key,
    year int not null,
    created_at TIMESTAMP not null,
    comment text,
    amount DECIMAL(10,4),

    INDEX idx_by_year(year)
) DEFAULT CHARACTER SET = utf8;

-- +goose Down

DROP TABLE tax_savings;