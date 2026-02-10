from setuptools import setup, find_packages
from pathlib import Path

# Read long description from README
long_description = (Path(__file__).parent / "README.md").read_text(encoding="utf-8")

setup(
    name="hyper-zilla",
    version="1.0.0",
    description="Advanced AI-powered security intelligence platform with proprietary Hyper-ZiLLA technology",
    long_description=long_description,
    long_description_content_type="text/markdown",
    author="Hyper-ZiLLA Team",
    packages=find_packages(include=["HyperZilla", "HyperZilla.*"]),
    package_data={
        "HyperZilla": [
            "COMMAND_CENTER/templates/*.html",
            "COMMAND_CENTER/static/css/*.css",
            "COMMAND_CENTER/static/js/*.js",
            "COMMAND_CENTER/static/js/facial_intel/*.js",
            "SUPPORT_SYSTEMS/configs/*.yaml",
        ]
    },
    install_requires=[
        "flask>=2.3.0",
        "numpy>=1.24.0",
        "opencv-python>=4.7.0",
        "pillow>=9.5.0",
        "requests>=2.31.0",
        "pyyaml>=6.0",
        "scikit-learn>=1.2.0",  # Math library for custom AI
        "torch>=2.0.0",         # ML framework for custom models
        "torchvision>=0.15.0",
        "selenium>=4.10.0",
        "beautifulsoup4>=4.12.0",
        "cryptography>=41.0.0",
        "psutil>=5.9.0",
    ],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "black>=23.0.0",
            "flake8>=6.0.0",
            "bandit>=1.7.5",
        ]
    },
    python_requires=">=3.8",
    entry_points={
        "console_scripts": [
            "hyperzilla=main:main",
        ],
    },
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Intended Audience :: Information Technology",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Topic :: Security",
        "Topic :: Internet :: WWW/HTTP",
        "Topic :: Scientific/Engineering :: Artificial Intelligence",
    ],
)
