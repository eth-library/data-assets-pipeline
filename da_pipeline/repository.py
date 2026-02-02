from dagster import define_asset_job, repository

from da_pipeline.assets import files, fixities, intellectual_entities, representations, sip_asset
from da_pipeline.sensors import xml_file_sensor

# Define a job that materializes all assets in the correct order
ingest_sip_job = define_asset_job(
    "ingest_sip_job", description="Processes METS XML files into OAIS-compliant SIPs"
)


@repository
def sip_repository():
    """
    Dagster repository for METS XML to OAIS SIP processing pipeline.

    This repository defines a-compliant Submis complete workflow for processing METS XML files
    into OAISsion Information Packages (SIPs). It includes
    all necessary assets, jobs, and sensors for automated processing.

    Components:
        Assets (in processing order):
            1. sip_asset: Parsed SIP model
            2. intellectual_entities: Extracted IEs
            3. representations: File representations
            4. files: Individual file metadata
            5. fixities: File checksums

        Jobs:
            - ingest_sip_job: Orchestrates the complete SIP creation process
              by materializing all assets in the correct order

        Sensors:
            - xml_file_sensor: Monitors for new METS XML files and triggers
              the ingest job automatically

    The processing flow ensures:
        - Proper asset dependency resolution
        - OAIS compliance validation
        - Automated processing of new files
        - Complete metadata extraction
        - Data integrity verification

    Returns:
        List containing all repository components (assets, jobs, sensors)
    """
    return [
        sip_asset,
        intellectual_entities,
        representations,
        files,
        fixities,
        ingest_sip_job,
        xml_file_sensor,
    ]
