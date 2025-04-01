from setuptools import find_packages, setup

setup(
    name="da_pipeline",
    packages=find_packages(exclude=["da_pipeline_tests"]),
    install_requires=[
        "dagster",
        "dagster-cloud",
        "pydantic"
    ],
    extras_require={"dev": ["dagster-webserver", "pytest"]},
)
