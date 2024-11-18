CREATE TABLE voucher_item (
    id BIGSERIAL PRIMARY KEY,
    voucher_id BIGINT NOT NULL REFERENCES voucher(id) ON DELETE CASCADE,
    sl_id BIGINT NOT NULL REFERENCES sl(id),
    dl_id BIGINT REFERENCES dl(id),
    debit INT CHECK (debit >= 0),
    credit INT CHECK (credit >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT sl_must_have_dl CHECK (
        (SELECT has_dl FROM sl WHERE id = sl_id) = false OR dl_id IS NOT NULL
    ),
    CONSTRAINT debit_or_credit_must_be_positive CHECK (
        (debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0)
    )
);
