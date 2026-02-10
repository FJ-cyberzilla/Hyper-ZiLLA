from setuptools import setup, find_packages

setup(
    name='codex-enterprise',
    version='0.1.0',
    packages=find_packages(where='src'),
    package_dir={'': 'src'},
    install_requires=[
        'jsonschema',
    ],
    entry_points={
        'console_scripts': [
            'codex = codex_enterprise.cli:main',
        ],
    },
    author='Your Name',
    author_email='your.email@example.com',
    description='Cyberzilla Codex - Enterprise Code Quality Analyzer',
    long_description=open('README.md').read(),
    long_description_content_type='text/markdown',
    url='https://github.com/yourusername/codex-enterprise',
    classifiers=[
        'Programming Language :: Python :: 3',
        'License :: OSI Approved :: MIT License',
        'Operating System :: OS Independent',
    ],
    python_requires='>=3.8',
)
