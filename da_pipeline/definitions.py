from dagster import Definitions, load_assets_from_modules, EnvVar

from da_pipeline import assets
from da_pipeline.sensors import (
    xml_file_sensor,
    ingest_sip_job,
    TestDataPathResource,
    _default_test_data_path,
)

all_assets = load_assets_from_modules([assets])

defs = Definitions(
    assets=all_assets,
    jobs=[ingest_sip_job],
    sensors=[xml_file_sensor],
    resources={
        "test_data_path": TestDataPathResource(
            path=EnvVar("DAGSTER_TEST_DATA_PATH").get_value(default=_default_test_data_path),
        ),
    },
)
