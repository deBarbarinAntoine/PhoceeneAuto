CREATE TABLE IF NOT EXISTS cars_catalog (
                                            id               bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
                                            created_at      TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                            updated_at      TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                            make               TEXT NOT NULL,
                                            model              TEXT NOT NULL,
                                            cylinders          INT,
                                            drive             TEXT,
                                            engine_descriptor TEXT,
                                            fuel1             TEXT,
                                            fuel2             TEXT,
                                            luggage_volume    FLOAT,
                                            passenger_volume  FLOAT,
                                            transmission      TEXT,
                                            size_class        TEXT,
                                            model_year              INT,
                                            electric_motor    FLOAT,
                                            base_model        TEXT,
                                            version       INT DEFAULT 1,
                                            CONSTRAINT unique_car UNIQUE (
                                                                          make, model, model_year, cylinders, drive, engine_descriptor, fuel1, fuel2,
                                                                          luggage_volume, passenger_volume, transmission, size_class, electric_motor, base_model, version
                                                )
);

CREATE OR REPLACE FUNCTION check_duplicate_car()
    RETURNS TRIGGER AS $$
BEGIN
    -- Check if the new row being inserted already exists in the table
    IF EXISTS (SELECT 1 FROM cars_catalog WHERE
        cars_catalog.make = NEW.make AND
        cars_catalog.model = NEW.model AND
        cars_catalog.model_year = NEW.model_year AND
        cars_catalog.cylinders = NEW.cylinders AND
        cars_catalog.drive = NEW.drive AND
        cars_catalog.engine_descriptor = NEW.engine_descriptor AND
        cars_catalog.fuel1 = NEW.fuel1 AND
        cars_catalog.fuel2 = NEW.fuel2 AND
        cars_catalog.luggage_volume = NEW.luggage_volume AND
        cars_catalog.passenger_volume = NEW.passenger_volume AND
        cars_catalog.transmission = NEW.transmission AND
        cars_catalog.size_class = NEW.size_class AND
        cars_catalog.electric_motor = NEW.electric_motor AND
        cars_catalog.base_model = NEW.base_model AND
        cars_catalog.version = NEW.version
    ) THEN
        RAISE EXCEPTION 'Exact duplicate row: This car already exists in the catalog';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Attach the trigger to the cars_catalog table
CREATE TRIGGER check_duplicate_car_insert
    BEFORE INSERT ON cars_catalog
    FOR EACH ROW
EXECUTE FUNCTION check_duplicate_car();