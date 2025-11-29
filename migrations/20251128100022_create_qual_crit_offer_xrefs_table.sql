CREATE TABLE IF NOT EXISTS blc_qual_crit_offer_xref (
    offer_qual_crit_id BIGSERIAL PRIMARY KEY,
    offer_id BIGINT NOT NULL,
    offer_item_criteria_id BIGINT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_blc_qual_crit_offer_xref_offer_id FOREIGN KEY (offer_id) REFERENCES blc_offer(offer_id),
    CONSTRAINT fk_blc_qual_crit_offer_xref_offer_item_criteria_id FOREIGN KEY (offer_item_criteria_id) REFERENCES blc_offer_item_criteria(offer_item_criteria_id)
);
