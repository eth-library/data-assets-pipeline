from dagster import Definitions, load_assets_from_modules

from . import assets
from .sensors import xml_file_sensor, ingest_sip_job

all_assets = load_assets_from_modules([assets])

defs = Definitions(
    assets=all_assets,
    jobs=[ingest_sip_job],
    sensors=[xml_file_sensor],
)
