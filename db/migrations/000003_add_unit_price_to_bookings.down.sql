-- Add unit_price column to bookings table
ALTER TABLE bookings ADD COLUMN IF NOT EXISTS unit_price BIGINT;

-- Backfill unit_price for existing rows
UPDATE bookings SET unit_price = total_amount / NULLIF(quantity, 0) WHERE unit_price IS NULL;

-- Enforce NOT NULL constraint for future inserts
ALTER TABLE bookings ALTER COLUMN unit_price SET NOT NULL;

-- Document column
COMMENT ON COLUMN bookings.unit_price IS 'Snapshot of ticket price at booking time';
